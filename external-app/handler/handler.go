package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"hashing/client"
	pb "hashing/hashing"
	"hashing/hashring"
	websocketserver "hashing/websocket-server"
	"log"
	"net"
	"net/http"
	"os/exec"
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

// func AddFirstServer(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
// 	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
// 	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

// 	cmd := exec.Command("kubectl", "apply", "-f", "../kubernetes/deployment.yaml")
// 	output, err := cmd.CombinedOutput()
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("Failed to apply manifest: %s", output), http.StatusInternalServerError)
// 		return
// 	}
// 	hashring := hashring.GetRingInstance()

// 	timeout := time.After(30 * time.Second)
// 	ticker := time.NewTicker(1 * time.Second)
// 	defer ticker.Stop()

// 	for {
// 		select {
// 		case <-timeout:
// 			http.Error(w, "Timeout waiting for node to be added to the hash ring", http.StatusGatewayTimeout)
// 			return
// 		case <-ticker.C:
// 			if hashring.GetLenNodes() == 1 {
// 				hashring.RLock()
// 				hashValue := hashring.NodeMetaData[0]
// 				hashring.RUnlock()

// 				// Respond with the hash value in JSON format
// 				w.Header().Set("Content-Type", "application/json")
// 				json.NewEncoder(w).Encode(map[string]int{"hash": hashValue})
// 				return
// 			}
// 		}
// 	}
// }

func GetServer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if clientIP, err := getIP(r); err == nil {

		hashring.GetRingInstance().GetNode(clientIP)
		websocketserver.Broadcast(&pb.WebSocketMetadata{
			Type: "client",
			Data: &pb.WebSocketMetadata_RequestMetaData{
				RequestMetaData: &pb.RequestMetaData{
					AssignedNodeName: hashring.GetRingInstance().RequestIpToMetaData[clientIP].AssignedNodeName,
					AssignedNodeIp:   hashring.GetRingInstance().RequestIpToMetaData[clientIP].AssignedNodeIP,
					RequestHash:      hashring.GetRingInstance().RequestIpToMetaData[clientIP].RequestHash,
					AssignedNodeHash: hashring.GetRingInstance().RequestIpToMetaData[clientIP].AssignedNodeHash,
				},
			},
		})

		err = runCurlFromCurlPod(hashring.GetRingInstance().RequestIpToMetaData[clientIP].AssignedNodeIP)
		if err != nil {
			log.Printf("Error running curl from curl-pod: %v", err)
		} else {
			log.Printf("Successfully ran curl to %s", hashring.GetRingInstance().RequestIpToMetaData[clientIP].AssignedNodeIP)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]int{"hash": int(hashring.GetRingInstance().RequestIpToMetaData[clientIP].RequestHash)})
	} else {

	}
}

func Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	c := client.GetClient()
	list, _ := c.GetHashRingData(context.Background(), &pb.Empty{})
	fmt.Println(list, "list")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[int]string{1: "hello"})
}
