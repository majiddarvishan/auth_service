package main

import (
    "net/http"
    "github.com/gin-gonic/gin"
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
