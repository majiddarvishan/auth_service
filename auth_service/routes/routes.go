package routes

import (
    "auth_service/handlers"
    "auth_service/middleware"
    "auth_service/proxy"

    "github.com/gin-gonic/gin"
)

// SetupRoutes configures and returns the Gin engine with all routes defined.
func SetupRoutes() *gin.Engine {
    r := gin.Default()

    // PUBLIC ROUTES:
    // Registration and login endpoints.
    r.POST("/register", handlers.RegisterHandler)
    r.POST("/login", handlers.LoginHandler)

    // DELETE User Endpoint: only admins can delete users.
    r.DELETE("/user/:username",
        middleware.AuthMiddleware,
        middleware.RoleMiddleware("admin"),
        handlers.DeleteUserHandler,
    )

    // Example Admin-only Endpoint.
    r.GET("/admin",
        middleware.AuthMiddleware,
        middleware.RoleMiddleware("admin"),
        func(c *gin.Context) {
            c.JSON(200, gin.H{"message": "Welcome, Admin!"})
        },
    )

    // PROTECTED ROUTES:
    // Create a separate Gin engine for protected routes that require JWT validation.
    protected := gin.New()
    protected.Use(middleware.AuthMiddleware)
    // Catch-all route: forward any unmatched request to your Final-Service.
    protected.Any("/*path", proxy.ProxyRequest)
    // Delegate any unmatched routes from the main router to the protected engine.
    r.NoRoute(func(c *gin.Context) {
        protected.HandleContext(c)
    })

    return r
}
