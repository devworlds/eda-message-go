package db

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if !db.Migrator().HasTable(&Message{}) {
		if err := db.AutoMigrate(&Message{}); err != nil {
			log.Printf("Error migrating Message table: %v", err)
		}
	}
	if !db.Migrator().HasTable(&User{}) {
		if err := db.AutoMigrate(&User{}); err != nil {
			log.Printf("Error migrating User table: %v", err)
		}
		// Insert initial users
		initialUsers := []User{
			{ID: "1", Username: "client1", Password: "password1"},
			{ID: "2", Username: "client2", Password: "password2"},
			{ID: "3", Username: "client3", Password: "password3"},
		}
		db.Create(&initialUsers)
	}
	return db, nil
}
