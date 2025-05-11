package db

import "time"

type Message struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}
