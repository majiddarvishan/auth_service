package middleware

import (
	"auth_service/constants"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestIDMiddleware adds a unique request ID to each request for tracing
func RequestIDMiddleware(c *gin.Context) {
	requestID := c.GetHeader(constants.HeaderRequestID)
	if requestID == "" {
		requestID = uuid.New().String()
	}

	c.Set("request_id", requestID)
	c.Header(constants.HeaderRequestID, requestID)
	c.Next()
}
