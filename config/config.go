package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	GeminiAPIKey string
	Port         string
	GinMode      string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	config := &Config{
		GeminiAPIKey: os.Getenv("GEMINI_API_KEY"),
		Port:         getEnvOrDefault("PORT", "8080"),
		GinMode:      getEnvOrDefault("GIN_MODE", "release"),
	}

	return config
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func (c *Config) IsValid() bool {
	return c.GeminiAPIKey != ""
}
