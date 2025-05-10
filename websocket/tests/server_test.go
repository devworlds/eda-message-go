package websocket

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/mock"
)

type MockHub struct {
	mock.Mock
}

func (m *MockHub) AddClient(client *websocket.Conn) {
	m.Called(client)
}

func (m *MockHub) RemoveClient(client *websocket.Conn) {
	m.Called(client)
}

func (m *MockHub) BroadcastMessage(message []byte) {
	m.Called(message)
}

func TestHandleWebSocket(t *testing.T) {
	mockHub := new(MockHub)

	// Mock the AddClient method
	mockHub.On("AddClient", mock.Anything).Return()
	// Mock the RemoveClient method
	mockHub.On("RemoveClient", mock.Anything).Return()
	// Mock the BroadcastMessage method
	mockHub.On("BroadcastMessage", mock.Anything).Return()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Upgrade(w, r, nil, 1024, 1024)
		if err != nil {
			t.Fatalf("Failed to upgrade connection: %v", err)
		}
		mockHub.AddClient(conn)
		defer func() {
			mockHub.RemoveClient(conn)
			conn.Close()
		}()
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				break
			}
			mockHub.BroadcastMessage(message)
		}
	}))
	defer server.Close()

	wsURL := "ws" + server.URL[4:]
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect to WebSocket server: %v", err)
	}
	defer conn.Close()

	// Simulate sending a message to the server
	message := []byte("Hello, WebSocket!")
	err = conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		t.Fatalf("Failed to send message: %v", err)
	}

	// Allow some time for the message to be processed
	time.Sleep(1 * time.Second)

	// Assert that the BroadcastMessage method was called with the correct message
	mockHub.AssertCalled(t, "BroadcastMessage", message)

	// Close the connection and ensure RemoveClient is called
	err = conn.Close()
	if err != nil {
		t.Fatalf("Failed to close WebSocket connection: %v", err)
	}
	mockHub.AssertCalled(t, "RemoveClient", mock.Anything)
}
