package websocket

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func Start() {
	hub := NewHub()
	go hub.Run()

	go func() {
		for {
			select {
			case <-TestBroadcastTrigger:
				hub.Broadcast <- []byte("external message to all clients!")
				TestMessageSent <- struct{}{}
			case <-time.After(10 * time.Second):
				hub.Broadcast <- []byte("external message to all clients!")
			}
		}
	}()

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
