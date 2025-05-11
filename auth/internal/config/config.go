package config

import "os"

type Config struct {
	JWTSecret   string
	DatabaseURL string
}

func Load() *Config {
	return &Config{
		JWTSecret:   os.Getenv("JWT_SECRET"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
	}
}
