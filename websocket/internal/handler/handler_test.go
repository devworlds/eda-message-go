package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/mock"

	m "github.com/devworlds/eda-message-go/websocket/internal/mock"
)

func TestHandleWebSocket(t *testing.T) {
	mockHub := new(m.MockHub)           // Assuming MockHub is defined in the mock.go file
	mockProducer := new(m.MockProducer) // MockProducer for testing

	// Mock the AddClient method
	mockHub.On("AddClient", mock.Anything).Return()
	// Mock the SendMessage method
	mockProducer.On("SendMessage", "websocket-messages", mock.Anything).Return(nil)

	handler := HandleWebSocket(mockHub, mockProducer)
	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	wsURL := "ws" + ts.URL[4:]
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect to WebSocket server: %v", err)
	}
	defer conn.Close()

	// Send a test message through the WebSocket
	testMessage := []byte("test message")
	err = conn.WriteMessage(websocket.TextMessage, testMessage)
	if err != nil {
		t.Fatalf("Failed to send message: %v", err)
	}

	// Allow some time for the AddClient and SendMessage calls to be processed
	time.Sleep(100 * time.Millisecond)

	mockHub.AssertNumberOfCalls(t, "AddClient", 1)
	mockProducer.AssertNumberOfCalls(t, "SendMessage", 1) // Message sent
	mockHub.AssertExpectations(t)
	mockProducer.AssertExpectations(t)
}

func TestHandleWebSocket_UpgradeConnection(t *testing.T) {
	mockHub := new(m.MockHub)
	mockProducer := new(m.MockProducer) // MockProducer for testing

	// Mock the AddClient method
	mockHub.On("AddClient", mock.Anything).Return()
	// Mock the SendMessage method
	mockProducer.On("SendMessage", "websocket-messages", mock.Anything).Return(nil)

	handler := HandleWebSocket(mockHub, mockProducer)
	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	wsURL := "ws" + ts.URL[4:]
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect to WebSocket server: %v", err)
	}
	defer conn.Close()

	// Log connection success
	t.Log("WebSocket connection established")

	// Send a test message through the WebSocket
	testMessage := []byte("test message")
	err = conn.WriteMessage(websocket.TextMessage, testMessage)
	if err != nil {
		t.Fatalf("Failed to send message: %v", err)
	}

	// Allow some time for the AddClient and SendMessage calls to be processed
	time.Sleep(500 * time.Millisecond)

	// Assert that AddClient was called exactly once
	mockHub.AssertNumberOfCalls(t, "AddClient", 1)
	mockProducer.AssertNumberOfCalls(t, "SendMessage", 1) // Message sent
	mockHub.AssertExpectations(t)
	mockProducer.AssertExpectations(t)
}
