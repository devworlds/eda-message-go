package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/mock"

	m "github.com/devworlds/eda-message-go/websocket/internal/mock"
)

func TestHandleWebSocket(t *testing.T) {
	mockHub := new(m.MockHub) // Assuming MockHub is defined in the mock.go file

	// Mock the AddClient method
	mockHub.On("AddClient", mock.Anything).Return()

	handler := HandleWebSocket(mockHub)
	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	wsURL := "ws" + ts.URL[4:]
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect to WebSocket server: %v", err)
	}
	defer conn.Close()

	// Allow some time for the AddClient call to be processed
	mockHub.AssertNumberOfCalls(t, "AddClient", 1)
	mockHub.AssertExpectations(t)
}

func TestHandleWebSocket_UpgradeConnection(t *testing.T) {
	mockHub := new(m.MockHub)

	// Mock the AddClient method
	mockHub.On("AddClient", mock.Anything).Return()

	handler := HandleWebSocket(mockHub)
	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	wsURL := "ws" + ts.URL[4:]
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect to WebSocket server: %v", err)
	}
	defer conn.Close()

	// Assert that AddClient was called exactly once
	mockHub.AssertNumberOfCalls(t, "AddClient", 1)
	mockHub.AssertExpectations(t)
}
