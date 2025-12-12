package config

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ValidateConfig checks all required environment variables and formats
func ValidateConfig(dbMode string) error {

	requiredVars := map[string]string{
		"BASE_API":                  os.Getenv("BASE_API"),
		"TLS_PATH":                  os.Getenv("TLS_PATH"),
		"SECRET_KEY":                os.Getenv("SECRET_KEY"),
		"TOKEN_EXPIRATION_PERIOD":   os.Getenv("TOKEN_EXPIRATION_PERIOD"),
	}

	// Only require DB variables if using PostgreSQL
	if dbMode != "mock" {
		requiredVars["DB_HOST"] = os.Getenv("DB_HOST")
		requiredVars["DB_PORT"] = os.Getenv("DB_PORT")
		requiredVars["DB_USER_NAME"] = os.Getenv("DB_USER_NAME")
		requiredVars["DB_PASSWORD"] = os.Getenv("DB_PASSWORD")
		requiredVars["DB_NAME"] = os.Getenv("DB_NAME")
	}

	// Check required variables are set
	for varName, value := range requiredVars {
		if value == "" {
			return fmt.Errorf("required environment variable not set: %s", varName)
		}
	}

	// Validate BASE_API format
	if !strings.HasPrefix(BaseApi, "/") {
		return fmt.Errorf("BASE_API must start with /")
	}
	if strings.Contains(BaseApi, " ") {
		return fmt.Errorf("BASE_API contains invalid characters")
	}

	// Validate SECRET_KEY length (should be at least 32 chars for security)
	if len(SecretKey) < 32 {
		return fmt.Errorf("SECRET_KEY must be at least 32 characters long for security")
	}

	// Validate TOKEN_EXPIRATION_PERIOD is valid duration
	if _, err := time.ParseDuration(os.Getenv("TOKEN_EXPIRATION_PERIOD")); err != nil {
		return fmt.Errorf("invalid TOKEN_EXPIRATION_PERIOD format: %v", err)
	}

	// Validate DB_PORT is numeric
	dbPort := os.Getenv("DB_PORT")
	if _, err := strconv.Atoi(dbPort); err != nil {
		return fmt.Errorf("DB_PORT must be numeric, got: %s", dbPort)
	}

	// Validate optional integer configs
	if dbMaxOpenConnsStr := os.Getenv("DB_MAX_OPEN_CONNS"); dbMaxOpenConnsStr != "" {
		if _, err := strconv.Atoi(dbMaxOpenConnsStr); err != nil {
			return fmt.Errorf("DB_MAX_OPEN_CONNS must be numeric, got: %s", dbMaxOpenConnsStr)
		}
	}

	if dbMaxIdleConnsStr := os.Getenv("DB_MAX_IDLE_CONNS"); dbMaxIdleConnsStr != "" {
		if _, err := strconv.Atoi(dbMaxIdleConnsStr); err != nil {
			return fmt.Errorf("DB_MAX_IDLE_CONNS must be numeric, got: %s", dbMaxIdleConnsStr)
		}
	}

	// Validate ACCOUNTING_ENDPOINT if provided
	if accountingEndpoint := os.Getenv("ACCOUNTING_ENDPOINT"); accountingEndpoint != "" {
		if !regexp.MustCompile(`^https?://`).MatchString(accountingEndpoint) {
			return fmt.Errorf("ACCOUNTING_ENDPOINT must start with http:// or https://")
		}
	}

	return nil
}

// GetEnvInt parses an environment variable as integer with default
func GetEnvInt(key string, defaultVal int) int {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	intVal, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}
	return intVal
}

// GetEnvDuration parses an environment variable as duration with default
func GetEnvDuration(key string, defaultVal time.Duration) time.Duration {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	duration, err := time.ParseDuration(val)
	if err != nil {
		return defaultVal
	}
	return duration
}
