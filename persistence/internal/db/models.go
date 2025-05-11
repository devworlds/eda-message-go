package db

import "time"

type Message struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

type User struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
