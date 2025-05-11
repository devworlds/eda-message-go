package kafka

// Package kafka provides a Kafka consumer that listens for messages and broadcasts them to connected WebSocket clients.

import (
	"context"
	"fmt"
	"log"

	"github.com/devworlds/eda-message-go/websocket/internal/hub"
	"github.com/segmentio/kafka-go"
)

// startKafkaConsumer starts a Kafka consumer that listens for messages
func StartKafkaConsumer(h *hub.Hub) {
	go func() {
		r := kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{"kafka:9092"},
			Topic:   "websocket-messages",
			//GroupID: "consumer-group",
		})
		defer r.Close()

		for {
			m, err := r.ReadMessage(context.Background())
			if err != nil {
				log.Printf("Error reading message from Kafka: %v", err)
				continue
			}
			fmt.Printf("Received message from Kafka: %s\n", string(m.Value))
			h.Broadcast <- m.Value
		}
	}()
}
