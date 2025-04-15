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
	// Load configuration variables.
	config.LoadConfig()

	// Initialize the database.
	database.InitDB()

	// Create the main Gin router.
	r := gin.Default()

	// PUBLIC ROUTES: register and login.
	r.POST("/register", handlers.RegisterHandler)
	r.POST("/login", handlers.LoginHandler)

    // ----------------------------
    // NEW: DELETE USER ENDPOINT
    // Only an authenticated user with the "admin" role can delete a user.
    r.DELETE("/user/:username", middleware.AuthMiddleware, middleware.RoleMiddleware("admin"), handlers.DeleteUserHandler)
    // ----------------------------

	// Example: an admin-only public endpoint.
	// This route is protected both by the AuthMiddleware and a role check.
	r.GET("/admin", middleware.AuthMiddleware, middleware.RoleMiddleware("admin"), func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Welcome, Admin!"})
	})

	// PROTECTED ROUTES:
	// Create a separate Gin engine for protected routes.
	protected := gin.New()
	protected.Use(middleware.AuthMiddleware)
	// Catch-all route: forward any unmatched URL to the Final-Service.
	protected.Any("/*path", proxy.ProxyRequest)

	// Delegate any request not matched above to the protected engine.
	r.NoRoute(func(c *gin.Context) {
		protected.HandleContext(c)
	})

	// Run the Auth-Service on port 8080.
	r.Run(":8080")
}
