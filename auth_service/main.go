package main

import (
    "auth_service/config"
    "auth_service/database"
    "auth_service/routes"
)

func main() {
    // Load configuration from .env.
    config.LoadConfig()

    // Initialize the database.
    database.InitDB()

    // Setup routes.
    r := routes.SetupRoutes()

    // Run the server on port 8080.
    r.Run(":8080")
}
