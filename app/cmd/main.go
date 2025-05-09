package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/devworlds/eda-message-go/internal/handlers"
)

func main() {
	http.HandleFunc("/ws", handlers.HandleWebSocket)
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
	port := ":8080"
	fmt.Printf("Server is running on http://localhost%s\n", port)

	server := &http.Server{
		Addr:         port,
		Handler:      nil, // use o DefaultServeMux
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}
