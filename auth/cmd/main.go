package main

import (
	"log"
	"net/http"

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

	log.Println("Auth service running on port 8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
