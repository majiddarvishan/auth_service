package handlers

import (
	"log"
	"net/http"
	"strings"

	"auth_service/database"
	"auth_service/middleware"
	"auth_service/proxy"

	"github.com/gin-gonic/gin"
)

type CreateCustomEndpointRequest struct {
	Path           string `json:"path"`
	Method         string `json:"method"`
	Endpoint       string `json:"endpoint"`
	NeedAccounting bool   `json:"needAccounting"`
}

func registerCustomEndpointDynamic(r *gin.RouterGroup, ep *database.CustomEndpoint) {
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
		r.GET(ep.Path, handlersChain...)
	case "POST":
		r.POST(ep.Path, handlersChain...)
	case "PUT":
		r.PUT(ep.Path, handlersChain...)
	case "DELETE":
		r.DELETE(ep.Path, handlersChain...)
	default:
		r.Any(ep.Path, handlersChain...)
	}
	log.Printf("Registered dynamic route: %s [%s] -> %s", ep.Path, ep.Method, ep.Endpoint)
}

func RegisterCustomEndpoints(r *gin.Engine) {
	var endpoints []database.CustomEndpoint
	if err := database.DB.Where("enabled = ?", true).Find(&endpoints).Error; err != nil {
		log.Println("Error fetching custom endpoints:", err)
		return
	}

	for _, endpoint := range endpoints {
		registerCustomEndpointDynamic(&r.RouterGroup, &endpoint)
	}
}

func CreateCustomEndpointHandler(dynamicGroup *gin.RouterGroup) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateCustomEndpointRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
			return
		}

		// Validate endpoint format
		if req.Endpoint == "" || !strings.HasPrefix(req.Endpoint, "http") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or missing endpoint URL"})
			return
		}

		if req.Method == "" {
			req.Method = "ANY"
		}

		endpoint := database.CustomEndpoint{
			Path:           req.Path,
			Method:         req.Method,
			Endpoint:       req.Endpoint,
			NeedAccounting: req.NeedAccounting,
			Enabled:        true,
		}

		if err := database.DB.Create(&endpoint).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create custom endpoint"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Custom endpoint created successfully", "endpoint": endpoint})

		registerCustomEndpointDynamic(dynamicGroup, &endpoint)

		c.Next()
	}
}
