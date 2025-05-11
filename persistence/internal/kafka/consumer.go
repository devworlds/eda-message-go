package kafka

import (
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
		log.Fatalf("Failed to create consumer: %v", err)
	}
	c.SubscribeTopics([]string{topic}, nil)
	return c
}

func ConsumeLoop(consumer *kafka.Consumer, database *gorm.DB) {
	for {
		msg, err := consumer.ReadMessage(-1)
		if err != nil {
			log.Printf("Error reading message: %v", err)
			continue
		}
		if err := db.SaveMessage(database, msg); err != nil {
			log.Printf("Error saving message: %v", err)
			continue
		}
		consumer.CommitMessage(msg)
	}
}
