package routes

import (
    "auth_service/handlers"
    "auth_service/middleware"
    "auth_service/proxy"

    "github.com/gin-gonic/gin"
)

// SetupRoutes configures and returns the Gin engine.
func SetupRoutes() *gin.Engine {
    r := gin.Default()

    // PUBLIC ROUTES:
    r.POST("/register", handlers.RegisterHandler)
    r.POST("/login", handlers.LoginHandler)

    // DELETE User Endpoint (Admin Only)
    r.DELETE("/user/:username",
        middleware.AuthMiddleware,
        middleware.RoleMiddleware("admin"),
        handlers.DeleteUserHandler,
    )

    // Update User Role (Admin Only)
    r.PUT("/user/:username/role",
        middleware.AuthMiddleware,
        middleware.RoleMiddleware("admin"),
        handlers.UpdateUserRoleHandler,
    )

    // Create New Role (Admin Only)
    r.POST("/roles",
        middleware.AuthMiddleware,
        middleware.RoleMiddleware("admin"),
        handlers.CreateRoleHandler,
    )

    // Premium endpoint example with accounting middleware.
    r.GET("/premium_data",
        middleware.AuthMiddleware,
        middleware.ChargeUserMiddleware(10), // Charge $10 for premium data access
        func(c *gin.Context) {
            c.JSON(200, gin.H{"message": "Premium data accessed"})
        },
    )

    // SMS endpoint: check balance and, if sufficient, forward the request to the final component.
    // Note: We are reusing the generic proxy handler and not hardcoding any SMS logic.
    r.POST("/sms",
        middleware.AuthMiddleware,
        middleware.ChargeUserMiddleware(5), // Charge $5 for sending an SMS
        proxy.ProxyRequest, // Forwards the entire request to config.FinalEndpoint (e.g. http://localhost:8081)
    )

    // Example Admin-only endpoint.
    r.GET("/admin",
        middleware.AuthMiddleware,
        middleware.RoleMiddleware("admin"),
        func(c *gin.Context) {
            c.JSON(200, gin.H{"message": "Welcome, Admin!"})
        },
    )

    // PROTECTED ROUTES:
    // Create a separate engine for protected routes.
    protected := gin.New()
    protected.Use(middleware.AuthMiddleware)
    // Catch-all route: forward any unmatched requests to the final component.
    protected.Any("/*path", proxy.ProxyRequest)
    r.NoRoute(func(c *gin.Context) {
        protected.HandleContext(c)
    })

    return r
}
