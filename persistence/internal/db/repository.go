package db

import (
	"encoding/json"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"gorm.io/gorm"
)

func SaveMessage(database *gorm.DB, msg *kafka.Message) error {
	var message Message
	if err := json.Unmarshal(msg.Value, &message); err != nil {
		return err
	}

	var existingMessage Message
	if err := database.First(&existingMessage, "id = ?", message.ID).Error; err == nil {
		return fmt.Errorf("mensage with ID %s already exist", message.ID)
	} else if err != gorm.ErrRecordNotFound {
		return err
	}

	if err := database.Create(&message).Error; err != nil {
		return err
	}
	return nil
}
