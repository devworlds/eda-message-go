package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/devworlds/eda-message-go/websocket/internal/handler"
	"github.com/devworlds/eda-message-go/websocket/internal/hub"
	"github.com/devworlds/eda-message-go/websocket/internal/kafka"
)

// Start initializes the WebSocket server and starts the Kafka consumer.
func Start() {
	hub := hub.NewHub()
	go hub.Run()

	kafka.StartKafkaConsumer(hub)

	http.HandleFunc("/ws", handler.HandleWebSocket(hub))
	port := ":8080"
	fmt.Printf("Server is running on http://localhost%s\n", port)

	server := &http.Server{
		Addr:         port,
		Handler:      nil,
		ReadTimeout:  5 * 1e9,
		WriteTimeout: 10 * 1e9,
		IdleTimeout:  120 * 1e9,
	}

	log.Fatal(server.ListenAndServe())
}
