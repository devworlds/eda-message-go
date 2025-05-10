package websocket

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	ws "github.com/devworlds/eda-message-go/websocket/internal"
	"github.com/gorilla/websocket"
)

func TestNewHub(t *testing.T) {
	hub := ws.NewHub()
	if hub == nil {
		t.Fatal("Expected NewHub to return a non-nil Hub")
	}
	if len(hub.Clients) != 0 {
		t.Fatal("Expected no clients in a new Hub")
	}
}

func TestAddClient(t *testing.T) {
	hub := ws.NewHub()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, _ := ws.Upgrader.Upgrade(w, r, nil)
		hub.AddClient(conn)
	}))
	defer server.Close()

	wsURL := "ws" + server.URL[4:]
	_, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect to WebSocket server: %v", err)
	}

	hub.Mu.Lock()
	if len(hub.Clients) != 1 {
		t.Fatalf("Expected 1 client, got %d", len(hub.Clients))
	}
	hub.Mu.Unlock()
}

func TestHandleWebSocket(t *testing.T) {
	hub := ws.NewHub()
	go hub.Run() // Ensure the Hub's Run method is running

	h := ws.HandleWebSocket(hub)
	server := httptest.NewServer(h)
	defer server.Close()

	wsURL := "ws" + server.URL[4:]
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect to WebSocket server: %v", err)
	}
	defer conn.Close()

	// Trigger an immediate broadcast for testing
	ws.TestBroadcastTrigger <- struct{}{}

	// Wait for the message to be sent
	select {
	case <-ws.TestMessageSent:
		fmt.Println("Test: Message sent signal received")
	case <-time.After(5 * time.Second):
		t.Fatal("Timeout waiting for message to be sent")
	}

	messageType, message, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("Failed to read message: %v", err)
	}

	if messageType != websocket.TextMessage {
		t.Fatalf("Expected TextMessage, got %d", messageType)
	}

	expectedMessage := "external message to all clients!"
	if string(message) != expectedMessage {
		t.Fatalf("Expected message '%s', got '%s'", expectedMessage, string(message))
	}
}
