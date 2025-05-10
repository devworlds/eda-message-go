package websocket

import (
	"testing"

	"github.com/gorilla/websocket"
)

func TestMockHub_AddClient(t *testing.T) {
	mockHub := new(MockHub)
	client := &websocket.Conn{}

	mockHub.On("AddClient", client).Return().Once()

	mockHub.AddClient(client)

	mockHub.AssertCalled(t, "AddClient", client)
	mockHub.AssertExpectations(t)
}

func TestMockHub_Run(t *testing.T) {
	mockHub := new(MockHub)

	mockHub.On("Run").Return().Once()

	mockHub.Run()

	mockHub.AssertCalled(t, "Run")
	mockHub.AssertExpectations(t)
}

func TestMockHub_RemoveClient(t *testing.T) {
	mockHub := new(MockHub)
	client := &websocket.Conn{}

	mockHub.On("RemoveClient", client).Return().Once()

	mockHub.RemoveClient(client)

	mockHub.AssertCalled(t, "RemoveClient", client)
	mockHub.AssertExpectations(t)
}

func TestMockHub_BroadcastMessage(t *testing.T) {
	mockHub := new(MockHub)
	message := []byte("test message")

	mockHub.On("BroadcastMessage", message).Return().Once()

	mockHub.BroadcastMessage(message)

	mockHub.AssertCalled(t, "BroadcastMessage", message)
	mockHub.AssertExpectations(t)
}
