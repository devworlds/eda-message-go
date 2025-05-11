package websocket

import (
	"fmt"
	"log"
	"net/http"
)

// Start initializes the WebSocket server and starts the Kafka consumer.
func Start() {
	hub := NewHub()
	go hub.Run()

	startKafkaConsumer(hub)

	http.HandleFunc("/ws", HandleWebSocket(hub))
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
