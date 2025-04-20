package middleware

import (
	"net/http"
	"strings"

	"auth_service/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware validates the JWT token on protected routes.
func AuthMiddleware(c *gin.Context) {
    authHeader := c.GetHeader("Authorization")
    if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        c.Abort()
        return
    }
    tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

    token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
        // Ensure the signing method is HMAC.
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, jwt.ErrSignatureInvalid
        }
        return []byte(config.SecretKey), nil
    })
    if err != nil || !token.Valid {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
        c.Abort()
        return
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
        c.Abort()
        return
    }

    c.Set("claims", claims)
    c.Next()
}
