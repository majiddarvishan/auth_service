# Improvement Ideas for Auth Service Proxy

## 1. Input Validation & Security

### 1.1 Password Complexity (Beyond Length)
**Current**: Only checks minimum 8 characters
**Improvement**: Add uppercase, lowercase, digit, special character requirements
```go
// In handlers/user.go RegisterHandler
func ValidatePasswordStrength(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}
	hasUpper, hasLower, hasDigit, hasSpecial := false, false, false, false
	for _, char := range password {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= '0' && char <= '9':
			hasDigit = true
		case strings.ContainsRune("!@#$%^&*()_+-=[]{}|;:,.<>?", char):
			hasSpecial = true
		}
	}
	if !hasUpper || !hasLower || !hasDigit || !hasSpecial {
		return fmt.Errorf("password must contain uppercase, lowercase, digit, and special character")
	}
	return nil
}
```

### 1.2 Username Sanitization
**Current**: Only checks minimum 3 characters
**Improvement**: Validate alphanumeric + underscore only, prevent SQL injection patterns
```go
// Add to handlers/user.go
func ValidateUsername(username string) error {
	if len(username) < 3 || len(username) > 32 {
		return fmt.Errorf("username must be 3-32 characters")
	}
	if !regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(username) {
		return fmt.Errorf("username can only contain letters, numbers, and underscores")
	}
	return nil
}
```

### 1.3 Endpoint Path Validation
**Current**: No validation on custom endpoint paths
**Improvement**: Prevent malicious patterns in dynamic route registration
```go
// In handlers/custom_endpoints.go CreateCustomEndpointHandler
func ValidateEndpointPath(path string) error {
	if !strings.HasPrefix(path, "/") {
		return fmt.Errorf("path must start with /")
	}
	if strings.Contains(path, "..") {
		return fmt.Errorf("path traversal not allowed")
	}
	if strings.Contains(path, "//") {
		return fmt.Errorf("double slashes not allowed")
	}
	return nil
}
```

---

## 2. Error Handling & Logging

### 2.1 Structured Logging
**Current**: Uses standard `log` package with inconsistent formats
**Improvement**: Add structured logging with levels
```go
// Add to main.go or new package logger/logger.go
import "log/slog"

var Logger *slog.Logger

func init() {
	Logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
}

// Usage in handlers:
Logger.Info("user registered", "username", req.Username, "timestamp", time.Now())
Logger.Error("database error", "error", err, "operation", "create_user")
```
**Benefits**:
- Consistent JSON format for log aggregation
- Searchable levels (info, warn, error, debug)
- Better for production monitoring

### 2.2 API Error Response Standard
**Current**: Inconsistent error responses (some include "details", some don't)
**Improvement**: Standardize all error responses
```go
// Create types/errors.go
type APIError struct {
	Code      string `json:"code"`      // Machine-readable: USER_NOT_FOUND, INVALID_PASSWORD
	Message   string `json:"message"`   // Human-readable message
	Details   string `json:"details,omitempty"`   // Optional debug info
	Timestamp int64  `json:"timestamp"`
	RequestID string `json:"request_id,omitempty"`
}

// Usage helper:
func RespondError(c *gin.Context, status int, code, message string) {
	c.JSON(status, APIError{
		Code:      code,
		Message:   message,
		Timestamp: time.Now().Unix(),
	})
}
```
**Usage**:
```go
RespondError(c, http.StatusUnauthorized, "INVALID_PASSWORD", "Password does not match")
```

### 2.3 Request ID Tracking
**Current**: No correlation between logs and requests
**Improvement**: Add request ID middleware for tracing
```go
// In middleware/request_id.go
func RequestIDMiddleware(c *gin.Context) {
	requestID := c.GetHeader("X-Request-ID")
	if requestID == "" {
		requestID = uuid.New().String()
	}
	c.Set("request_id", requestID)
	c.Header("X-Request-ID", requestID)
	c.Next()
}

// Use in handlers:
requestID := c.GetString("request_id")
Logger.Info("operation completed", "request_id", requestID)
```

---

## 3. Database & Performance

### 3.1 Connection Pooling Configuration
**Current**: GORM default pool settings
**Improvement**: Make pool settings configurable
```go
// In config/config.go
var (
	DBMaxOpenConns int
	DBMaxIdleConns int
	DBConnMaxLifetime time.Duration
)

// In config loading:
DBMaxOpenConns = getEnvInt("DB_MAX_OPEN_CONNS", 25)
DBMaxIdleConns = getEnvInt("DB_MAX_IDLE_CONNS", 5)

// In database/pgstore.go Init():
sqlDB, _ := db.DB()
sqlDB.SetMaxOpenConns(config.DBMaxOpenConns)
sqlDB.SetMaxIdleConns(config.DBMaxIdleConns)
sqlDB.SetConnMaxLifetime(config.DBConnMaxLifetime)
```

### 3.2 Query Optimization - N+1 Prevention
**Current**: `GetUserAndRoleByUsername` might load role inefficiently
**Improvement**: Use GORM Preload
```go
// In database/pgstore.go
func (s *PGStore) GetUserAndRoleByUsername(username string) (*User, error) {
	var user User
	// GORM preload prevents N+1 query
	if err := s.db.Preload("Role").Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
```

### 3.3 Soft Deletes for Audit Trail
**Current**: Hard deletes lose historical data
**Improvement**: Use GORM soft deletes for audit
```go
// In database/models.go
type User struct {
	gorm.Model // Includes DeletedAt field
	// ... fields ...
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// GORM automatically filters deleted records in queries
// To view deleted: db.Unscoped().Where("deleted_at IS NOT NULL")...
```

---

## 4. API Design & Documentation

### 4.1 Pagination for List Endpoints
**Current**: `GetAllUsers()`, `GetAllCustomEndpoints()` return everything
**Improvement**: Add pagination to prevent memory exhaustion
```go
// Create types/pagination.go
type PaginationRequest struct {
	Page  int `query:"page" binding:"min=1"` // Default 1
	Limit int `query:"limit" binding:"min=1,max=100"` // Default 10, max 100
}

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalPages int         `json:"total_pages"`
}

// Usage in handlers:
var users []User
db.Offset((req.Page-1)*req.Limit).Limit(req.Limit).Find(&users)
```

### 4.2 Consistent Query Parameter Naming
**Current**: Mix of patterns in endpoint handling
**Improvement**: Standardize (filter, sort, search patterns)
```go
// Standard patterns:
// GET /users?search=john&role=admin&sort=-created_at&page=1&limit=20
// GET /custom-endpoints?enabled=true&sort=path&page=1&limit=50
```

### 4.3 API Versioning in Code
**Current**: All routes under `/v1/api`
**Improvement**: Make versioning consistent and documented
```go
// In routes/routes.go - document version strategy
// v1 endpoints deprecated in favor of v2
// Transition plan: Deprecate v1 by 2025-12-31

v1Group := httpsRouter.Group("/v1/api")
v2Group := httpsRouter.Group("/v2/api") // New endpoints here
```

---

## 5. Testing & Quality

### 5.1 Integration Tests for Custom Endpoints
**Current**: Basic unit tests only
**Improvement**: Test proxy behavior with multiple targets
```go
// handlers/custom_endpoints_test.go
func TestProxyToMultipleTargets(t *testing.T) {
	// Setup two mock backend servers
	backend1 := httptest.NewServer(http.HandlerFunc(...))
	backend2 := httptest.NewServer(http.HandlerFunc(...))

	// Create custom endpoint pointing to both
	// Test load balancing works across requests
	// Test fallback if one backend down
}
```

### 5.2 Rate Limiting Implementation
**Current**: `RateLimitMiddleware()` is a placeholder
**Improvement**: Implement actual rate limiting
```go
// middleware/rate_limit.go
import "github.com/ulule/limiter/v3"

func RateLimitMiddleware(limit int, period time.Duration) gin.HandlerFunc {
	store := memory.NewStore()
	limiter := limiter.New(store, limiter.Rate{
		Limit:  int64(limit),
		Period: period,
	})

	return func(c *gin.Context) {
		context, err := limiter.Get(c.Request.Context(), c.ClientIP())
		if err != nil || context.Reached {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded",
				"retry_after": context.ResetAfter,
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
```

---

## 6. Security Enhancements

### 6.1 JWT Token Rotation / Refresh Tokens
**Current**: Single JWT with fixed expiry
**Improvement**: Add refresh token pattern
```go
// handlers/user.go - LoginHandler response
c.JSON(http.StatusOK, gin.H{
	"access_token": accessToken,    // Short-lived (15 min)
	"refresh_token": refreshToken,  // Long-lived (7 days)
	"expires_in": 900,              // seconds
})

// New POST /refresh endpoint
// Exchange refresh_token for new access_token
```

### 6.2 HTTPS Only with HSTS
**Current**: TLS configured but HSTS header missing
**Improvement**: Add security headers
```go
// middleware/security_headers.go
func SecurityHeadersMiddleware(c *gin.Context) {
	c.Header("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
	c.Header("X-Content-Type-Options", "nosniff")
	c.Header("X-Frame-Options", "DENY")
	c.Header("X-XSS-Protection", "1; mode=block")
	c.Header("Content-Security-Policy", "default-src 'self'")
	c.Next()
}
```

### 6.3 Password Reset Mechanism
**Current**: No password reset functionality
**Improvement**: Add secure password reset
```go
// handlers/auth.go - ForgotPasswordHandler
// Send email with time-limited reset token
// Validate token expiry before allowing reset
// Audit trail for password changes
```

---

## 7. Monitoring & Observability

### 7.1 Health Check Endpoint
**Current**: No health check for load balancers
**Improvement**: Add /health endpoint
```go
// handlers/health.go
func HealthHandler(c *gin.Context) {
	// Check database connectivity
	err := database.DB.Init()
	status := "healthy"
	code := http.StatusOK
	if err != nil {
		status = "unhealthy"
		code = http.StatusServiceUnavailable
	}
	c.JSON(code, gin.H{"status": status})
}

// In routes: httpsRouter.GET("/health", handlers.HealthHandler)
```

### 7.2 Metrics Collection
**Current**: No metrics on requests, latency, errors
**Improvement**: Add Prometheus metrics
```go
// middleware/metrics.go
import "github.com/prometheus/client_golang/prometheus"

var (
	httpRequestsTotal = prometheus.NewCounterVec(...)
	httpRequestDuration = prometheus.NewHistogramVec(...)
	dbQueryErrors = prometheus.NewCounterVec(...)
)

// Register and track in middleware/handlers
```

---

## 8. Configuration & Deployment

### 8.1 Configuration Validation Helper
**Current**: Manual validation of each env var
**Improvement**: Centralized validation
```go
// config/validator.go
func ValidateConfig() error {
	requiredVars := []string{"BASE_API", "TLS_PATH", "SECRET_KEY", ...}
	for _, v := range requiredVars {
		if os.Getenv(v) == "" {
			return fmt.Errorf("required env var missing: %s", v)
		}
	}

	// Validate formats
	if !strings.HasPrefix(os.Getenv("BASE_API"), "/") {
		return fmt.Errorf("BASE_API must start with /")
	}

	if _, err := time.ParseDuration(os.Getenv("TOKEN_EXPIRATION_PERIOD")); err != nil {
		return fmt.Errorf("invalid TOKEN_EXPIRATION_PERIOD: %v", err)
	}

	return nil
}
```

### 8.2 Graceful Shutdown Enhancement
**Current**: Basic signal handling
**Improvement**: Drain in-flight requests
```go
// main.go - improve shutdown
var wg sync.WaitGroup
shutdownChan := make(chan struct{})

// In middleware:
func GracefulShutdownMiddleware(c *gin.Context) {
	wg.Add(1)
	defer wg.Done()
	select {
	case <-shutdownChan:
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "server shutting down"})
		c.Abort()
		return
	default:
		c.Next()
	}
}

// On signal:
close(shutdownChan)
wg.Wait() // Wait for in-flight requests
```

---

## 9. Developer Experience

### 9.1 Constants Package for Magic Strings
**Current**: Hard-coded strings scattered across code
**Improvement**: Centralize constants
```go
// constants/constants.go
const (
	RoleAdmin = "admin"
	RoleGuest = "guest"

	ErrorUserNotFound = "USER_NOT_FOUND"
	ErrorInvalidPassword = "INVALID_PASSWORD"

	MinPasswordLength = 8
	MinUsernameLength = 3
	MaxUsernameLength = 32
)

// Usage: database.DB.GetRoleByName(constants.RoleGuest)
```

### 9.2 Dependency Injection for Testing
**Current**: Direct use of global `database.DB`
**Improvement**: Accept Store as parameter
```go
// handlers/user.go
type UserHandlers struct {
	store database.Store
	config *config.Config
	logger *slog.Logger
}

func NewUserHandlers(store database.Store, cfg *config.Config, logger *slog.Logger) *UserHandlers {
	return &UserHandlers{store, cfg, logger}
}

// In handlers:
func (h *UserHandlers) RegisterHandler(c *gin.Context) {
	if err := h.store.CreateUser(&user); err != nil {
		...
	}
}
```

### 9.3 Mock Improvements
**Current**: MockStore exists but could be more realistic
**Improvement**: Add mock helpers for testing
```go
// database/mockstore.go
func (m *MockStore) WithUser(username, password string) *MockStore {
	// Pre-populate with user for test
	return m
}

// Usage in tests:
store := NewMockStore().
	WithUser("admin", "hashedpass").
	WithRole("admin").
	WithRole("guest")
```

---

## Priority Implementation Order

| Priority | Feature | Effort | Impact |
|----------|---------|--------|--------|
| 🔴 High | Structured Logging | Medium | High (observability) |
| 🔴 High | API Error Standardization | Low | High (developer UX) |
| 🔴 High | Password Complexity | Low | High (security) |
| 🟡 Medium | Health Check Endpoint | Low | High (ops) |
| 🟡 Medium | Rate Limiting Implementation | Medium | High (security) |
| 🟡 Medium | Pagination | Medium | Medium (scalability) |
| 🟡 Medium | Query Optimization (Preload) | Low | Medium (performance) |
| 🟠 Low | HSTS Headers | Low | Medium (security) |
| 🟠 Low | Request ID Tracking | Low | Low (troubleshooting) |

