package main

import (
    "auth_service/config"
    "auth_service/database"
    "auth_service/routes"

    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
)

func main() {
    // Load configuration from .env.
    config.LoadConfig()

    // Initialize the database.
    database.InitDB()

    // Initialize Gin router.
    r := gin.Default()

    // Enable CORS for frontend requests.
    r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"*"}, // Allow all origins (or specify "http://localhost:3000")
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Authorization", "Content-Type", "Accept"},
        AllowCredentials: true,
    }))

    // Setup routes.
    routes.SetupRoutes(r)

    // Run the server on port 8080.
    r.Run(":8080")
}
