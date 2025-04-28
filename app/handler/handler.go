package handler

import (
	"fmt"
	metrics "goserver/metric"
	"net"
	"net/http"
	"runtime"
	"sync"
	"time"
)

const (
	loop   = 1e10
	worker = 4
)

var (
	mu     sync.Mutex
	ipSeen = make(map[string]time.Time)
)

func trackIp(ip string) {
	mu.Lock()
	defer mu.Unlock()
	ipSeen[ip] = time.Now()
}

func PruneOldIps() {

	for {
		time.Sleep(10 * time.Second)
		mu.Lock()
		cutoff := time.Now().Add(-1 * time.Minute)
		for ip, lastSeen := range ipSeen {
			if lastSeen.Before(cutoff) {
				delete(ipSeen, ip)
			}
		}
		metrics.UniqueIPsGauge.Set(float64(len(ipSeen)))
		mu.Unlock()
	}
}

func GetServer(w http.ResponseWriter, r *http.Request) {

	runtime.GOMAXPROCS(4)
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		fmt.Println(forwarded, "forwarded")
	}

	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		fmt.Println(realIP, "real")
	}

	// Fallback to RemoteAddr
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	trackIp(ip)
	var wg sync.WaitGroup
	chunk := int(loop) / worker

	for i := range worker {
		start := i * chunk
		end := start + chunk
		if i == worker-1 {
			end = int(loop)
		}

		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			sum := 0
			for j := start; j < end; j++ {
				sum += j % 10
			}
			fmt.Println("Thread : ", i)

		}(start, end)
	}

	wg.Wait()
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("<h1>Hello</h1>"))
}
