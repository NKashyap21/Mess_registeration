package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnvVariables() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Set default values if not provided
	setDefaultEnv("DB_HOST", "localhost")
	setDefaultEnv("DB_PORT", "5432")
	setDefaultEnv("DB_NAME", "mess_registration")
	setDefaultEnv("DB_USER", "postgres")
	setDefaultEnv("DB_PASSWORD", "password")
	setDefaultEnv("JWT_SECRET", "your-super-secret-jwt-key-change-in-production")
	setDefaultEnv("GIN_MODE", "debug")
	setDefaultEnv("MESS_API_KEY", "mess-scanner-api-key-12345")
}

func setDefaultEnv(key, defaultValue string) {
	if os.Getenv(key) == "" {
		os.Setenv(key, defaultValue)
	}
}
