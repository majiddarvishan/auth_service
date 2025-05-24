package routes

import (
	"auth_service/config"
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
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
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

    // httpsRouter.RedirectTrailingSlash = false
    // httpsRouter.RemoveExtraSlash = true

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
	captchaGroup := httpsRouter.Group(config.BaseApi + "/captcha")
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

    rootGroup := httpsRouter.Group(config.BaseApi)

    // Create a dedicated group for dynamic endpoints.
	var dynamicGroup *gin.RouterGroup = httpsRouter.Group(config.BaseApi)


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

	// redirect /swagger to /swagger/index.html
	// rootGroup.GET("/swagger", func(c *gin.Context) {
    //     c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	// })
	// also catch the slash variant if you like
	// httpsRouter.GET("/swagger/", func(c *gin.Context) {
    //     c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	// })

    // httpsRouter.GET("/swagger", func(c *gin.Context) {
    //     c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
    // })

    rootGroup.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
    rootGroup.GET("/doc", ginSwagger.WrapHandler(swaggerFiles.Handler))

	rootGroup.POST("/login", handlers.LoginHandler)

	rootGroup.GET("/admin",
		middleware.AuthMiddleware,          // Ensure user is authenticated.
		middleware.RoleMiddleware("admin"), // Ensure only admins can access.
		handlers.AdminDashboardHandler,
	)

	rootGroup.POST("/admin/customendpoints",
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
	handlers.RegisterCustomEndpoints(rootGroup)

	// Add new User Endpoint (Admin Only)
	rootGroup.POST("/users",
		middleware.AuthMiddleware,          // Ensure the request is authenticated.
		middleware.RoleMiddleware("admin"), // Ensure the requester is an admin.
		handlers.RegisterHandler,           // Handler to create a new user.
	)

	// DELETE User Endpoint (Admin Only)
	rootGroup.DELETE("/user/:username",
		middleware.AuthMiddleware,
		middleware.RoleMiddleware("admin"),
		handlers.DeleteUserHandler,
	)

	// Update User Role (Admin Only)
	rootGroup.PUT("/user/:username/role",
		middleware.AuthMiddleware,
		middleware.RoleMiddleware("admin"),
		handlers.UpdateUserRoleHandler,
	)

    // Get User Phones (Admin Only)
	// rootGroup.GET("/user/:username/phones",
	// 	middleware.AuthMiddleware,
	// 	middleware.RoleMiddleware("admin"),
	// 	handlers.GetUserPhonesHandler,
	// )

    // rootGroup.POST("/user/:username/phones",
	// 	middleware.AuthMiddleware,
	// 	middleware.RoleMiddleware("admin"),
	// 	handlers.AddUserPhonesHandler,
	// )

	// Create New Role (Admin Only)
	rootGroup.POST("/roles",
		middleware.AuthMiddleware,
		middleware.RoleMiddleware("admin"),
		handlers.CreateRoleHandler,
	)

	rootGroup.GET("/roles",
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
	if err := httpsRouter.RunTLS(httpsAddr, config.TLSPath + "/localhost.pem", config.TLSPath + "/localhost-key.pem"); err != nil {
		log.Fatal("Failed to start HTTPS server:", err)
	}
}
