package utils

import (
	"os"
	"strconv"
	"time"
)

// GetEnv retrieves environment variable with a default value
func GetEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// GetEnvAsInt retrieves environment variable as integer with a default value
func GetEnvAsInt(key string, defaultValue int) int {
	value := GetEnv(key, strconv.Itoa(defaultValue))
	if intVal, err := strconv.Atoi(value); err == nil {
		return intVal
	}
	return defaultValue
}

// GetEnvAsDuration retrieves environment variable as duration with a default value
func GetEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	value := GetEnv(key, defaultValue.String())
	if duration, err := time.ParseDuration(value); err == nil {
		return duration
	}
	return defaultValue
}
