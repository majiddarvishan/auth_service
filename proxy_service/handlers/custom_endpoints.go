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

func registerCustomEndpointDynamic(r *gin.RouterGroup, ep *database.CustomEndpoint) {
	// Wrap the handler with the endpoint parameter.
	wrappedHandler := func(c *gin.Context) {
		proxy.ProxyToEndpoint(c, ep.Endpoints)
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
	log.Printf("Registered dynamic route: %s [%s] -> %s", ep.Path, ep.Method, ep.Endpoints[0])
}

// func RegisterCustomEndpoints(r *gin.Engine) {
func RegisterCustomEndpoints(routerGroup *gin.RouterGroup) {
	endpoints, err := database.DB.GetAllCustomEndpoints()
	if err != nil {
		log.Println("Error fetching custom endpoints:", err)
		return
	}

	for _, endpoint := range endpoints {
		registerCustomEndpointDynamic(routerGroup, &endpoint)
	}
}

// SwaggerCustomEndpoint represents the payload for Custom Endpoint.
// swagger:model SwaggerCustomEndpoint
// @Description SwaggerCustomEndpoint defines the expected request body for custom endpoint.
// @Property path body string true "Prefix of URI"
// @Property method body string true "method of request"
// @Property endpoints body []string true  "ID of the captcha challenge"
// @Property needAccounting body string true  "Needs check accouting before redirect it"
type SwaggerCustomEndpoint struct {
	Path           string
	Method         string
	Endpoints      []string
	NeedAccounting bool
}

// CreateCustomEndpointHandler create custom endpoint.
// @Summary      Create Custom Endpoint
// @Description  Create Custom Endpoint to redirect its requests to another endpoints
// @Tags         CustomEndpoints
// @Accept       json
// @Produce      json
// @Param        request  body      SwaggerCustomEndpoint  true  "CustomEndpoint payload"
// @Success      200      {object}  map[string]string  "JWT token"
// @Failure      400      {object}  map[string]string  "Invalid JSON format"
// @Failure      401      {object}  map[string]string  "Unauthorized: invalid credentials"
// @Failure      500      {object}  map[string]string  "Server error during token generation"
// @Router       /admin/custom-endpoints [post]
func CreateCustomEndpointHandler(dynamicGroup *gin.RouterGroup) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req database.CustomEndpoint
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
			return
		}

		// Validate endpoints format
		for _, endpoint := range req.Endpoints {
			if endpoint == "" || !strings.HasPrefix(endpoint, "http") {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or missing endpoint URL"})
				return
			}
		}

		if req.Method == "" {
			req.Method = "ANY"
		}

		req.Path += "/*path"
		req.Enabled = true

		if err := database.DB.CreateCustomEndpoint(&req); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create custom endpoint"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Custom endpoint created successfully", "endpoint": req})

		registerCustomEndpointDynamic(dynamicGroup, &req)

		c.Next()
	}
}

// DeleteCustomEndpointHandler deletes a custom endpoint dynamically.
//
// @Summary      Delete a custom endpoint
// @Description  Deletes a previously registered custom endpoint by path and restores a 404 handler for it.
// @Tags         CustomEndpoints
// @Accept       json
// @Produce      json
// @Param        request body SwaggerCustomEndpoint true "Custom endpoint object (must include path)"
// @Success      200 {object} map[string]interface{} "Custom endpoint deleted successfully"
// @Failure      400 {object} map[string]interface{} "Invalid JSON payload"
// @Failure      500 {object} map[string]interface{} "Failed to delete or find custom endpoint"
// @Router       /admin/custom-endpoints [delete]
func DeleteCustomEndpointHandler(dynamicGroup *gin.RouterGroup) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req database.CustomEndpoint
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
			return
		}

		req.Path += "/*path"
		endPoint, err := database.DB.GetCustomEndpointByPath(req.Path)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find custom endpoint", "endpoint": req.Path})
			return
		}

		if err = database.DB.DeleteCustomEndpointByPath(req.Path); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete custom endpoint", "endpoint": req.Path})
			return
		}

		notFoundHandler := func(c *gin.Context) {
			c.JSON(404, gin.H{"error": "Not found"})
			c.Abort() // Important: stop the chain
		}

		switch endPoint.Method {
		case "GET":
			dynamicGroup.GET(endPoint.Path, notFoundHandler)
		case "POST":
			dynamicGroup.POST(endPoint.Path, notFoundHandler)
		case "PUT":
			dynamicGroup.PUT(endPoint.Path, notFoundHandler)
		case "DELETE":
			dynamicGroup.DELETE(endPoint.Path, notFoundHandler)
		default:
			dynamicGroup.Any(endPoint.Path, notFoundHandler)
		}

		c.JSON(http.StatusOK, gin.H{"message": "Custom endpoint deleted successfully", "endpoint": req.Path})

		c.Next()
	}
}
