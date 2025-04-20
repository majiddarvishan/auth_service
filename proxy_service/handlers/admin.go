package handlers

import (
	"net/http"

	"auth_service/database"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AdminDashboardHandler retrieves user details and system statistics for the admin dashboard.
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

	// Retrieve all users.
	var users []database.User
	// if err := database.DB.Find(&users).Error; err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
	// 	return
	// }
	if err := database.DB.Select("username", "role").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
		return
	}

	// Format response to include only username and role.
	var userData []map[string]string
	for _, user := range users {
		userData = append(userData, map[string]string{
			"username": user.Username,
			"role":     user.Role,
		})
	}

	// Retrieve accounting rules (if applicable).
	var rules []database.AccountingRule
	if err := database.DB.Find(&rules).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve accounting rules"})
		return
	}

	// Format response.
	c.JSON(http.StatusOK, gin.H{
		"message": "Admin dashboard loaded successfully",
		"users":   userData,
		"rules":   rules,
	})
}

// func CreateUserHandler(c *gin.Context) {
//     var req CreateUserRequest
//     if err := c.ShouldBindJSON(&req); err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
//         return
//     }

//     // Optionally, add validation for the username, password, and role.

//     newUser := database.User{
//         Username: req.Username,
//         Password: req.Password, // Remember: always hash the password in production!
//         Role:     req.Role,
//     }

//     if err := database.DB.Create(&newUser).Error; err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
//         return
//     }

//     c.JSON(http.StatusOK, gin.H{"message": "User created successfully", "user": newUser})
// }
