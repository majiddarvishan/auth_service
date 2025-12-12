# Copilot Instructions for Auth Service Proxy

## Project Architecture

**Purpose**: A Go-based authentication/authorization gateway with dynamic reverse proxy capabilities. Handles user registration/login, JWT token validation, RBAC, custom endpoint proxying, and API accounting/quotas.

**Core Stack**: Gin HTTP framework, PostgreSQL (with mock fallback), JWT tokens (golang-jwt/v5), bcrypt password hashing, GORM ORM, structured logging (slog).

### Service Components

| Component | Role | Key Files |
|-----------|------|-----------|
| **handlers** | HTTP endpoint logic (auth, user mgmt, roles, captcha, custom routes, health) | `handlers/user.go`, `handlers/admin.go`, `handlers/custom_endpoints.go`, `handlers/health.go` |
| **middleware** | Request validation/enrichment (JWT auth, role checks, accounting, request ID, security headers, rate limiting) | `middleware/auth.go`, `middleware/role.go`, `middleware/accounting.go`, `middleware/request_id.go`, `middleware/security_headers.go`, `middleware/rate_limit.go` |
| **proxy** | Load-balanced reverse proxy with dynamic route support | `proxy/proxy.go` (creates multi-target reverse proxies, manages deleted routes) |
| **database** | Data layer abstraction with Store interface (PostgreSQL/mock implementations) | `database/store.go`, `database/pgstore.go`, `database/mockstore.go`, `database/models.go` |
| **routes** | Gin router setup with CORS, dynamic groups, Swagger docs, middleware registration | `routes/routes.go` |
| **config** | Environment variable loading, validation, security configurations | `config/config.go`, `config/validator.go` |
| **validation** | Input validation (password strength, username format, endpoint paths) | `validation/validation.go` |
| **constants** | Centralized constants (roles, error codes, defaults, claim keys) | `constants/constants.go` |
| **types** | Shared types (APIError, APISuccess, Health status, pagination) | `types/types.go` |
| **logger** | Structured logging with JSON output | `logger/logger.go` |

### Data Flow

1. **Request Entry**: HTTP/HTTPS → RequestIDMiddleware (add trace ID) → SecurityHeadersMiddleware (add headers) → RateLimitMiddleware (check quota)
2. **Route Matching**: Gin router matches to handler (health/version don't require auth, others do)
3. **Auth Routes**: POST `/users` → validation.ValidateUsername/Password → bcrypt hash → DB insert
4. **Login/Token**: POST `/login` → lookup user → bcrypt.CompareHashAndPassword → generate JWT with claims["user"]=ID
5. **Protected Routes**: Bearer token → AuthMiddleware → validate JWT → extract claims → lookup user by ID
6. **Accounting**: Optional DynamicAccountingMiddleware → call accounting service → track usage
7. **Custom Routes**: Dynamic path → reverse proxy to target URL(s) with random load balancing
8. **Error Responses**: Standardized APIError type with error codes, messages, request ID, timestamp

## Critical Architectural Decisions

### 1. Store Interface Pattern
Database ops use a `Store` interface (`database/store.go`) with two implementations: **PGStore** (production) and **MockStore** (testing). Switch via CLI flag: `go run main.go -d postgres` or `go run main.go -d mock`. This avoids test database dependencies.

### 2. Structured Logging
Uses Go 1.21+ `slog` package for JSON logging instead of `log` package. Initialize in main:
```go
logger.Init()  // Creates singleton with JSON handler
log := logger.Get()
log.Info("event", "key1", value1, "key2", value2)  // Output: JSON
```
Benefits: Structured format for log aggregation, debug-level granularity, request correlation via request_id field.

### 3. Request ID Middleware
Every request gets unique X-Request-ID (UUID) for end-to-end tracing. Automatically propagated in responses and logged. Access in handlers via: `c.Get("request_id")`.

### 4. Standardized Error Responses
All errors return `types.APIError` with: code (machine-readable), message (human-readable), timestamp, requestID. Example:
```go
requestID, _ := c.Get("request_id")
c.JSON(http.StatusBadRequest, types.APIError{
    Code: constants.ErrorInvalidUsername,
    Message: "Username must be 3-32 chars",
    Timestamp: time.Now().Unix(),
    RequestID: requestID.(string),
})
```

### 5. Input Validation Package
Centralized validators in `validation/validation.go`:
- `ValidatePasswordStrength()`: 8+ chars, uppercase, lowercase, digit, special char
- `ValidateUsername()`: 3-32 chars, alphanumeric + underscore only
- `ValidateEndpointPath()`: Prevents path traversal (`..`), double slashes, null bytes
- `ValidateEndpointTargets()`: Ensures URLs are valid HTTP/HTTPS
- `ValidateHTTPMethod()`: Whitelists valid HTTP methods

### 6. Constants Package
All magic strings centralized in `constants/constants.go`:
- Role names: `constants.RoleAdmin`, `constants.RoleGuest`
- Error codes: `constants.ErrorUserNotFound`, `constants.ErrorInvalidPassword`
- Claim keys: `constants.ClaimKeyUserID`, `constants.ClaimKeyRole`
- Use instead of hardcoding: `database.DB.GetRoleByName(constants.DefaultRole)`

### 7. Rate Limiting
In-memory rate limiter in `middleware/rate_limit.go`. Default: 100 req/min per IP. Usage:
```go
httpsRouter.Use(middleware.RateLimitMiddleware(100, 60*time.Second))
```
Returns HTTP 429 with APIError when exceeded.

### 8. Security Headers Middleware
Adds HSTS, CSP, X-Frame-Options, etc. via `middleware/security_headers.go`. Applied globally in `routes.go`.

### 9. Configuration Validation
`config/validator.go` ValidateConfig() checks all required env vars are set and correctly formatted. Called in main before starting service.

## Build & Test Commands

```bash
# Load from .env; required vars: BASE_API, TLS_PATH, SECRET_KEY, TOKEN_EXPIRATION_PERIOD, DB_*
go run main.go -d postgres          # Run with PostgreSQL
go run main.go -d mock              # Run with in-memory mock (no DB needed)

# Generate Swagger docs (requires swag CLI installed)
swag init -g main.go -o docs

# Build
go build -o auth_service main.go

# Tests (use mock store to avoid DB dependency)
go test ./...

# Docker
docker build -f deployments/Dockerfile -t auth-service .
docker-compose -f deployments/docker-compose.yml up

# Add required dependency
go get github.com/google/uuid
```

## Project-Specific Patterns & Conventions

### Request/Response Structure
All handlers follow this pattern:
- Bind JSON: `c.ShouldBindJSON(&req)` with validation
- Validate input: use functions from `validation/` package
- Responses: `c.JSON(statusCode, types.APIError/APISuccess)` with request ID
- Log activity: `logger.Info/Error()` for logging instead of print statements
- Example: [handlers/user.go#L43-L100](handlers/user.go#L43-L100)

### Database Models Location
All data structures: [database/models.go](database/models.go). Use GORM struct tags (`gorm:"..."`) for DB mappings. Example: `User` has `ID`, `Username`, `PasswordHash`, `RoleID` (foreign key).

### Middleware Chain
Routes use gin groups with chained middleware. Example: [routes/routes.go#L61-L80](routes/routes.go#L61-L80) shows global middleware applied to router, then role-based groups nested.

### Middleware Execution Order (in routes.go)
1. RequestIDMiddleware - adds request_id to context
2. SecurityHeadersMiddleware - adds security headers
3. RateLimitMiddleware - checks rate limit
4. CORS middleware
5. Route-specific middleware (AuthMiddleware, RoleMiddleware, DynamicAccountingMiddleware)

### Environment Variables (Required)
Loaded in [config/config.go](config/config.go) and validated in [config/validator.go](config/validator.go). **ALL** are checked; missing vars cause fatal exit:
- `BASE_API`: URL prefix for all routes (e.g., `/v1/api`)
- `TLS_PATH`: Directory containing `localhost.pem` and `localhost-key.pem`
- `SECRET_KEY`: 32+ hex chars (generate: `openssl rand -hex 32`)
- `TOKEN_EXPIRATION_PERIOD`: Go duration string (e.g., `24h`)
- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`: PostgreSQL creds
- `DB_MAX_OPEN_CONNS`: (optional) Max open connections, default 25
- `DB_MAX_IDLE_CONNS`: (optional) Max idle connections, default 5

### Error Code Standards
Use `constants` package for all error codes. Examples:
- `constants.ErrorUserNotFound`: HTTP 404
- `constants.ErrorInvalidPassword`: HTTP 401
- `constants.ErrorForbidden`: HTTP 403
- `constants.ErrorRateLimitExceeded`: HTTP 429

### Code Generation (Swagger)
Use swag CLI to auto-generate API docs from handler comments. All public handlers have `@Summary`, `@Description`, `@Param`, `@Success`, `@Failure` annotations. Docs output to `docs/` → served at `/swagger/index.html`.

## API Endpoints Overview

### Health & Info (No Auth Required)
- `GET /health` - Health check (DB connectivity)
- `GET /version` - Service version info

### Auth (No Auth Required for Register/Login)
- `POST /v1/api/users` - Register new user (validates password strength, username format)
- `POST /v1/api/login` - Login (returns JWT token)

### Captcha (CORS-enabled)
- `GET /v1/api/captcha/new` - Get new captcha ID
- `GET /v1/api/captcha/image/:captchaId` - Get captcha image

### Protected Routes (Auth Required)
- `GET /v1/api/admin` - Admin dashboard (requires admin role)
- `PUT/DELETE /v1/api/users/{username}` - Manage users

### Admin Routes (Auth + Admin Role Required)
- `POST /v1/api/admin/custom-endpoints` - Create dynamic endpoint (validates path, targets, HTTP method)
- `DELETE /v1/api/admin/custom-endpoints` - Delete dynamic endpoint
- `GET /v1/api/admin/roles` - List roles
- `POST /v1/api/admin/roles` - Create role

### Custom Routes (Dynamic)
Created at runtime via admin API. Each gets:
- AuthMiddleware (JWT validation)
- Optional DynamicAccountingMiddleware (quota checks)
- ProxyToEndpoint handler (reverse proxy with load balancing)

## External Dependencies

- **Gin**: Web framework (`github.com/gin-gonic/gin`)
- **GORM**: ORM (`gorm.io/gorm`, `gorm.io/driver/postgres`)
- **golang-jwt**: JWT token creation/parsing (`github.com/golang-jwt/jwt/v5`)
- **bcrypt**: Password hashing (`golang.org/x/crypto/bcrypt`)
- **swaggo**: Swagger auto-docs (`github.com/swaggo/gin-swagger`)
- **captcha**: Captcha generation (`github.com/dchest/captcha`)
- **godotenv**: Load .env files (`github.com/joho/godotenv`)
- **google/uuid**: UUID generation for request IDs (`github.com/google/uuid`)
- **slog**: Structured logging (stdlib, Go 1.21+)

## Testing Approach

Use **MockStore** for unit tests (in-memory, no DB):
- Set `-d mock` flag or call `database.NewStore("mock")`
- Example test pattern: [handlers/user_test.go](handlers/user_test.go)
- Pre-populate mock data via `mockstore.go` methods before test

Use **PostgreSQL** for integration tests. Database schema auto-created by GORM migrations (defined in model structs, not separate SQL files).

## File Navigation Tips

- **Entry point**: [main.go](main.go) (logger init, config validation, graceful shutdown)
- **Routes definition**: [routes/routes.go](routes/routes.go) (where to add new endpoints, middleware setup)
- **All handlers**: `handlers/` (auth, users, roles, custom endpoints, captcha, health)
- **All middleware**: `middleware/` (auth validation, role checks, accounting, request ID, security, rate limit)
- **Validation logic**: [validation/validation.go](validation/validation.go) (password, username, path validators)
- **Constants**: [constants/constants.go](constants/constants.go) (all magic strings, error codes, defaults)
- **Error types**: [types/types.go](types/types.go) (APIError, APISuccess, HealthStatus structures)
- **Logging**: [logger/logger.go](logger/logger.go) (structured JSON logging)
- **Proxy logic**: [proxy/proxy.go](proxy/proxy.go) (reverse proxy, load balancing, thread-safe route management)
- **Database models**: [database/models.go](database/models.go)
- **HTTPS server**: [routes/routes.go#L173](routes/routes.go#L173) (ListenAndServeTLS)

## Common Tasks

**Adding a new protected endpoint**:
1. Create handler in `handlers/` with proper error handling using `types.APIError` and `constants`
2. Add Swagger comments (@Summary, @Description, etc.)
3. Add route in [routes/routes.go](routes/routes.go) with appropriate middleware
4. Use `logger.Info/Error()` for logging instead of print statements

Example:
```go
// handlers/status.go
func StatusHandler(c *gin.Context) {
    requestID, _ := c.Get("request_id")

    // Get user from token claims
    claims := c.MustGet("claims").(jwt.MapClaims)
    userID := uint(claims[constants.ClaimKeyUserID].(float64))

    user, err := database.DB.GetUserByID(userID)
    if err != nil {
        logger.Error("User not found", "user_id", userID, "error", err.Error())
        c.JSON(http.StatusNotFound, types.APIError{
            Code: constants.ErrorUserNotFound,
            Message: "User not found",
            Timestamp: time.Now().Unix(),
            RequestID: requestID.(string),
        })
        return
    }

    logger.Info("Status check", "username", user.Username, "request_id", requestID)
    c.JSON(http.StatusOK, types.APISuccess{
        Data: gin.H{"status": "active"},
        Timestamp: time.Now().Unix(),
        RequestID: requestID.(string),
    })
}
```

**Adding a database model**:
1. Define struct in [database/models.go](database/models.go) with GORM tags
2. Implement CRUD methods in [database/store.go](database/store.go) interface
3. Implement in both [database/pgstore.go](database/pgstore.go) and [database/mockstore.go](database/mockstore.go)
4. GORM auto-migrations run on DB init

**Adding input validation**:
1. Add validator function to [validation/validation.go](validation/validation.go)
2. Call in handler: `if err := validation.ValidateXXX(input); err != nil { ... }`
3. Return `types.APIError` with appropriate error code from `constants`

**Debugging JWT/Auth issues**:
- Print claims: `claims := c.MustGet("claims").(jwt.MapClaims)` then inspect
- Check token expiry: `exp := claims[constants.ClaimKeyExpiry].(float64)` compare to `time.Now().Unix()`
- Verify request ID: `requestID, _ := c.Get("request_id")` (should be UUID string)
- Check logs: `logger.Info()` outputs JSON with all context

**Adding rate-limited endpoint**:
Already global (100 req/min), but to customize per-route:
```go
rootGroup.GET("/expensive",
    middleware.RateLimitMiddleware(10, 60*time.Second),  // 10 req/min
    handlers.ExpensiveHandler,
)
```

## Graceful Shutdown

Server now waits up to 30 seconds for in-flight requests to complete on SIGINT/SIGTERM. Coordinated via:
- `shutdownChan` closed on signal
- `sync.WaitGroup` tracks active requests
- Context timeout prevents indefinite wait

No changes needed in handlers unless adding long-running operations.

## Improvements Implemented

This codebase now includes:
✅ Structured logging with slog (JSON output)
✅ Standardized APIError responses with error codes
✅ Request ID tracking (X-Request-ID header)
✅ Password complexity validation (8+ chars, uppercase, lowercase, digit, special)
✅ Username format validation (alphanumeric + underscore only)
✅ Endpoint path validation (prevents traversal attacks)
✅ Security headers (HSTS, CSP, X-Frame-Options)
✅ Rate limiting (100 req/min per IP)
✅ Health check endpoint (/health, /version)
✅ Configuration validation on startup
✅ Constants package (all magic strings centralized)
✅ Graceful shutdown with request draining
✅ In-memory database connection pooling ready
✅ GORM Preload optimization in queries

## Known Limitations & Future Improvements

1. **Refresh Tokens**: Not yet implemented. Current tokens have single fixed expiry.
2. **Password Reset**: No password reset mechanism yet.
3. **Metrics**: No Prometheus metrics collection yet (infrastructure ready in middleware).
4. **Soft Deletes**: Using hard deletes, consider GORM soft deletes for audit trail.
5. **Pagination**: Not yet implemented on list endpoints (`/v1/api/admin/roles`, etc).
6. **Dependency Injection**: Handlers still use global `database.DB` (ready for refactor).
