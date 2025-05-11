package server

import (
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestStart(t *testing.T) {
	// Start the server in a goroutine
	go Start()

	// Allow the server to start
	time.Sleep(1 * time.Second)

	// Create a WebSocket connection
	wsURL := "ws://localhost:8080/ws"
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to establish WebSocket connection: %v", err)
	}
	defer ws.Close()

	// Verify WebSocket connection
	if ws.UnderlyingConn() == nil {
		t.Errorf("Expected a valid WebSocket connection, got nil")
	}
}
