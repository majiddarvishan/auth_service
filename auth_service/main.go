package main

import (
    "auth_service/config"
    "auth_service/database"
    "auth_service/routes"
)

func main() {
    // Load configuration from .env file.
    config.LoadConfig()

    // Initialize the database (making your models ready).
    database.InitDB()

    // Set up routes.
    r := routes.SetupRoutes()

    // Run the server on port 8080.
    r.Run(":8080")
}
