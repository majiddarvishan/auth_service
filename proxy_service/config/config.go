package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

var (
	BaseApi string

	// TLSPath is used for read TLS files from that path.
	TLSPath string

	// SecretKey is used for signing JWT tokens.
	SecretKey string

	// tokenExpirationPeriod is the duration for which the JWT token is valid.
	TokenExpirationPeriod time.Duration

	// AccountingEndpoint is the URL for the Accounting-Service
	AccountingEndpoint string

	DatabaseHost string

	DatabasePort string

	DatabaseUserName string

	DatabasePassword string

	DatabaseName string
)

// LoadConfig loads environment variables from a .env file.
func LoadConfig() {
	// Attempt to load .env only if it exists.
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found (or not needed), continuing with environment variables.")
	}

	BaseApi = os.Getenv("BASE_API")
	if BaseApi == "" {
		log.Fatal("BASE_API is not set in .env file")
	}

	TLSPath = os.Getenv("TLS_PATH")
	if BaseApi == "" {
		log.Fatal("TLS_PATH is not set in .env file")
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

	AccountingEndpoint = os.Getenv("ACCOUNTING_ENDPOINT")
	if AccountingEndpoint == "" {
		// Default to local accounting port.
		AccountingEndpoint = "http://localhost:8082"
	}

	DatabaseHost = os.Getenv("DB_HOST")
	if DatabaseHost == "" {
		log.Fatal("DB_HOST is not set in .env file")
	}

	DatabasePort = os.Getenv("DB_PORT")
	if DatabasePort == "" {
		log.Fatal("DB_PORT is not set in .env file")
	}

	DatabaseUserName = os.Getenv("DB_USER_NAME")
	if DatabaseUserName == "" {
		log.Fatal("DB_USER_NAME is not set in .env file")
	}

	DatabasePassword = os.Getenv("DB_PASSWORD")
	if DatabasePassword == "" {
		log.Fatal("DB_PASSWORD is not set in .env file")
	}

	DatabaseName = os.Getenv("DB_NAME")
	if DatabaseName == "" {
		log.Fatal("DB_NAME is not set in .env file")
	}
}
