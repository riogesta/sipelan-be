package config

import (
	"os"
	"strconv"
)

type Config struct {
	ServerPort  string
	DatabaseURL string
}

func Load() *Config {
	return &Config{
		ServerPort:  getEnv("SERVER_PORT", "8081"),
		DatabaseURL: getEnv("DATABASE_URL", "sipelan.db"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}
