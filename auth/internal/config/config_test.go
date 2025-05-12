package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("JWT_SECRET", "test-secret")
	os.Setenv("DATABASE_URL", "test-database-url")

	config := Load()

	if config == nil {
		t.Fatal("expected a valid config object, got nil")
	}

	if config.JWTSecret != "test-secret" {
		t.Errorf("expected JWTSecret to be 'test-secret', got %v", config.JWTSecret)
	}

	if config.DatabaseURL != "test-database-url" {
		t.Errorf("expected DatabaseURL to be 'test-database-url', got %v", config.DatabaseURL)
	}
}
