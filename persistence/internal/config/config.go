package config

import "os"

type Config struct {
	DatabaseURL  string
	KafkaBrokers []string
	KafkaGroup   string
	KafkaTopic   string
}

func Load() *Config {
	return &Config{
		DatabaseURL:  os.Getenv("DATABASE_URL"),
		KafkaBrokers: []string{os.Getenv("KAFKA_BROKER")},
		KafkaGroup:   os.Getenv("KAFKA_GROUP"),
		KafkaTopic:   os.Getenv("KAFKA_TOPIC"),
	}
}
