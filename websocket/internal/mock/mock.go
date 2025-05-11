package mock

import (
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/mock"
)

type MockHub struct {
	mock.Mock
}

func (m *MockHub) AddClient(client *websocket.Conn) {
	m.Called(client)
}

func (m *MockHub) Run() {
	m.Called()
}

func (m *MockHub) RemoveClient(client *websocket.Conn) {
	m.Called(client)
}

func (m *MockHub) BroadcastMessage(message []byte) {
	m.Called(message)
}
