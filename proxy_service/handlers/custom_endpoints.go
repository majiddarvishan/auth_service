package handlers

import (
	"net/http"
	"strings"

	"auth_service/database"

	"github.com/gin-gonic/gin"
)

type CreateCustomEndpointRequest struct {
	Path           string `json:"path"`
	Method         string `json:"method"`
	Endpoint       string `json:"endpoint"`
	NeedAccounting bool   `json:"needAccounting"`
}

func CreateCustomEndpointHandler(c *gin.Context) {
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

	// Successfully created endpoint, now re-register dynamic endpoints.
	// Assuming you have a reference to the dynamic router group:
	// For example, if you have a global variable for the dynamic group in main.go:
	// go func() {
	//     // The re-registration can be triggered asynchronously.
	//     // It might be necessary to use a mutex to prevent concurrent modifications.
	//     routes.RegisterDynamicRoutes() // This is a helper that calls RegisterCustomEndpointsDynamic(dynamicGroup)
	// }()

	c.JSON(http.StatusOK, gin.H{"message": "Custom endpoint created successfully", "endpoint": endpoint})

	c.Next()
}
