package routes

import (
	"auth_service/database"
	"auth_service/handlers"
	"auth_service/middleware"
	"auth_service/proxy"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/dchest/captcha"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "auth_service/docs"
)

// We'll keep a variable for the dynamic endpoints group.
var dynamicGroup *gin.RouterGroup

func RegisterCustomEndpoints(r *gin.Engine) {
	var endpoints []database.CustomEndpoint
	if err := database.DB.Where("enabled = ?", true).Find(&endpoints).Error; err != nil {
		log.Println("Error fetching custom endpoints:", err)
		return
	}

	for _, endpoint := range endpoints {
		// Wrap the handler to include the "endpoint" value.
		wrappedHandler := func(c *gin.Context) {
			proxy.ProxyToEndpoint(c, endpoint.Endpoint)
		}

		// Build the handler chain for the dynamic route.
		handlersChain := []gin.HandlerFunc{middleware.AuthMiddleware}
		if endpoint.NeedAccounting {
			handlersChain = append(handlersChain, middleware.DynamicAccountingMiddleware)
		}
		handlersChain = append(handlersChain, wrappedHandler)

		// Register the endpoint using the specified HTTP method.
		switch endpoint.Method {
		case "GET":
			r.GET(endpoint.Path, handlersChain...)
		case "POST":
			r.POST(endpoint.Path, handlersChain...)
		case "PUT":
			r.PUT(endpoint.Path, handlersChain...)
		case "DELETE":
			r.DELETE(endpoint.Path, handlersChain...)
		default: // ANY or unrecognized method defaults to ANY
			r.Any(endpoint.Path, handlersChain...)
		}
		log.Printf("Registered custom endpoint: %s [%s] -> %s", endpoint.Path, endpoint.Method, endpoint.Endpoint)
	}
}

// RegisterCustomEndpointsDynamic registers dynamic endpoints to the given router group.
func RegisterCustomEndpointsDynamic(group *gin.RouterGroup) {
	// First, clear the group routes if needed.
	// Note: Gin doesn't provide a built-in "Clear()" function; you may need to reinitialize
	// the group or use your own structure to track dynamic routes.
	var endpoints []database.CustomEndpoint
	if err := database.DB.Where("enabled = ?", true).Find(&endpoints).Error; err != nil {
		log.Println("Error fetching custom endpoints:", err)
		return
	}

	for _, ep := range endpoints {
		// Wrap the handler with the endpoint parameter.
		wrappedHandler := func(c *gin.Context) {
			proxy.ProxyToEndpoint(c, ep.Endpoint)
		}

		// Build the handler chain for the dynamic route.
		handlersChain := []gin.HandlerFunc{middleware.AuthMiddleware}
		if ep.NeedAccounting {
			handlersChain = append(handlersChain, middleware.DynamicAccountingMiddleware)
		}
		handlersChain = append(handlersChain, wrappedHandler)

		// Register based on the HTTP method.
		switch ep.Method {
		case "GET":
			group.GET(ep.Path, handlersChain...)
		case "POST":
			group.POST(ep.Path, handlersChain...)
		case "PUT":
			group.PUT(ep.Path, handlersChain...)
		case "DELETE":
			group.DELETE(ep.Path, handlersChain...)
		default:
			group.Any(ep.Path, handlersChain...)
		}
		log.Printf("Registered dynamic route: %s [%s] -> %s", ep.Path, ep.Method, ep.Endpoint)
	}
}

// RegisterDynamicRoutes is a helper to refresh dynamic endpoints.
func RegisterDynamicRoutes(c *gin.Context) {
	// use a mutex or other lock to update dynamicGroup safely.
	RegisterCustomEndpointsDynamic(dynamicGroup)
}

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
    httpRouter := gin.Default()

	// Enable CORS for frontend requests.
    corsConfig := cors.Config{
        AllowOrigins:     []string{"http://localhost:3000"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept", "X-Requested-With"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge:           12 * 60 * 60, // 12 hours
    }

    httpsRouter.Use(cors.New(corsConfig))

    // Create a separate group for captcha endpoints with explicit CORS
    captchaGroup := httpRouter.Group("/captcha")
    captchaGroup.Use(CaptchaCorsMiddleware())

    // Handle OPTIONS requests explicitly for captcha endpoints
    httpRouter.OPTIONS("/captcha/*path", func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
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
    httpRouter.NoRoute(func(c *gin.Context) {
        if !strings.HasPrefix(c.Request.URL.Path, "/captcha/") {
            host := c.Request.Host
            if h, _, err := net.SplitHostPort(host); err == nil {
                host = h
            }
            target := fmt.Sprintf("https://%s%s", host, c.Request.RequestURI)
            c.Redirect(http.StatusMovedPermanently, target)
        } else {
            // Handle 404 for captcha routes that don't exist
            c.AbortWithStatus(http.StatusNotFound)
        }
    })

    // HTTPS Router setup
    httpsRouter.GET("/captcha/new", func(c *gin.Context) {
        captchaId := captcha.NewLen(6)
        c.JSON(http.StatusOK, gin.H{"captchaId": captchaId})
    })

    httpsRouter.GET("/captcha/image/:captchaId", func(c *gin.Context) {
        captchaId := c.Param("captchaId")
        c.Header("Content-Type", "image/png")
        if err := captcha.WriteImage(c.Writer, captchaId, 240, 80); err != nil {
            c.AbortWithStatus(http.StatusInternalServerError)
        }
    })

	httpsRouter.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Create a dedicated group for dynamic endpoints.
	dynamicGroup = httpsRouter.Group("/") // or some subpath like "/dynamic"

	// PUBLIC ROUTES:
	httpsRouter.POST("/login", handlers.LoginHandler)

	httpsRouter.GET("/admin",
		middleware.AuthMiddleware,          // Ensure user is authenticated.
		middleware.RoleMiddleware("admin"), // Ensure only admins can access.
		handlers.AdminDashboardHandler,
	)

	httpsRouter.POST("/admin/customendpoints",
		middleware.AuthMiddleware,
		middleware.RoleMiddleware("admin"),
		handlers.CreateCustomEndpointHandler,
		RegisterDynamicRoutes,
	)

	// Add new User Endpoint (Admin Only)
	httpsRouter.POST("/users",
		middleware.AuthMiddleware,          // Ensure the request is authenticated.
		middleware.RoleMiddleware("admin"), // Ensure the requester is an admin.
		handlers.RegisterHandler,           // Handler to create a new user.
	)

	// DELETE User Endpoint (Admin Only)
	httpsRouter.DELETE("/user/:username",
		middleware.AuthMiddleware,
		middleware.RoleMiddleware("admin"),
		handlers.DeleteUserHandler,
	)

	// Update User Role (Admin Only)
	httpsRouter.PUT("/user/:username/role",
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

	// Dynamically register the custom endpoints from the database.
	RegisterCustomEndpoints(httpsRouter)

    // Launch both servers
    go func() {
        if err := httpRouter.Run(httpAddr); err != nil {
            log.Fatal("HTTP redirection server failed:", err)
        }
    }()
    if err := httpsRouter.RunTLS(httpsAddr, "cert.pem", "key.pem"); err != nil {
        log.Fatal("Failed to start HTTPS server:", err)
    }
}