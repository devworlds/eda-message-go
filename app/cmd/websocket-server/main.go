package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/devworlds/eda-message-go/internal/handlers"
)

func main() {
	hub := handlers.NewHub()
	go hub.Run()

	http.HandleFunc("/ws", handlers.HandleWebSocket(hub))

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
