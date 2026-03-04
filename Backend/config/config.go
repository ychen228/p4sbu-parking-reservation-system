package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds all configuration for the application
type Config struct {
	Server  ServerConfig
	MongoDB MongoDBConfig
	JWT     JWTConfig
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// MongoDBConfig holds MongoDB-related configuration
type MongoDBConfig struct {
	URI              string
	Username         string
	Password         string
	Host             string
	Database         string
	ConnectionOptions string
	Timeout          time.Duration
}

// JWTConfig holds JWT-related configuration
type JWTConfig struct {
	Secret        string
	ExpirationHrs int
}

// NewConfig creates a new Config instance
func NewConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         getEnv("SERVER_PORT", "8080"),
			ReadTimeout:  getEnvAsDuration("SERVER_READ_TIMEOUT", 15*time.Second),
			WriteTimeout: getEnvAsDuration("SERVER_WRITE_TIMEOUT", 15*time.Second),
			IdleTimeout:  getEnvAsDuration("SERVER_IDLE_TIMEOUT", 60*time.Second),
		},
		MongoDB: MongoDBConfig{
			Username:         getEnv("MONGODB_USERNAME", ""),
			Password:         getEnv("MONGODB_PASSWORD", ""),
			Host:             getEnv("MONGODB_HOST", "localhost:27017"),
			Database:         getEnv("MONGODB_DATABASE", "campus_parking"),
			ConnectionOptions: getEnv("MONGODB_OPTIONS", ""),
			Timeout:          getEnvAsDuration("MONGODB_TIMEOUT", 30*time.Second),
		},
		JWT: JWTConfig{
			Secret:        getEnv("JWT_SECRET", "your-secret-key"),
			ExpirationHrs: getEnvAsInt("JWT_EXPIRATION_HOURS", 24),
		},
	}
}

// Helper functions for environment variables
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if valueStr, exists := os.LookupEnv(key); exists {
		if value, err := strconv.Atoi(valueStr); err == nil {
			return value
		}
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if valueStr, exists := os.LookupEnv(key); exists {
		if value, err := time.ParseDuration(valueStr); err == nil {
			return value
		}
	}
	return defaultValue
}