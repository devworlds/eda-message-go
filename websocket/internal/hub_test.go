package websocket

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/gorilla/websocket"
)

type MockConn struct {
	websocket.Conn
	writeMessageError error
}

func (m *MockConn) WriteMessage(messageType int, data []byte) error {
	return m.writeMessageError
}

func TestHub_AddClient(t *testing.T) {
	hub := NewHub()
	conn := &websocket.Conn{} // Mock connection

	hub.AddClient(conn)

	hub.Mu.Lock()
	defer hub.Mu.Unlock()
	if _, exists := hub.Clients[conn]; !exists {
		t.Fatalf("Expected client to be added to hub")
	}
}

func TestHub_RemoveClient(t *testing.T) {
	hub := NewHub()
	conn := &websocket.Conn{} // Mock connection

	hub.AddClient(conn)
	hub.RemoveClient(conn)

	hub.Mu.Lock()
	defer hub.Mu.Unlock()
	if _, exists := hub.Clients[conn]; exists {
		t.Fatalf("Expected client to be removed from hub")
	}
}

func TestHub_Run(t *testing.T) {
	hub := &Hub{
		Clients:   make(map[*websocket.Conn]bool),
		Broadcast: make(chan []byte, 1),
		Mu:        sync.Mutex{},
	}

	// Create a test WebSocket server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatalf("Failed to upgrade connection: %v", err)
		}
		hub.AddClient(conn)
	}))
	defer ts.Close()

	// Connect to the test WebSocket server
	wsURL := "ws" + ts.URL[4:]
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect to WebSocket server: %v", err)
	}
	defer conn.Close()

	// Simulate a message broadcast
	message := []byte("test message")
	hub.Broadcast <- message

	// Run the hub in a separate goroutine
	go func() {
		hub.Run()
	}()

	// Read the message from the WebSocket connection
	_, receivedMessage, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("Failed to read message: %v", err)
	}

	// Check if the client received the message
	if !bytes.Equal(receivedMessage, message) {
		t.Fatalf("Expected message %s, but got %s", message, receivedMessage)
	}
}

func TestHub_Run_FailureScenario(t *testing.T) {
	hub := &Hub{
		Clients:   make(map[*websocket.Conn]bool),
		Broadcast: make(chan []byte, 1),
		Mu:        sync.Mutex{},
	}

	// Create a test WebSocket server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatalf("Failed to upgrade connection: %v", err)
		}
		// Simulate a write error by closing the connection immediately
		conn.Close()
		hub.AddClient(conn)
	}))
	defer ts.Close()

	// Connect to the test WebSocket server
	wsURL := "ws" + ts.URL[4:]
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect to WebSocket server: %v", err)
	}
	defer conn.Close()

	// Simulate a message broadcast
	message := []byte("test message")
	hub.Broadcast <- message

	// Run the hub in a separate goroutine
	go func() {
		hub.Run()
	}()

	// Allow some time for the message to be processed
	// In a real test, you might use synchronization primitives instead
	// to ensure the message is processed before making assertions.
	// For simplicity, we use a small sleep here.
	// time.Sleep(100 * time.Millisecond)

	// Check if the client was removed from the hub after the write error
	hub.Mu.Lock()
	defer hub.Mu.Unlock()
	if len(hub.Clients) != 0 {
		t.Fatalf("Expected client to be removed from hub after write error")
	}
}
