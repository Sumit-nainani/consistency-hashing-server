package websocketserver

import (
	pb "hashing/hashing"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

var (
	once        sync.Once
	hubInstance *Hub
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // allow all origins
}

type Hub struct {
	clients   map[*websocket.Conn]bool
	broadcast chan *pb.WebSocketMetadata
	lock      sync.RWMutex
}

// Singleton Pattern.
func GetHubInstance() *Hub {
	once.Do(func() {
		hubInstance = &Hub{
			clients:   make(map[*websocket.Conn]bool),
			broadcast: make(chan *pb.WebSocketMetadata, 100),
		}
	})
	return hubInstance
}

// Registering new Websocket client after reloading browser or network error.
func (h *Hub) Register(conn *websocket.Conn) {
	h.lock.Lock()
	defer h.lock.Unlock()
	h.clients[conn] = true
}

// Unregistering old Websocket client.
func (h *Hub) Unregister(conn *websocket.Conn) {
	h.lock.Lock()
	defer h.lock.Unlock()
	delete(h.clients, conn)
	conn.Close()
}

// This method is used to listen from channel continuously when an event is inserted into broadcast channel.
// Here We are sending serialized/marshalled data for efficiently using protobuff and reducing data size in the network.
// We are broadcasting the data to all active websocket clients.
func (h *Hub) Run() {
	for messageEvent := range h.broadcast {
		serializedEvent, err := proto.Marshal(messageEvent)
		if err != nil {
			log.Println("Failed to serialize:", err)
			continue
		}

		h.lock.RLock()
		for client := range h.clients {
			if err := client.WriteMessage(websocket.BinaryMessage, serializedEvent); err != nil {
				log.Println("Write error:", err)
				h.lock.RUnlock()
				h.Unregister(client)
				h.lock.RLock()
			}
		}
		h.lock.RUnlock()
	}
}

func Broadcast(messageEvent *pb.WebSocketMetadata) {
	GetHubInstance().broadcast <- messageEvent
}

func HandleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	log.Println("WebSocket client connected")

	GetHubInstance().Register(conn)

	// Here this method is used to unregister the dead websocket client.
	// If It is not getting anything from client side then the websocket client is dead , so just unregister it.
	go func() {
		defer GetHubInstance().Unregister(conn)
		for {
			if _, _, err := conn.NextReader(); err != nil {
				break
			}
		}
	}()
}
