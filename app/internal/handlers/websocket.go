package handlers

import (
	//standard
	"fmt"
	"net/http"
	"sync"
	"time"

	//third-party
	"github.com/gorilla/websocket"
)

// upgrader is used to upgrade HTTP connections to WebSocket connections.
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Update the testBroadcastTrigger and testMessageSent channels to be buffered.
var testBroadcastTrigger = make(chan struct{}, 1)
var testMessageSent = make(chan struct{}, 1)

// Hub maintains the set of active clients and broadcasts messages to them.
type Hub struct {
	clients   map[*websocket.Conn]bool
	broadcast chan []byte
	mu        sync.Mutex
}

// NewHub creates a new Hub instance.
func NewHub() *Hub {
	return &Hub{
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan []byte),
	}
}

// Run starts the Hub's message broadcasting loop.
func (h *Hub) Run() {
	for {
		msg := <-h.broadcast
		h.mu.Lock()
		for client := range h.clients {
			err := client.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				client.Close()
				delete(h.clients, client)
			}
		}
		h.mu.Unlock()
	}
}

// AddClient adds a new client to the Hub's list of clients.
func (h *Hub) AddClient(conn *websocket.Conn) {
	h.mu.Lock()
	h.clients[conn] = true
	h.mu.Unlock()
}

// Add more detailed logging to debug the message flow.
func HandleWebSocket(hub *Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println("Failed to upgrade connection:", err)
			return
		}
		defer func() {
			fmt.Println("Closing connection")
			conn.Close()
		}()

		hub.AddClient(conn)
		fmt.Println("Client connected")

		for {
			select {
			case <-testBroadcastTrigger:
				fmt.Println("Trigger received, broadcasting message...")
				hub.broadcast <- []byte("external message to all clients!")
				select {
				case testMessageSent <- struct{}{}:
					fmt.Println("Message sent signal emitted")
				default:
					fmt.Println("Message sent signal skipped (buffer full)")
				}
			case <-time.After(10 * time.Second):
				fmt.Println("Broadcasting periodic message...")
				hub.broadcast <- []byte("external message to all clients!")
			}
		}
	}
}
