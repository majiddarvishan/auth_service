package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	// SecretKey is used for signing JWT tokens.
	SecretKey string

	// DatabaseURL is the DSN for connecting to your PostgreSQL (or MySQL) database.
	DatabaseURL string

	// AccountingEndpoint is the URL for the Accounting-Service
	AccountingEndpoint string

	// FinalEndpoint is the URL for the Final-Service
	SmsEndpoint string
)

// LoadConfig loads environment variables from a .env file.
func LoadConfig() {
	// Load .env file (log fatal if not found).
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	SecretKey = os.Getenv("SECRET_KEY")
	if SecretKey == "" {
		log.Fatal("SECRET_KEY is not set in .env file")
	}

	DatabaseURL = os.Getenv("DATABASE_URL")
	if DatabaseURL == "" {
		log.Fatal("DATABASE_URL is not set in .env file")
	}

	AccountingEndpoint = os.Getenv("ACCOUNTING_ENDPOINT")
	if AccountingEndpoint == "" {
		// Default to local accounting port.
		AccountingEndpoint = "http://localhost:8082"
	}

	// Load the final endpoint URL.
	SmsEndpoint = os.Getenv("SMS_ENDPOINT")
	if SmsEndpoint == "" {
		// Provide a default if not set, or log fatal if you require it.
		SmsEndpoint = "http://localhost:8081"
	}
}
