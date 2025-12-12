package handlers

import (
	"net/http"
	"time"

	"auth_service/database"
	"auth_service/types"
	"github.com/gin-gonic/gin"
)

// HealthHandler provides health check endpoint for load balancers
// @Summary      Health Check
// @Description  Check service health and database connectivity
// @Tags         Health
// @Produce      json
// @Success      200      {object}  types.HealthStatus  "Service is healthy"
// @Failure      503      {object}  types.HealthStatus  "Service is unhealthy"
// @Router       /health [get]
func HealthHandler(c *gin.Context) {
	health := types.HealthStatus{
		Status:    "healthy",
		Timestamp: time.Now().Unix(),
		Database:  "connected",
	}

	// Check database connectivity by pinging it
	if database.DB != nil {
		// Try a simple query - GetAllUsers as a connectivity check
		_, err := database.DB.GetAllUsers()
		if err != nil {
			health.Status = "unhealthy"
			health.Database = "disconnected: " + err.Error()
			c.JSON(http.StatusServiceUnavailable, health)
			return
		}
	} else {
		health.Status = "unhealthy"
		health.Database = "not initialized"
		c.JSON(http.StatusServiceUnavailable, health)
		return
	}

	c.JSON(http.StatusOK, health)
}

// VersionHandler returns service version
// @Summary      Get Version
// @Description  Get service version information
// @Tags         Info
// @Produce      json
// @Success      200  {object}  map[string]string  "Version information"
// @Router       /version [get]
func VersionHandler(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"service":   "auth-service",
		"version":   "1.0.0",
		"timestamp": time.Now().Unix(),
	})
}
