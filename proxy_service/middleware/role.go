package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// RoleMiddleware accepts a list of allowed roles and permits access only if the user's role is allowed.
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
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
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "Role not present in token"})
			c.Abort()
			return
		}

		allowed := false
		for _, r := range allowedRoles {
			if r == role {
				allowed = true
				break
			}
		}

		if !allowed {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient privileges"})
			c.Abort()
			return
		}

		c.Next()
	}
}
