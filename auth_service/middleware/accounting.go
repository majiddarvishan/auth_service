package middleware

import (
    "net/http"

    "auth_service/database"
    "github.com/gin-gonic/gin"
)

// ChargeUserMiddleware deducts a charge from the user's balance before allowing access.
func ChargeUserMiddleware(chargeAmount float64) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Get claims from AuthMiddleware
        claimsVal, exists := c.Get("claims")
        if !exists {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Token claims not found"})
            c.Abort()
            return
        }
        claims, ok := claimsVal.(map[string]interface{})
        if !ok {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
            c.Abort()
            return
        }

        username, ok := claims["sub"].(string)
        if !ok {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username in token"})
            c.Abort()
            return
        }

        // Find user balance
        var user database.User
        if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
            c.Abort()
            return
        }

        // Check if balance is sufficient
        if user.Balance < chargeAmount {
            c.JSON(http.StatusPaymentRequired, gin.H{"error": "Insufficient balance"})
            c.Abort()
            return
        }

        // Deduct charge and update balance
        user.Balance -= chargeAmount
        database.DB.Save(&user)

        // Record transaction
        transaction := database.Transaction{
            UserID:   user.ID,
            Amount:   chargeAmount,
            Endpoint: c.Request.URL.Path,
        }
        database.DB.Create(&transaction)

        // Continue request handling
        c.Next()
    }
}
