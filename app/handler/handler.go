package handler

import (
	"encoding/json"
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

// This method will watch last 1 minute IP burst , If it is greater then threshold (custom metric) then upscaling will be done.
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
	startTime := time.Now()
	forwarded := r.Header.Get("X-Forwarded-For")
	realIP := r.Header.Get("X-Real-IP")
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)

	var clientIP string
	if len(forwarded) != 0 {
		clientIP = forwarded
	} else if len(realIP) != 0 {
		clientIP = realIP
	} else {
		clientIP = ip
	}

	trackIp(clientIP)
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
		}(start, end)
	}

	wg.Wait()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "request served", "duration": time.Since(startTime).String(), "clientIP": clientIP})
}
