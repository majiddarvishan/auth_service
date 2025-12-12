package constants

// Role constants
const (
	RoleAdmin = "admin"
	RoleGuest = "guest"
	RoleUser  = "user"
)

// Error codes (machine-readable)
const (
	ErrorUserNotFound       = "USER_NOT_FOUND"
	ErrorInvalidPassword    = "INVALID_PASSWORD"
	ErrorUserAlreadyExists  = "USER_ALREADY_EXISTS"
	ErrorInvalidUsername    = "INVALID_USERNAME"
	ErrorInvalidPassword2   = "INVALID_PASSWORD_STRENGTH"
	ErrorUnauthorized       = "UNAUTHORIZED"
	ErrorForbidden          = "FORBIDDEN"
	ErrorRoleNotFound       = "ROLE_NOT_FOUND"
	ErrorEndpointNotFound   = "ENDPOINT_NOT_FOUND"
	ErrorInternalServer     = "INTERNAL_SERVER_ERROR"
	ErrorInvalidJSON        = "INVALID_JSON"
	ErrorInvalidPathTravers = "INVALID_PATH"
	ErrorRateLimitExceeded  = "RATE_LIMIT_EXCEEDED"
	ErrorInvalidRefreshToken = "INVALID_REFRESH_TOKEN"
	ErrorRefreshTokenExpired = "REFRESH_TOKEN_EXPIRED"
)

// Password validation constraints
const (
	MinPasswordLength = 8
	MinUsernameLength = 3
	MaxUsernameLength = 32
	MinPathLength     = 1
	MaxPathLength     = 255
)

// Token claims keys
const (
	ClaimKeyUserID   = "user"
	ClaimKeyRole     = "role"
	ClaimKeyExpiry   = "exp"
	ClaimKeyIssuedAt = "iat"
)

// HTTP header constants
const (
	HeaderAuthorization = "Authorization"
	HeaderContentType   = "Content-Type"
	HeaderRequestID     = "X-Request-ID"
	HeaderXForwardedFor = "X-Forwarded-For"
)

// Default values
const (
	DefaultRole               = RoleGuest
	DefaultAccountingEndpoint = "http://localhost:8082"
	DefaultHTTPTimeout        = 5 // seconds
	DefaultRateLimitPerSecond = 100
	DefaultRateLimitWindow    = 60 // seconds
	DefaultPageLimit          = 10
	DefaultMaxPageLimit       = 100
)

// API paths
const (
	PathHealthCheck = "/health"
	PathMetrics     = "/metrics"
	PathVersion     = "/version"
)

// Database operation modes
const (
	DBModePGSQL = "postgres"
	DBModeMock  = "mock"
)
