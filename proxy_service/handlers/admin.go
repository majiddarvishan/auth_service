package handlers

import (
	"net/http"
	"strings"

	"auth_service/database"
	// "auth_service/routes"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// type CreateUserRequest struct {
//     Username string `json:"username"`
//     Password string `json:"password"`
//     Role     string `json:"role"`
// }

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

type CreateCustomEndpointRequest struct {
    Path        string `json:"path"`
    HandlerName string `json:"handler"`  // Maps to a registered handler
    Method      string `json:"method"`   // Optional: "GET", "POST", etc. Defaults to "ANY".
    Endpoint    string `json:"endpoint"` // Required: Target endpoint URL.
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
        Path:        req.Path,
        HandlerName: req.HandlerName,
        Method:      req.Method,
        Endpoint:    req.Endpoint,
        Enabled:     true,
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
