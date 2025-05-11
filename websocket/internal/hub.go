package websocket

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

// Hub maintains the set of active clients and broadcasts messages to them.
type Hub struct {
	Clients   map[*websocket.Conn]bool
	Broadcast chan []byte
	Mu        sync.Mutex
}

// NewHub creates a new Hub instance.
func NewHub() *Hub {
	return &Hub{
		Clients:   make(map[*websocket.Conn]bool),
		Broadcast: make(chan []byte),
	}
}

// HandleWebSocket handles WebSocket connections.
func (h *Hub) Run() {
	for {
		msg := <-h.Broadcast
		h.Mu.Lock()
		for client := range h.Clients {
			err := client.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				fmt.Printf("Failed to send message to client: %v", err)
				if cerr := client.Close(); cerr != nil {
					fmt.Printf("Error closing client: %v\n", cerr)
				}
				delete(h.Clients, client)
			}
		}
		h.Mu.Unlock()
	}
}

// AddClient adds a new client to the Hub.
func (h *Hub) AddClient(conn *websocket.Conn) {
	h.Mu.Lock()
	h.Clients[conn] = true
	h.Mu.Unlock()
}

// RemoveClient removes a client from the Hub.
func (h *Hub) RemoveClient(conn *websocket.Conn) {
	h.Mu.Lock()
	defer h.Mu.Unlock()
	delete(h.Clients, conn)
}
