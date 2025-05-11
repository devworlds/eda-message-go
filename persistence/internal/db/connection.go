package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if !db.Migrator().HasTable(&Message{}) {
		db.AutoMigrate(&Message{})
	}
	if !db.Migrator().HasTable(&User{}) {
		db.AutoMigrate(&User{})
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
