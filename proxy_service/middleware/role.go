package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// RoleMiddleware accepts a list of allowed roles and permits access only if the user has any of the allowed roles.
// Supports both single role (legacy "role" claim) and multi-role (new "roles" claim as array).
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

		allowed := false

		// Try multi-role first (new format: "roles" as []interface{})
		if rolesVal, ok := claims["roles"]; ok {
			if rolesArray, ok := rolesVal.([]interface{}); ok {
				for _, userRoleVal := range rolesArray {
					if userRole, ok := userRoleVal.(string); ok {
						for _, allowedRole := range allowedRoles {
							if allowedRole == userRole {
								allowed = true
								break
							}
						}
						if allowed {
							break
						}
					}
				}
			}
		}

		// Fall back to single role (legacy format: "role" as string)
		if !allowed {
			if roleVal, ok := claims["role"].(string); ok {
				for _, allowedRole := range allowedRoles {
					if allowedRole == roleVal {
						allowed = true
						break
					}
				}
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
