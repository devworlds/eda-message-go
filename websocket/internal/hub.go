package websocket

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

// Update the testBroadcastTrigger and testMessageSent channels to be buffered.
var TestBroadcastTrigger = make(chan struct{}, 1)
var TestMessageSent = make(chan struct{}, 1)

type Hub struct {
	Clients   map[*websocket.Conn]bool
	Broadcast chan []byte
	Mu        sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		Clients:   make(map[*websocket.Conn]bool),
		Broadcast: make(chan []byte),
	}
}

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

func (h *Hub) AddClient(conn *websocket.Conn) {
	h.Mu.Lock()
	h.Clients[conn] = true
	h.Mu.Unlock()
}
