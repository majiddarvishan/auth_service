package main

import (
	"auth_service/config"
	"auth_service/database"
	"auth_service/logger"
	"auth_service/routes"
	"context"
	"flag"
	"os"
	"os/signal"
	"sync"
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

	// Initialize structured logger
	logger.Init()
	log := logger.Get()

	// Load configuration from .env.
	config.LoadConfig()

	// Validate configuration
	if err := config.ValidateConfig(*dbMode); err != nil {
		log.Error("Configuration validation failed", "error", err.Error())
		os.Exit(1)
	}

	log.Info("Configuration loaded successfully")

	// Initialize the database.
	_, err := database.NewStore(*dbMode)
	if err != nil {
		log.Error("Error creating database connection", "error", err.Error())
		os.Exit(1)
	}

	log.Info("Database connected successfully", "mode", *dbMode)

	// Graceful shutdown coordination
	var wg sync.WaitGroup
	shutdownChan := make(chan struct{})

	// Setup graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		log.Info("Received shutdown signal", "signal", sig.String())
		close(shutdownChan)

		// Give requests 30 seconds to complete
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Wait for in-flight requests
		done := make(chan struct{})
		go func() {
			wg.Wait()
			close(done)
		}()

		select {
		case <-done:
			log.Info("All requests completed")
		case <-ctx.Done():
			log.Warn("Shutdown timeout, forcing exit")
		}

		os.Exit(0)
	}()

	// Setup routes and start server.
	routes.SetupRoutes(":8080", ":8443")
}
