package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"hashing/client"
	ph "hashing/hashing"
	"hashing/hashring"
	"log"
	"net"
	"net/http"
	"os/exec"
	"time"
)

const (
	containerPort string = "8080"
)

func getIP(r *http.Request) (string, error) {

	forwarded := r.Header.Get("X-Forwarded-For")

	if forwarded != "" {
		return forwarded, nil
	}

	realIP := r.Header.Get("X-Real-IP")

	if realIP != "" {
		return realIP, nil
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	fmt.Println("ipppp", ip)
	if err != nil {
		return "", err
	} else {
		return ip, err
	}
}
func runCurlFromCurlPod(podIP string) error {
	cmd := exec.Command("kubectl", "exec", "-n", "curl-pod", "curlpod", "--",
		"curl", fmt.Sprintf("http://%s:%s", podIP, containerPort),
	)
	_, err := cmd.CombinedOutput()

	if err != nil {
		return err
	}
	return nil
}

func AddFirstServer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	cmd := exec.Command("kubectl", "apply", "-f", "../kubernetes/deployment.yaml")
	output, err := cmd.CombinedOutput()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to apply manifest: %s", output), http.StatusInternalServerError)
		return
	}
	hashring := hashring.GetRingInstance()

	timeout := time.After(30 * time.Second)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			http.Error(w, "Timeout waiting for node to be added to the hash ring", http.StatusGatewayTimeout)
			return
		case <-ticker.C:
			if hashring.GetLenNodes() == 1 {
				hashring.RLock()
				hashValue := hashring.Nodes[0]
				hashring.RUnlock()

				// Respond with the hash value in JSON format
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]int{"hash": hashValue})
				return
			}
		}
	}
}

func GetServer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if clientIP, err := getIP(r); err == nil {
		c := client.GetClient()
		response, err := c.GetNodeForRequest(context.Background(), &ph.NodeRequest{Key: clientIP})
		if err != nil {
			log.Fatalf("Error getting node IP for clientIP %s: %v", clientIP, err)
		}

		err = runCurlFromCurlPod(response.Node)
		if err != nil {
			log.Printf("Error running curl from curl-pod: %v", err)
		} else {
			log.Printf("Successfully ran curl to %s", response.Node)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]int{"hash": int(response.Hash)})
	} else {

	}
}
