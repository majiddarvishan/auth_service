package config

import (
    "log"
    "os"

    "github.com/joho/godotenv"
)

var (
    SecretKey   string
    DatabaseURL string
)

func LoadConfig() {
    // Load environment variables from .env file
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
