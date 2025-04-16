package routes

import (
    "auth_service/handlers"
    "auth_service/middleware"
    "auth_service/proxy"

    "github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
    r := gin.Default()

    // PUBLIC ROUTES:
    r.POST("/register", handlers.RegisterHandler)
    r.POST("/login", handlers.LoginHandler)
    r.DELETE("/user/:username",
        middleware.AuthMiddleware,
        middleware.RoleMiddleware("admin"),
        handlers.DeleteUserHandler,
    )
    r.PUT("/user/:username/role",
        middleware.AuthMiddleware,
        middleware.RoleMiddleware("admin"),
        handlers.UpdateUserRoleHandler,
    )

    // NEW: Route for creating new roles (admin only):
    r.POST("/roles",
        middleware.AuthMiddleware,
        middleware.RoleMiddleware("admin"),
        handlers.CreateRoleHandler,
    )

    // Example: An admin-only endpoint.
    r.GET("/admin",
        middleware.AuthMiddleware,
        middleware.RoleMiddleware("admin"),
        func(c *gin.Context) {
            c.JSON(200, gin.H{"message": "Welcome, Admin!"})
        },
    )

    // PROTECTED ROUTES:
    protected := gin.New()
    protected.Use(middleware.AuthMiddleware)
    // Catch-all route: forward any unmatched requests to the Final-Service.
    protected.Any("/*path", proxy.ProxyRequest)
    r.NoRoute(func(c *gin.Context) {
        protected.HandleContext(c)
    })

    return r
}
