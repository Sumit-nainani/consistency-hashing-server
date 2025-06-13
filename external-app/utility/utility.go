package utility

import (
	"fmt"
	pb "hashing/hashing"
	websocketserver "hashing/websocket-server"
	"log"
	"net"
	"net/http"
	"os/exec"
)

const (
	pod           string = "pod"
	client        string = "client"
	containerPort string = "8080"
)

// Utility methods for broadcasting websocket events for node creation/deletion and client registeration.
func BroadcastRequestIPMetaData(params ...interface{}) {
	websocketserver.Broadcast(&pb.WebSocketMetadata{
		Type: client,
		Data: &pb.WebSocketMetadata_RequestMetaData{
			RequestMetaData: &pb.RequestMetaData{
				RequestHash:      params[0].(int32),
				AssignedNodeName: params[1].(string),
				AssignedNodeIp:   params[2].(string),
				AssignedNodeHash: params[3].(int32),
			},
		},
	})
}

func BroadcastNodeMetaData(params ...interface{}) {
	websocketserver.Broadcast(&pb.WebSocketMetadata{
		Type:   pod,
		Action: params[0].(string),
		Data: &pb.WebSocketMetadata_NodeMetaData{
			NodeMetaData: &pb.NodeMetaData{
				NodeIp:   params[1].(string),
				NodeName: params[2].(string),
				NodeHash: params[3].(int32),
			},
		},
	})
}

// Utility method for extracting ip address from a client request.
func GetIP(r *http.Request) (string, error) {

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

// Utility method for hitting the node which is deployed on k8s cluster.
// We are hitting directly the kubernets pod for our usecase , but it is not the production idea.
func RunCurlFromCurlPod(podIP string) error {
	cmd := exec.Command("kubectl", "exec", "-n", "curl-pod", "curlpod", "--",
		"curl", "--max-time", "3", fmt.Sprintf("http://%s:%s", podIP, containerPort),
	)
	if _, err := cmd.CombinedOutput(); err != nil {
		return err
	}
	return nil
}

// Utility method for recovery from panic , safe method.
func SafeGo(name string, fn func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Panic in goroutine [%s]: %v", name, r)
			}
		}()
		log.Printf("Starting goroutine: %s", name)
		fn()
	}()
}

// Utility function for finding whether there is at least one node available or not.
func IsNodeAvailable(nodeIp string, nodeName string, nodeHash int32) bool {
	return len(nodeIp) != 0 && len(nodeName) != 0 && nodeHash != 0
}
