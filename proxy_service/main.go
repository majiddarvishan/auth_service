package main

import (
	"auth_service/config"
	"auth_service/database"
	"auth_service/routes"
	"flag"
	"fmt"
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
	// _, err := database.NewStore("mock")
	// _, err := database.NewStore("postgres")
     _, err := database.NewStore(*dbMode)
	if err != nil {
		fmt.Println("Error in creating databse connection")
		return
	}

	// Setup routes.
	routes.SetupRoutes(":8080", ":8443")
}
