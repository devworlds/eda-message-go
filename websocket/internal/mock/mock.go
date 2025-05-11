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

type MockProducer struct {
	mock.Mock
}

func (m *MockProducer) SendMessage(topic string, message []byte) error {
	args := m.Called(topic, message)
	return args.Error(0)
}

type MockTokenValidator struct {
	mock.Mock
}

func (m *MockTokenValidator) ValidateJWT(token string) bool {
	args := m.Called(token)
	return args.Bool(1)
}
