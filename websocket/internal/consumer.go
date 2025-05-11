package websocket

import (
	"context"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

// startKafkaConsumer starts a Kafka consumer that listens for messages
func startKafkaConsumer(hub *Hub) {
	go func() {
		r := kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{"kafka:9092"},
			Topic:   "websocket-messages",
		})
		defer r.Close()

		for {
			m, err := r.ReadMessage(context.Background())
			if err != nil {
				log.Printf("Error reading message from Kafka: %v", err)
				continue
			}
			fmt.Printf("Received message from Kafka: %s\n", string(m.Value))
			hub.Broadcast <- m.Value
		}
	}()
}
