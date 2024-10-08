package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	DBSSLMode      string
	ExternalAPIURL string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Printf("[ERROR] Unable to load .env file: %v", err)
		return nil, err
	}

	config := &Config{
		DBHost:         os.Getenv("DB_HOST"),
		DBPort:         os.Getenv("DB_PORT"),
		DBUser:         os.Getenv("DB_USER"),
		DBPassword:     os.Getenv("DB_PASSWORD"),
		DBName:         os.Getenv("DB_NAME"),
		DBSSLMode:      os.Getenv("DB_SSLMODE"),
		ExternalAPIURL: os.Getenv("EXTERNAL_API_URL"),
	}

	if err := validateConfig(config); err != nil {
		return nil, err
	}

	return config, nil
}

// Validate the configuration to ensure all required fields are set
func validateConfig(config *Config) error {
	if config.DBHost == "" || config.DBUser == "" || config.DBPassword == "" || config.DBName == "" {
		return fmt.Errorf("missing required configuration")
	}
	return nil
}
