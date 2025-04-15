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

    // FinalEndpoint is the URL for the Final-Service
    FinalEndpoint string
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

    // Load the final endpoint URL.
    FinalEndpoint = os.Getenv("FINAL_ENDPOINT")
    if FinalEndpoint == "" {
        // Provide a default if not set, or log fatal if you require it.
        FinalEndpoint = "http://localhost:8081"
    }
}
