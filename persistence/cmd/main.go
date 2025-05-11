package main

import (
	"fmt"

	"github.com/devworlds/eda-message-go/persistence/internal/config"
	"github.com/devworlds/eda-message-go/persistence/internal/db"
	"github.com/devworlds/eda-message-go/persistence/internal/kafka"
)

func main() {
	cfg := config.Load()
	database, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		panic(err)
	}

	// Initialize Kafka consumer
	consumer := kafka.NewConsumer(cfg.KafkaBrokers, cfg.KafkaGroup, cfg.KafkaTopic)
	fmt.Println("Persistence service initialized")
	kafka.ConsumeLoop(consumer, database)
}
