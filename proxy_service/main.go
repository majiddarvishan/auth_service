package main

import (
	"auth_service/config"
	"auth_service/database"
	"auth_service/routes"
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// @title Auth service API
// @version 1.0
// @description A Auth-service gateway.
// @termsOfService https://example.com/terms

// @contact.name Majid Darvishan
// @contact.url https://github.com/shpd
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /v1/api
func main() {
	dbMode := flag.String("d", "postgres", "Database mode. postgres or mock")
	flag.Parse()

	// Load configuration from .env.
	config.LoadConfig()

	// Initialize the database.
	_, err := database.NewStore(*dbMode)
	if err != nil {
		log.Fatal("Error creating database connection:", err)
	}

	// Setup graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		log.Printf("Received signal: %v. Shutting down gracefully...\n", sig)
		// Give requests 10 seconds to complete
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = ctx
		os.Exit(0)
	}()

	// Setup routes and start server.
	routes.SetupRoutes(":8080", ":8443")
}
