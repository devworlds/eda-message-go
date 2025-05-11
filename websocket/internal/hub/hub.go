package hub

import (
	"fmt"
	"sync"

	"github.com/devworlds/eda-message-go/websocket/internal/interfaces"
	"github.com/gorilla/websocket"
)

// Hub maintains the set of active clients and broadcasts messages to them.
type Hub struct {
	Clients   map[*websocket.Conn]bool
	Broadcast chan []byte
	Mu        sync.Mutex
	Validator interfaces.IValidator
}

// NewHub creates a new Hub instance.
func NewHub(validator interfaces.IValidator) *Hub {
	return &Hub{
		Clients:   make(map[*websocket.Conn]bool),
		Broadcast: make(chan []byte),
		Validator: validator,
	}
}

// HandleWebSocket handles WebSocket connections.
func (h *Hub) Run() {
	for {
		msg := <-h.Broadcast
		fmt.Printf("Run: Message received for broadcast: %s\n", string(msg))

		h.Mu.Lock()
		for client := range h.Clients {
			fmt.Printf("Run: Sending message to client: %v\n", client.RemoteAddr())
			err := client.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				fmt.Printf("Run: Failed to send message to client: %v\n", err)
				if cerr := client.Close(); cerr != nil {
					fmt.Printf("Run: Error closing client: %v\n", cerr)
				}
				delete(h.Clients, client)
			}
		}
		h.Mu.Unlock()
	}
}

// AddClient adds a new client to the Hub after validating the first message as a JWT.
func (h *Hub) AddClient(conn *websocket.Conn) {
	fmt.Printf("AddClient: New connection waiting auth from: %v\n", conn.RemoteAddr())
	_, msg, err := conn.ReadMessage()
	if err != nil {
		fmt.Printf("AddClient: Failed to read first message: %v\n", err)
		conn.Close()
		return
	}

	fmt.Printf("AddClient: Received first message: %s\n", string(msg))
	if !h.Validator.ValidateJWT(string(msg)) {
		fmt.Printf("AddClient: Invalid JWT: %s\n", msg)
		conn.Close()
		return
	}

	h.Mu.Lock()
	h.Clients[conn] = true
	h.Mu.Unlock()
	fmt.Println("HandleWebSocket: Adding client to hub")
}

// RemoveClient removes a client from the Hub.
func (h *Hub) RemoveClient(conn *websocket.Conn) {
	h.Mu.Lock()
	defer h.Mu.Unlock()
	delete(h.Clients, conn)
}
