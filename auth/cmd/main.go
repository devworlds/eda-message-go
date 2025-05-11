package main

import (
	"log"
	"net/http"
	"time"

	"github.com/devworlds/eda-message-go/auth/internal/config"
	"github.com/devworlds/eda-message-go/auth/internal/db"
	"github.com/devworlds/eda-message-go/auth/internal/jwt"
)

func main() {
	cfg := config.Load()
	database, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		return
	}

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		jwt.LoginHandler(w, r, database)
	})

	server := &http.Server{
		Addr:         ":8081",
		Handler:      nil,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Println("Auth service running on port 8081")
	log.Fatal(server.ListenAndServe())
}
