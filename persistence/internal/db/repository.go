package db

import (
	"encoding/json"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SaveMessage(database *gorm.DB, msg *kafka.Message) error {
	var message Message
	if err := json.Unmarshal(msg.Value, &message); err != nil {
		return err
	}
	// Idempotent insert (upsert)
	if err := database.Clauses(clause.OnConflict{DoNothing: true}).Create(&message).Error; err != nil {
		return err
	}
	return nil
}
