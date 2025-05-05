package websocketserver

import (
	"fmt"
	pb "hashing/hashing"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // allow all origins
}

type Hub struct {
	clients   map[*websocket.Conn]bool
	broadcast chan *pb.WebSocketMetadata
	lock      sync.RWMutex
}

var hub = &Hub{
	clients:   make(map[*websocket.Conn]bool),
	broadcast: make(chan *pb.WebSocketMetadata, 100),
}

func GetHub() *Hub {
	return hub
}

func (h *Hub) Register(conn *websocket.Conn) {
	h.lock.Lock()
	defer h.lock.Unlock()
	h.clients[conn] = true
}

func (h *Hub) Unregister(conn *websocket.Conn) {
	h.lock.Lock()
	defer h.lock.Unlock()
	delete(h.clients, conn)
	conn.Close()
}

func (h *Hub) Run() {
	for msg := range h.broadcast {
		fmt.Println(msg, "message")
		serialized, err := proto.Marshal(msg)
		if err != nil {
			log.Println("Failed to serialize:", err)
			continue
		}

		h.lock.RLock()
		for client := range h.clients {
			if err := client.WriteMessage(websocket.BinaryMessage, serialized); err != nil {
				log.Println("Write error:", err)
				h.lock.RUnlock()
				h.Unregister(client)
				h.lock.RLock()
			}
		}
		h.lock.RUnlock()
	}
}

func Broadcast(msg *pb.WebSocketMetadata) {
	GetHub().broadcast <- msg
}

func HandleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	// defer conn.Close()

	log.Println("WebSocket client connected")

	// hashinstance := hashring.GetRingInstance()

	// // Listen for new pods
	// go func() {
	// 	for podData := range kubeclient.Ch {
	// 		pod := Pod{
	// 			Podip:   podData.Podip,
	// 			PodName: podData.Podname,
	// 			Podhash: hashinstance.GetNodeHashValue(podData.Podip),
	// 		}

	// 		podMsg := Message{
	// 			Type: "pod",
	// 			Data: pod,
	// 		}

	// 		if err := conn.WriteJSON(podMsg); err != nil {
	// 			log.Println("Error sending pod message:", err)
	// 			continue
	// 		} else {
	// 			fmt.Println("data sent")
	// 		}
	// 	}
	// }()

	// // Listen for new client IPs
	// go func() {
	// 	for ip := range handler.Ch {
	// 		clientIP := ClientIP{
	// 			ClientIP: ip,
	// 			IPhash:   hashinstance.IPToKey[ip],
	// 		}

	// 		clientMsg := Message{
	// 			Type: "client",
	// 			Data: clientIP,
	// 		}

	// 		if err := conn.WriteJSON(clientMsg); err != nil {
	// 			log.Println("Error sending client IP message:", err)
	// 			continue
	// 		}
	// 	}
	// }()
	// // Keep the connection open
	// select {}
	GetHub().Register(conn)
	go func() {
		defer GetHub().Unregister(conn)
		for {
			if _, _, err := conn.NextReader(); err != nil {
				break
			}
		}
	}()
}
