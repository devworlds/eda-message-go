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

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatalf("Failed to upgrade connection: %v", err)
		}
		mockHub.AddClient(conn)
		defer func() {
			mockHub.RemoveClient(conn)
			conn.Close()
		}()
		mockHub.On("RemoveClient", mock.Anything).Run(func(args mock.Arguments) {
			client := args.Get(0).(*websocket.Conn)
			if client.LocalAddr().String() != conn.LocalAddr().String() || client.RemoteAddr().String() != conn.RemoteAddr().String() {
				t.Fatalf("RemoveClient called with unexpected connection: got %v, want %v", client, conn)
			}
		}).Return()
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

	// Mock the BroadcastMessage method to validate the argument
	mockHub.On("BroadcastMessage", mock.MatchedBy(func(message []byte) bool {
		return string(message) == "Hello, WebSocket!"
	})).Return()

	// Allow some time for the message to be processed
	time.Sleep(1 * time.Second)

	// Assert that the BroadcastMessage method was called with the correct message
	mockHub.AssertCalled(t, "BroadcastMessage", message)

	// Close the connection and ensure RemoveClient is called
	err = conn.Close()
	if err != nil {
		t.Fatalf("Failed to close WebSocket connection: %v", err)
	}
}
