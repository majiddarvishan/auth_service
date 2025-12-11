package config

import (
	"log"
	"os"
	"path/filepath"
	"strings"
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

	// AllowedCORSOrigins is a comma-separated list of allowed CORS origins
	AllowedCORSOrigins []string

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
	if TLSPath == "" {
		log.Fatal("TLS_PATH is not set in .env file")
	}

	// Validate TLS certificate files exist
	certFile := filepath.Join(TLSPath, "localhost.pem")
	keyFile := filepath.Join(TLSPath, "localhost-key.pem")
	if _, err := os.Stat(certFile); err != nil {
		log.Fatal("Certificate file not found at", certFile, ":", err)
	}
	if _, err := os.Stat(keyFile); err != nil {
		log.Fatal("Key file not found at", keyFile, ":", err)
	}

	SecretKey = os.Getenv("SECRET_KEY")
	if SecretKey == "" {
		log.Fatal("SECRET_KEY is not set in .env file")
	}

	p := os.Getenv("TOKEN_EXPIRATION_PERIOD")
	if p == "" {
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

	// Load CORS allowed origins from environment variable
	corsOriginsEnv := os.Getenv("ALLOWED_CORS_ORIGINS")
	if corsOriginsEnv == "" {
		// Default to localhost for development
		AllowedCORSOrigins = []string{"http://localhost:3000", "http://localhost:8080"}
	} else {
		AllowedCORSOrigins = strings.Split(corsOriginsEnv, ",")
	}
}
