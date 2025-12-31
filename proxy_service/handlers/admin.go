package handlers

import (
	"net/http"
	"time"

	"auth_service/database"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// SwaggerAccountingRule defines the accounting rule model for Swagger.
// swagger:model SwaggerAccountingRule
//
// Fields:
//   id:         Unique identifier
//   endpoint:   Endpoint the rule applies to
//   charge:     Charge amount for the endpoint
//   created_at: Creation timestamp
//   updated_at: Update timestamp
//
type SwaggerAccountingRule struct {
	ID        uint      `json:"id"`
	Endpoint  string    `json:"endpoint"`
	Charge    float64   `json:"charge"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserRolePair defines username + role in the admin response.
// swagger:model UserRolePair
//
// Fields:
//   username: Username of the user
//   role:     Role assigned to the user
//
type SwaggerUserRolePair struct {
	Username string `json:"username"`
	Role     string `json:"role"`
}

// AdminDashboardResponse is the payload returned by AdminDashboardHandler.
// swagger:model AdminDashboardResponse
//
// Fields:
//   message: Confirmation message
//   users:   List of users with roles
//   rules:   List of accounting rules
//
type AdminDashboardResponse struct {
	Message string                   `json:"message"`
	Users   []SwaggerUserRolePair    `json:"users"`
	Rules   []SwaggerAccountingRule  `json:"rules"`
}

// AdminDashboardHandler retrieves user details and system statistics for the admin dashboard.
// @Summary      Get admin dashboard data
// @Description  Retrieves a list of all users (username + role) and all accounting rules. Admins only.
// @Tags         Admin
// @Produce      json
// @Success      200  {object}  AdminDashboardResponse
// @Failure      401  {object}  map[string]string  "Token claims missing or invalid"
// @Failure      403  {object}  map[string]string  "Not an admin"
// @Failure      500  {object}  map[string]string  "Database error"
// @Security     ApiKeyAuth
// @Router       /admin/dashboard [get]
func AdminDashboardHandler(c *gin.Context) {
	// Ensure only admins can access this route.
	claimsVal, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token claims not found"})
		c.Abort()
		return
	}
	claims, ok := claimsVal.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		c.Abort()
		return
	}

	role, ok := claims["role"].(string)
	if !ok || role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access restricted to admins only"})
		c.Abort()
		return
	}

    users, err := database.DB.GetAllUsers()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err})
        return
    }

	// Format response to include only username and role.
	var userData []map[string]string
	for _, user := range users {
		userData = append(userData, map[string]string{
			"username": user.Username,
			"role":     user.Role.Name,
		})
	}

	// Format response.
	c.JSON(http.StatusOK, gin.H{
		"message": "Admin dashboard loaded successfully",
		"users":   userData,
	})
}

