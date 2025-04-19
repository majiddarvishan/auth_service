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

//         username, ok := claims["sub"].(string)
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

//         // Continue request handling
//         c.Next()
//     }
// }

package middleware

import (
    "net/http"

    "auth_service/database"
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
)

// DynamicAccountingMiddleware dynamically checks for an accounting rule for the current endpoint.
// If no rule exists for the requested endpoint, it rejects the request with a 403 Forbidden.
// Otherwise, it verifies that the user has sufficient balance, deducts the required charge, and continues.
func DynamicAccountingMiddleware(c *gin.Context) {
    // Query the accounting rule for the requested path.
    var rule database.AccountingRule
    err := database.DB.Where("endpoint = ?", c.Request.URL.Path).First(&rule).Error
    if err != nil {
        // If no accounting rule exists for this endpoint, reject the request.
        c.JSON(http.StatusForbidden, gin.H{"error": "Access forbidden: no accounting rule defined for this endpoint"})
        c.Abort()
        return
    }

    // Retrieve JWT claims for the current user.
    claimsVal, exists := c.Get("claims")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Token claims not found"})
        c.Abort()
        return
    }
    // Cast to jwt.MapClaims instead of map[string]interface{}
    claims, ok := claimsVal.(jwt.MapClaims)
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

    // Find the user in the database.
    var user database.User
    if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        c.Abort()
        return
    }

    // Check if the user has enough balance for the charge.
    if user.Balance < rule.Charge {
        c.JSON(http.StatusPaymentRequired, gin.H{"error": "Insufficient balance"})
        c.Abort()
        return
    }

    // Deduct the charge from the user's balance.
    user.Balance -= rule.Charge
    if err := database.DB.Save(&user).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to update user balance"})
        c.Abort()
        return
    }

    // Continue to the next handler.
    c.Next()
}