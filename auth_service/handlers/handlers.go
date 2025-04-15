package handlers

import (
	"net/http"
	"time"

	"auth_service/config"
	"auth_service/database"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

// RegisterHandler handles new user registrations.
func RegisterHandler(c *gin.Context) {
	type RegisterRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Role     string `json:"role"` // Optionally allow role specification (make sure to restrict this on production!)
	}

	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	if req.Username == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username and password are required"})
		return
	}

	// Hash the password.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 14)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not hash password"})
		return
	}

	// Use provided role or assign a default role.
	role := req.Role
	if role == "" {
		role = "user"
	}

	user := database.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Role:     role,
	}

	// Create user in the database.
	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

// LoginHandler authenticates the user and returns a JWT token.
func LoginHandler(c *gin.Context) {
	type LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	var user database.User
	// Look up the user by username.
	if err := database.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// Compare the stored hashed password with the incoming password.
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	// Create JWT claims: subject, role, and expiry.
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := jwt.MapClaims{
		"sub":  user.Username,
		"role": user.Role,
		"exp":  expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(config.SecretKey))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenStr})
}

// DeleteUserHandler deletes a user based on the username passed in the URL parameter.
// This endpoint should be accessible only to admins.
func DeleteUserHandler(c *gin.Context) {
    // Get the username from the URL parameter.
    username := c.Param("username")
    if username == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Username is required"})
        return
    }

    // Find the user in the database.
    var user database.User
    if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    // Delete the user.
    if err := database.DB.Delete(&user).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete user", "details": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
