package routes

import (
	"auth_service/database"
	"auth_service/handlers"
	"auth_service/middleware"
	"auth_service/proxy"
	"log"

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
	// Note: Gin doesn’t provide a built-in "Clear()" function; you may need to reinitialize
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

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Allow requests from any origin - you can restrict this by replacing "*" with your allowed domain.
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Authorization, Accept, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		// If the method is OPTIONS then return straight away.
		// if c.Request.Method == "OPTIONS" {
		//     c.AbortWithStatus(http.StatusNoContent)
		//     return
		// }

		c.Next()
	}
}

// SetupRoutes configures and returns the Gin engine.
func SetupRoutes() *gin.Engine {
	r := gin.Default()

	// Use the custom CORS middleware
	// r.Use(CORSMiddleware())

	// Enable CORS for frontend requests.
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Allow all origins (or specify "http://localhost:3000")
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Create a dedicated group for dynamic endpoints.
	dynamicGroup = r.Group("/") // or some subpath like "/dynamic"

	// PUBLIC ROUTES:
	r.POST("/login", handlers.LoginHandler)

	r.GET("/admin",
		middleware.AuthMiddleware,          // Ensure user is authenticated.
		middleware.RoleMiddleware("admin"), // Ensure only admins can access.
		handlers.AdminDashboardHandler,
	)

	r.POST("/admin/customendpoints",
		middleware.AuthMiddleware,
		middleware.RoleMiddleware("admin"),
		handlers.CreateCustomEndpointHandler,
		RegisterDynamicRoutes,
	)

	// Add new User Endpoint (Admin Only)
	r.POST("/users",
		middleware.AuthMiddleware,          // Ensure the request is authenticated.
		middleware.RoleMiddleware("admin"), // Ensure the requester is an admin.
		handlers.RegisterHandler,           // Handler to create a new user.
	)

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

	r.GET("/roles",
		middleware.AuthMiddleware,
		middleware.RoleMiddleware("admin"),
		handlers.GetRolesHandler,
	)

	// Dynamically register the custom endpoints from the database.
	RegisterCustomEndpoints(r)

	return r
}
