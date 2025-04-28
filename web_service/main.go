package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Message struct {
    From   string `json:"from"`
    To     string `json:"to"`
    Body   string `json:"body"`
    Status string `json:"status"`
}

var messages = make(map[string]Message)

func main() {
    r := gin.Default()

    // Final Service Endpoints
    r.POST("/sms/sendsms", func(c *gin.Context) {
        var msg Message
        if err := c.ShouldBindJSON(&msg); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
            return
        }

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
            SECRET_KEY:="f78973efc0c0664995e2bb055bb2cac6779597a5294685f069229c909358f54a"
            return []byte(SECRET_KEY), nil
            // return []byte(config.SecretKey), nil
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

		user, ok := claims["sub"].(string)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "Username not present in token"})
			c.Abort()
			return
		}
        fmt.Printf("user is %s\n", user)

        // Generate a unique message ID
        messageID := uuid.New().String()
        msg.Status = "Sent"

        // Store the message
        messages[messageID] = msg

        // Respond with the message ID
        c.JSON(http.StatusOK, gin.H{"message-id": messageID})
    })

    r.GET("/getstatus", func(c *gin.Context) {
        messageID := c.Query("message-id")
        if messageID == "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "message-id is required"})
            return
        }

        // Retrieve the message
        msg, exists := messages[messageID]
        if !exists {
            c.JSON(http.StatusNotFound, gin.H{"error": "Message not found"})
            return
        }

        // Respond with the status
        c.JSON(http.StatusOK, gin.H{"status": msg.Status})
    })

    r.GET("/home", func(c *gin.Context) {
        c.String(http.StatusOK, "Welcome to the SMS Service!")
    })

    // Run the Final Service
    r.Run(":8081") // This runs the final service on port 8081
}
