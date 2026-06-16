package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Config holds all application configuration values.
type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	ServerPort string
}

// Load reads configuration from environment variables (with .env fallback).
func Load() *Config {
	// Load .env file if present; ignore errors (production won't have it).
	_ = godotenv.Load()

	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "user_management"),
		ServerPort: getEnv("SERVER_PORT", "3000"),
	}
}

// getEnv retrieves an environment variable or returns a fallback default.
func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
