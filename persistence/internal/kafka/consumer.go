package kafka

import (
	"fmt"
	"log"

	"github.com/devworlds/eda-message-go/persistence/internal/db"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"gorm.io/gorm"
)

func NewConsumer(brokers []string, group, topic string) *kafka.Consumer {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": brokers[0],
		"group.id":          group,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		log.Fatalf("Persistance: Failed to create consumer: %v", err)
	}
	if err := c.SubscribeTopics([]string{topic}, nil); err != nil {
		log.Printf("Error subscribing to topics: %v", err)
		return nil
	}
	return c
}

func ConsumeLoop(consumer *kafka.Consumer, database *gorm.DB) {
	for {
		msg, err := consumer.ReadMessage(-1)
		if err != nil {
			log.Printf("Persistance: Error reading message: %v", err)
			continue
		}
		if err := db.SaveMessage(database, msg); err != nil {
			log.Printf("Persistance: Error saving message: %v", err)
			continue
		}
		_, err = consumer.CommitMessage(msg)
		if err != nil {
			log.Printf("Error committing message: %v", err)
		}
		fmt.Printf("Persistance: Message consumed: %s\n", msg.Value)
	}
}
