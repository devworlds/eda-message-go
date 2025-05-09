package main

import (
	//standard
	"fmt"
	"log"
	"net/http"

	//internal
	"github.com/devworlds/eda-message-go/internal/handlers"
)

func main() {
	http.HandleFunc("/ws", handlers.HandleWebSocket)

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	port := ":8080"
	fmt.Printf("Server is running on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
