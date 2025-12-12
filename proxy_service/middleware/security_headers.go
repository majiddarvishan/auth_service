package middleware

import "github.com/gin-gonic/gin"

// SecurityHeadersMiddleware adds security headers to all responses
func SecurityHeadersMiddleware(c *gin.Context) {
	// Enforce HTTPS
	c.Header("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")

	// Prevent MIME type sniffing
	c.Header("X-Content-Type-Options", "nosniff")

	// Prevent clickjacking
	c.Header("X-Frame-Options", "DENY")

	// Enable XSS protection in older browsers
	c.Header("X-XSS-Protection", "1; mode=block")

	// Content Security Policy
	c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'")

	// Referrer Policy
	c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

	// Feature Policy
	c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

	c.Next()
}
