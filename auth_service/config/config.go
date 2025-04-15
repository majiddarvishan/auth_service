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
}
