// package middleware

// import (
//     "net/http"

//     "auth_service/database"
//     "github.com/gin-gonic/gin"
// )

// // ChargeUserMiddleware deducts a charge from the user's balance before allowing access.
// func ChargeUserMiddleware(chargeAmount float64) gin.HandlerFunc {
//     return func(c *gin.Context) {
//         // Get claims from AuthMiddleware
//         claimsVal, exists := c.Get("claims")
//         if !exists {
//             c.JSON(http.StatusUnauthorized, gin.H{"error": "Token claims not found"})
//             c.Abort()
//             return
//         }
//         claims, ok := claimsVal.(map[string]interface{})
//         if !ok {
//             c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
//             c.Abort()
//             return
//         }

//         username, ok := claims["user"].(string)
//         if !ok {
//             c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username in token"})
//             c.Abort()
//             return
//         }

//         // Find user balance
//         var user database.User
//         if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
//             c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
//             c.Abort()
//             return
//         }

//         // Check if balance is sufficient
//         if user.Balance < chargeAmount {
//             c.JSON(http.StatusPaymentRequired, gin.H{"error": "Insufficient balance"})
//             c.Abort()
//             return
//         }

//         // Deduct charge and update balance
//         user.Balance -= chargeAmount
//         database.DB.Save(&user)

//         // Record transaction
//         transaction := database.Transaction{
//             UserID:   user.ID,
//             Amount:   chargeAmount,
//             Endpoint: c.Request.URL.Path,
//         }
//         database.DB.Create(&transaction)

package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"auth_service/config"
	"auth_service/database"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// ChargePayload represents the request we send to accounting service
type ChargePayload struct {
	Username string `json:"username"`
	Endpoint string `json:"endpoint"`
}

// DynamicAccountingMiddleware calls the accounting service to check and deduct a charge.
// If the accounting service returns an error, the request is aborted.
func DynamicAccountingMiddleware(c *gin.Context) {
	// Retrieve JWT claims.
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

	// Extract user ID (stored as uint, not string)
	userIDVal, ok := claims["user"]
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not in token"})
		c.Abort()
		return
	}

	var userID uint
	switch v := userIDVal.(type) {
	case float64:
		userID = uint(v)
	default:
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
		c.Abort()
		return
	}

	// Fetch user by ID to get username
	user, err := database.DB.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		c.Abort()
		return
	}

	// Prepare the payload to send to the accounting service.
	payload := ChargePayload{
		Username: user.Username,
		Endpoint: c.Request.URL.Path,
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not marshal payload"})
		c.Abort()
		return
	}

	// Call the accounting service with timeout.
	accountingURL := config.AccountingEndpoint + "/accounting/charge"
	client := &http.Client{Timeout: 5 * time.Second} // 5 second timeout
	resp, err := client.Post(accountingURL, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		fmt.Printf("Error calling accounting service: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error calling accounting service"})
		c.Abort()
		return
	}
	defer resp.Body.Close()

	// If the accounting service doesn't return 200 OK, abort the request.
	if resp.StatusCode != http.StatusOK {
		c.JSON(resp.StatusCode, gin.H{"error": "Accounting service rejected the charge"})
		c.Abort()
		return
	}

	// If the accounting service returns OK, continue processing.
	c.Next()
}
