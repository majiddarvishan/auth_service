package config

import (
    "log"
    "os"

    "github.com/joho/godotenv"
)

var (
    // DatabaseURL is the connection string to the accounting database.
    DatabaseURL string
    // Port is the port where the accounting service runs.
    Port string
)

func LoadConfig() {
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, continuing with system env")
    }

    DatabaseURL = os.Getenv("DATABASE_URL")
    if DatabaseURL == "" {
        log.Fatal("DATABASE_URL is not set")
    }

    Port = os.Getenv("PORT")
    if Port == "" {
        // Default port if not set
        Port = "8082"
    }
}
