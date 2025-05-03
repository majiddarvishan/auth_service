package routes

import (
	"auth_service/handlers"
	"auth_service/middleware"
	"log"
	// "fmt"
	// "net"
	// "strings"
	"net/http"

	"github.com/dchest/captcha"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "auth_service/docs"
)

// Explicit CORS middleware for captcha endpoints
func CaptchaCorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

// SetupRoutes configures and returns the Gin engine.
func SetupRoutes(httpAddr, httpsAddr string) {
	httpsRouter := gin.Default()
	// httpRouter := gin.Default()

	// Enable CORS for frontend requests.
	corsConfig := cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60, // 12 hours
	}

	httpsRouter.Use(cors.New(corsConfig))
	// httpRouter.Use(cors.New(corsConfig))

	// Create a separate group for captcha endpoints with explicit CORS
	captchaGroup := httpsRouter.Group("/captcha")
	captchaGroup.Use(CaptchaCorsMiddleware())

	// Handle OPTIONS requests explicitly for captcha endpoints
	httpsRouter.OPTIONS("/captcha/*path", func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		c.AbortWithStatus(http.StatusNoContent)
	})

	// Add captcha endpoints to HTTP router with custom CORS handlers
	captchaGroup.GET("/new", func(c *gin.Context) {
		captchaId := captcha.NewLen(6)
		c.JSON(http.StatusOK, gin.H{"captchaId": captchaId})
	})

	captchaGroup.GET("/image/:captchaId", func(c *gin.Context) {
		captchaId := c.Param("captchaId")
		c.Header("Content-Type", "image/png")
		if err := captcha.WriteImage(c.Writer, captchaId, 240, 80); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
	})

	// Only redirect non-captcha routes to HTTPS
	// httpRouter.NoRoute(func(c *gin.Context) {
	//     if !strings.HasPrefix(c.Request.URL.Path, "/captcha/") {
	//         host := c.Request.Host
	//         if h, _, err := net.SplitHostPort(host); err == nil {
	//             host = h
	//         }
	//         target := fmt.Sprintf("https://%s%s", host, c.Request.RequestURI)
	//         c.Redirect(http.StatusMovedPermanently, target)
	//     } else {
	//         // Handle 404 for captcha routes that don't exist
	//         c.AbortWithStatus(http.StatusNotFound)
	//     }
	// })

	httpsRouter.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// PUBLIC ROUTES:
	httpsRouter.POST("/login", handlers.LoginHandler)

	httpsRouter.GET("/admin",
		middleware.AuthMiddleware,          // Ensure user is authenticated.
		middleware.RoleMiddleware("admin"), // Ensure only admins can access.
		handlers.AdminDashboardHandler,
	)

	// Create a dedicated group for dynamic endpoints.
	var dynamicGroup *gin.RouterGroup = httpsRouter.Group("/") // or some subpath like "/dynamic"

	httpsRouter.POST("/admin/customendpoints",
		middleware.AuthMiddleware,
		middleware.RoleMiddleware("admin"),
		handlers.CreateCustomEndpointHandler(dynamicGroup),
	)

    // httpsRouter.DELETE("/admin/customendpoints/:endpoint",
	// 	middleware.AuthMiddleware,
	// 	middleware.RoleMiddleware("admin"),
	// 	handlers.CreateCustomEndpointHandler(dynamicGroup),
	// )

	// Dynamically register the custom endpoints from the database.
	handlers.RegisterCustomEndpoints(httpsRouter)

	// Add new User Endpoint (Admin Only)
	httpsRouter.POST("/users",
		middleware.AuthMiddleware,          // Ensure the request is authenticated.
		middleware.RoleMiddleware("admin"), // Ensure the requester is an admin.
		handlers.RegisterHandler,           // Handler to create a new user.
	)

	// DELETE User Endpoint (Admin Only)
	httpsRouter.DELETE("/users/:username",
		middleware.AuthMiddleware,
		middleware.RoleMiddleware("admin"),
		handlers.DeleteUserHandler,
	)

	// Update User Role (Admin Only)
	httpsRouter.PUT("/users/:username/role",
		middleware.AuthMiddleware,
		middleware.RoleMiddleware("admin"),
		handlers.UpdateUserRoleHandler,
	)

	// Create New Role (Admin Only)
	httpsRouter.POST("/roles",
		middleware.AuthMiddleware,
		middleware.RoleMiddleware("admin"),
		handlers.CreateRoleHandler,
	)

	httpsRouter.GET("/roles",
		middleware.AuthMiddleware,
		middleware.RoleMiddleware("admin"),
		handlers.GetRolesHandler,
	)

	// Launch both servers
	// go func() {
	//     if err := httpRouter.Run(httpAddr); err != nil {
	//         log.Fatal("HTTP redirection server failed:", err)
	//     }
	// }()
	// if err := httpsRouter.RunTLS(httpsAddr, "cert.pem", "key.pem"); err != nil {
	// if err := httpsRouter.RunTLS(httpsAddr, "localhost.pem", "localhost-key.pem"); err != nil {
	if err := httpsRouter.RunTLS(httpsAddr, "172.26.249.184.pem", "172.26.249.184-key.pem"); err != nil {
		log.Fatal("Failed to start HTTPS server:", err)
	}
}
