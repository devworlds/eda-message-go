package handler

import (
	"fmt"
	"net/http"

	"github.com/devworlds/eda-message-go/websocket/internal/interfaces"
	"github.com/gorilla/websocket"
)

// upgrader is used to upgrade HTTP connections to WebSocket connections.
var Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// HandleWebSocket handles WebSocket connections, adds them to the Hub, and sends messages to Kafka.
func HandleWebSocket(hub interfaces.IHub, producer interfaces.IProducer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("HandleWebSocket: New connection")
		conn, err := Upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println("HandleWebSocket: Failed to upgrade connection:", err)
			return
		}

		hub.AddClient(conn)
		
		// Start a goroutine to read messages from the WebSocket connection
		go func() {
			defer conn.Close()
			for {
				fmt.Println("HandleWebSocket: Waiting to read message")
				_, message, err := conn.ReadMessage()
				if err != nil {
					fmt.Println("HandleWebSocket: Error reading message:", err)
					break
				}

				fmt.Println("HandleWebSocket: Message received, sending to producer")
				err = producer.SendMessage("websocket-messages", message)
				if err != nil {
					fmt.Println("HandleWebSocket: Failed to send message to Kafka:", err)
				} else {
					fmt.Println("HandleWebSocket: Message sent to Kafka successfully")
				}
			}
		}()
	}
}
