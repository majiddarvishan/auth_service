package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

var (
	// SecretKey is used for signing JWT tokens.
	SecretKey string

    // tokenExpirationPeriod is the duration for which the JWT token is valid.
    TokenExpirationPeriod time.Duration

	// DatabaseURL is the DSN for connecting to your PostgreSQL (or MySQL) database.
	DatabaseURL string

	// AccountingEndpoint is the URL for the Accounting-Service
	AccountingEndpoint string
)

// LoadConfig loads environment variables from a .env file.
func LoadConfig() {
    // Attempt to load .env only if it exists.
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found (or not needed), continuing with environment variables.")
    }

	SecretKey = os.Getenv("SECRET_KEY")
	if SecretKey == "" {
		log.Fatal("SECRET_KEY is not set in .env file")
	}

    p := os.Getenv("TOKEN_EXPIRATION_PERIOD")
	if SecretKey == "" {
		log.Fatal("TOKEN_EXPIRATION_PERIOD is not set in .env file")
	}

    var err error
    TokenExpirationPeriod, err = time.ParseDuration(p)
    if err != nil {
        log.Fatalf("Error parsing %s: %v\n", p, err)
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
}
