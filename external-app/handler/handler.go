package handler

import (
	"encoding/json"
	"fmt"
	"hashing/hashring"
	"hashing/utility"
	"net/http"
	"time"
)

// It is the entrypoint of our external application.
// It will work as a middleware which will forward the client request to the correct node/server.
// it will find the correct node according to hash values and hit the kubernetes pod directly.
// At the end it will broadcast the client registeration event to update the UI.
func GetServer(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if clientIP, err := utility.GetIP(r); err == nil {
		clientIP = fmt.Sprintf("%s-%d", clientIP, time.Now().UnixNano())
		hashringInstance := hashring.GetRingInstance()
		gotNode := hashringInstance.GetNode(clientIP)

		assignedNodeIp := hashringInstance.RequestIpToMetaData[clientIP].AssignedNodeIP
		requestHash := hashringInstance.RequestIpToMetaData[clientIP].RequestHash
		assignedNodeName := hashringInstance.RequestIpToMetaData[clientIP].AssignedNodeName
		assignedNodeHash := hashringInstance.RequestIpToMetaData[clientIP].AssignedNodeHash

		// When there is no server/node available right now then we will send dummy data to the UI.
		if !utility.IsNodeAvailable(assignedNodeIp, assignedNodeName, assignedNodeHash) {
			assignedNodeHash = -1
			assignedNodeIp = "0.0.0.0"
			assignedNodeName = "NA"
			json.NewEncoder(w).Encode(map[string]string{
				"error": "No Server Available To Serve You. Please Try Again.",
			})
		} else if err := utility.RunCurlFromCurlPod(assignedNodeIp); err != nil {
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Internal Server Error. Can't Reach Server at IP: " + assignedNodeIp,
			})
		} else {
			json.NewEncoder(w).Encode(map[string]string{
				"message": "You Got Server With IP: " + assignedNodeIp,
			})
		}

		// If a client has come first time then only send the event to UI.
		// because client will never be removed from UI , so no need to send same data again and again.
		if gotNode {
			utility.BroadcastRequestIPMetaData(requestHash, assignedNodeName, assignedNodeIp, assignedNodeHash)
		}
	} else {
		json.NewEncoder(w).Encode(map[string]string{
			"error": "We are not able to find the Ip address of request. Please try again.",
		})
	}
}
