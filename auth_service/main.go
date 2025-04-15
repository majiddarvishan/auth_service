package main

import (
    "auth_service/config"
    "auth_service/database"
    "auth_service/handlers"
    "auth_service/middleware"
    "auth_service/proxy"

    "github.com/gin-gonic/gin"
)

func main() {
    // Load configuration from .env
    config.LoadConfig()

    // Initialize the database connection
    database.InitDB()

    // Create the main Gin router.
    r := gin.Default()

    // PUBLIC ROUTES: Register and Login endpoints.
    r.POST("/register", handlers.RegisterHandler)
    r.POST("/login", handlers.LoginHandler)

    // PROTECTED ROUTES:
    // Create a separate Gin engine for protected routes,
    // so we avoid conflicts with the public routes.
    protected := gin.New()
    protected.Use(middleware.AuthMiddleware)
    // Use a catch-all route to forward any unmatched request to the Final-Service.
    protected.Any("/*path", proxy.ProxyToFinalService)

    // For any request not caught by public routes,
    // use the NoRoute handler to delegate to the protected engine.
    r.NoRoute(func(c *gin.Context) {
        protected.HandleContext(c)
    })

    // Run the Auth-Service on port 8080
    r.Run(":8080")
}
