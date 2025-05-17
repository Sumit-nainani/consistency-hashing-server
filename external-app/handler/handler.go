package handler

import (
	"encoding/json"
	"fmt"
	pb "hashing/hashing"
	"hashing/hashring"
	websocketserver "hashing/websocket-server"
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

func GetServer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
	if clientIP, err := getIP(r); err == nil {
		clientIP = fmt.Sprintf("%s-%d", clientIP, time.Now().UnixNano())
		gotNode := hashring.GetRingInstance().GetNode(clientIP)

		if gotNode {
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

			json.NewEncoder(w).Encode(map[string]string{"You Got Server With Ip": hashring.GetRingInstance().RequestIpToMetaData[clientIP].AssignedNodeIP})
		} else {
			json.NewEncoder(w).Encode(map[string]string{"error": "No Server Available To Serve You.Please Try Again."})
		}

	} else {

	}
}
