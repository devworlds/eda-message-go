package interfaces

import "github.com/gorilla/websocket"

type IHub interface {
	AddClient(client *websocket.Conn)
	Run()
}

type IProducer interface {
	SendMessage(topic string, message []byte) error
}

// Validator interface for dependency injection
type IValidator interface {
	ValidateJWT(token string) bool
}
