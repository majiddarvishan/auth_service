package middleware

import (
	"net/http"
	"sync"
	"time"

	"auth_service/constants"
	"auth_service/types"
	"github.com/gin-gonic/gin"
)

// RateLimitStore tracks request counts by IP
type RateLimitStore struct {
	mu      sync.RWMutex
	buckets map[string]*rateBucket
}

type rateBucket struct {
	count      int
	lastReset  time.Time
	windowSize time.Duration
}

var defaultStore *RateLimitStore

func init() {
	defaultStore = &RateLimitStore{
		buckets: make(map[string]*rateBucket),
	}

	// Cleanup old buckets every minute
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			defaultStore.cleanup()
		}
	}()
}

// cleanup removes expired buckets
func (s *RateLimitStore) cleanup() {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	for ip, bucket := range s.buckets {
		if now.Sub(bucket.lastReset) > bucket.windowSize*2 {
			delete(s.buckets, ip)
		}
	}
}

// checkLimit checks if client has exceeded rate limit
func (s *RateLimitStore) checkLimit(ip string, limit int, window time.Duration) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	bucket, exists := s.buckets[ip]

	if !exists || now.Sub(bucket.lastReset) > window {
		// New window
		s.buckets[ip] = &rateBucket{
			count:      1,
			lastReset:  now,
			windowSize: window,
		}
		return true // Allowed
	}

	if bucket.count >= limit {
		return false // Exceeded
	}

	bucket.count++
	return true // Allowed
}

// RateLimitMiddleware limits requests by IP address
// Default: 100 requests per 60 seconds
func RateLimitMiddleware(limit int, window time.Duration) gin.HandlerFunc {
	if limit == 0 {
		limit = constants.DefaultRateLimitPerSecond
	}
	if window == 0 {
		window = time.Duration(constants.DefaultRateLimitWindow) * time.Second
	}

	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		allowed := defaultStore.checkLimit(clientIP, limit, window)

		if !allowed {
			requestID, _ := c.Get("request_id")
			c.JSON(http.StatusTooManyRequests, types.APIError{
				Code:      constants.ErrorRateLimitExceeded,
				Message:   "rate limit exceeded",
				Timestamp: time.Now().Unix(),
				RequestID: requestID.(string),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
