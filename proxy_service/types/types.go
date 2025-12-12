package types

import (
	"net/http"
	"time"
)

// APIError represents a standardized error response
type APIError struct {
	Code      string `json:"code"`                 // Machine-readable: USER_NOT_FOUND
	Message   string `json:"message"`              // Human-readable message
	Details   string `json:"details,omitempty"`    // Optional debug info (not in prod)
	Timestamp int64  `json:"timestamp"`            // Unix timestamp
	RequestID string `json:"request_id,omitempty"` // For tracing
	Path      string `json:"path,omitempty"`       // Request path
}

// APISuccess represents a standardized success response
type APISuccess struct {
	Data      interface{} `json:"data,omitempty"`
	Message   string      `json:"message,omitempty"`
	Timestamp int64       `json:"timestamp"`
	RequestID string      `json:"request_id,omitempty"`
}

// PaginationRequest represents pagination query parameters
type PaginationRequest struct {
	Page  int `query:"page" binding:"min=1"`          // Default 1
	Limit int `query:"limit" binding:"min=1,max=100"` // Default 10, max 100
}

// PaginatedResponse represents a paginated response
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalPages int64       `json:"total_pages"`
	Timestamp  int64       `json:"timestamp"`
}

// TokenResponse represents login/refresh token response
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	ExpiresIn    int    `json:"expires_in"` // seconds
	TokenType    string `json:"token_type"` // "Bearer"
	Timestamp    int64  `json:"timestamp"`
}

// HealthStatus represents the health check response
type HealthStatus struct {
	Status    string `json:"status"` // "healthy" or "unhealthy"
	Timestamp int64  `json:"timestamp"`
	Database  string `json:"database,omitempty"` // "connected" or error message
	Version   string `json:"version,omitempty"`
}

// RequestMetadata contains request context
type RequestMetadata struct {
	RequestID string
	Timestamp time.Time
	UserID    uint
	Username  string
	Role      string
	Path      string
	Method    string
}

// ErrorResponse returns error status code for given error code
func ErrorStatusCode(errorCode string) int {
	statusMap := map[string]int{
		"USER_NOT_FOUND":            http.StatusNotFound,
		"INVALID_PASSWORD":          http.StatusUnauthorized,
		"USER_ALREADY_EXISTS":       http.StatusConflict,
		"INVALID_USERNAME":          http.StatusBadRequest,
		"INVALID_PASSWORD_STRENGTH": http.StatusBadRequest,
		"UNAUTHORIZED":              http.StatusUnauthorized,
		"FORBIDDEN":                 http.StatusForbidden,
		"ROLE_NOT_FOUND":            http.StatusNotFound,
		"ENDPOINT_NOT_FOUND":        http.StatusNotFound,
		"INVALID_JSON":              http.StatusBadRequest,
		"INVALID_PATH":              http.StatusBadRequest,
		"RATE_LIMIT_EXCEEDED":       http.StatusTooManyRequests,
		"INVALID_REFRESH_TOKEN":     http.StatusUnauthorized,
		"REFRESH_TOKEN_EXPIRED":     http.StatusUnauthorized,
	}
	if status, ok := statusMap[errorCode]; ok {
		return status
	}
	return http.StatusInternalServerError
}
