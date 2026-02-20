package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// loadEnv to read .env file and set environment variables

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}
}

// GetEnv to get environment variable with default value

func GetEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}