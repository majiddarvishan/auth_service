# Auth Service Proxy - Bug Analysis & Improvement Suggestions

## Critical Bugs

### 1. **Config Loading Bug - Wrong Variable Check (config.go:51)**
**Severity**: CRITICAL
**Issue**: Line 51 checks `SecretKey` instead of `p` for TOKEN_EXPIRATION_PERIOD
```go
p := os.Getenv("TOKEN_EXPIRATION_PERIOD")
if SecretKey == "" {  // ❌ WRONG - should check 'p'
    log.Fatal("TOKEN_EXPIRATION_PERIOD is not set in .env file")
}
```
**Impact**: Missing TOKEN_EXPIRATION_PERIOD env var won't be caught; service will fail with wrong error message
**Fix**: Change to `if p == ""`

---

### 2. **Token Claims Type Mismatch (middleware/accounting.go:27-31)**
**Severity**: CRITICAL
**Issue**: Extracting "user" from JWT claims as string, but it's stored as uint (ID)
```go
// In user.go (LoginHandler):
claims := jwt.MapClaims{
    "user": user.ID,  // ✅ Stored as ID (uint)
    ...
}

// In accounting.go:
username, ok := claims["user"].(string)  // ❌ Expecting string, getting ID
```
**Impact**: DynamicAccountingMiddleware will fail for any request; accounting checks will never work
**Fix**: Change to extract user ID, then fetch username from database

---

### 3. **Password Validation Missing (handlers/user.go:48-53)**
**Severity**: HIGH
**Issue**: RegisterHandler accepts any password length/complexity
```go
if req.Username == "" || req.Password == "" {
    c.JSON(http.StatusBadRequest, gin.H{"error": "Username and password are required"})
    return
}
// No validation on password strength
```
**Impact**: Users can create accounts with single-character passwords
**Fix**: Add minimum length (8+ chars) and complexity requirements

---

### 4. **SQL Injection in DeleteUserHandler (handlers/user.go:264)**
**Severity**: MEDIUM
**Issue**: While GORM uses parameterized queries, no existence check before delete
```go
func DeleteUserHandler(c *gin.Context) {
    username := c.Param("username")
    err := database.DB.DeleteUserByUsername(username)
    if err != nil {
        c.JSON(http.StatusInternalServerError, ...)  // Treats all errors the same
    }
}
```
**Impact**: Returns 500 for non-existent users instead of 404; poor UX
**Fix**: Check user existence first, return 404 if not found

---

### 5. **Race Condition in Dynamic Route Management (proxy/proxy.go:18-19)**
**Severity**: HIGH
**Issue**: `deletedCustomRoute` map is not thread-safe
```go
var deletedCustomRoute map[string]*database.CustomEndpoint

func DeleteEndpoint(ep *database.CustomEndpoint) {
    deletedCustomRoute[ep.Path+"/*path"] = ep  // ❌ No mutex
}

func ProxyToEndpoint(c *gin.Context, ep *database.CustomEndpoint) {
    _, exists := deletedCustomRoute[ep.Path+"/*path"]  // ❌ Concurrent read
}
```
**Impact**: Concurrent deletes/proxies can cause panics or undefined behavior
**Fix**: Use `sync.RWMutex` to protect map access

---

### 6. **Dropped Errors in Proxy Director (proxy/proxy.go:75-95)**
**Severity**: MEDIUM
**Issue**: JWT decoding errors are printed but not returned; request continues with incomplete headers
```go
decoded, err := base64.RawURLEncoding.DecodeString(payload)
if err != nil {
    fmt.Println("Error decoding payload:", err)
    return  // ❌ Request continues without user headers!
}
```
**Impact**: Upstream services receive requests without user-id/user-name headers when JWT parsing fails
**Fix**: Return error or abort request; don't silently continue

---

## High-Priority Issues

### 7. **Missing Request Timeout (proxy/proxy.go & middleware/accounting.go)**
**Issue**: HTTP calls to accounting service and reverse proxy have no timeout
```go
resp, err := http.Post(accountingURL, "application/json", bytes.NewBuffer(jsonPayload))
// No timeout - could hang indefinitely
```
**Impact**: Request can hang forever if accounting service is down
**Fix**: Use `http.Client` with `Timeout` field

---

### 8. **No Validation of Custom Endpoint URLs (handlers/custom_endpoints.go:100-107)**
**Severity**: MEDIUM
**Issue**: Weak URL validation allows invalid endpoints
```go
for _, endpoint := range req.Endpoints {
    if endpoint == "" || !strings.HasPrefix(endpoint, "http") {
        c.JSON(http.StatusBadRequest, ...)
    }
}
// ❌ Only checks for "http" prefix, doesn't validate URL structure
```
**Impact**: Can register endpoints like "http://" or malformed URLs; causes 500 in proxy
**Fix**: Use `url.Parse()` for proper validation

---

### 9. **Path Traversal Risk in Custom Endpoints (handlers/custom_endpoints.go:122)**
**Severity**: MEDIUM
**Issue**: Custom endpoint paths aren't validated
```go
req.Path += "/*path"  // ❌ No validation of req.Path
```
**Impact**: Admin could create routes like `/../../../admin` to bypass path restrictions
**Fix**: Validate that path doesn't contain `..` or other traversal patterns

---

### 10. **No Rate Limiting**
**Issue**: No rate limiting on auth endpoints (login, register)
**Impact**: Vulnerable to brute force attacks and DDoS
**Fix**: Add rate limiter middleware (e.g., `github.com/ulule/limiter`)

---

## Logic Issues

### 11. **Inconsistent Error Handling in AccountingMiddleware (middleware/accounting.go)**
**Issue**: Claims extraction expects string for "user" but JWT stores uint
```go
username, ok := claims["user"].(string)  // ❌ Won't work with uint
```
**Impact**: This middleware never works due to type mismatch; no error reported
**Fix**: Extract user ID (uint), query database for username

---

### 12. **Missing Cascade Delete Behavior (database/models.go)**
**Issue**: Foreign key constraint won't prevent orphaned records
```go
type User struct {
    RoleID   uint
    Role     Role    `gorm:"foreignKey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
```
**While CASCADE is set, if admin is deleted, users become role-less. Better to prevent admin deletion or handle gracefully.**

---

### 13. **No Deduplication of Accounting Service Calls**
**Issue**: Every dynamic endpoint with `NeedAccounting=true` calls accounting service separately
**Impact**: High latency; accounting service becomes bottleneck
**Improvement**: Cache accounting rules or batch checks

---

### 14. **Custom Endpoint Creation Response Issue (handlers/custom_endpoints.go:125)**
**Issue**: Response returns full endpoint object including Endpoints array
```go
c.JSON(http.StatusOK, gin.H{"message": "...", "endpoint": req})
```
**Risk**: If endpoints contain sensitive URLs, they're exposed in response
**Fix**: Only return minimal data (path, method)

---

## Security Issues

### 15. **JWT Secret Key Exposed in Logs (proxy/proxy.go)**
**Issue**: Debug prints contain JWT data
```go
fmt.Println("Error decoding payload:", err)
fmt.Printf("invalid user id %v\n", value)
fmt.Printf("call ProxyToEndpoint %v\n", ep)
```
**Impact**: If logs are accessed, auth data is exposed
**Fix**: Use structured logging; avoid printing sensitive data

---

### 16. **CORS Allows All Origins (routes/routes.go:41)**
**Issue**: `AllowOrigins: []string{"*"}` allows any origin
**Impact**: CSRF attacks possible; should whitelist origins
**Fix**: Load allowed origins from config

---

### 17. **No HTTPS Redirect (routes/routes.go:65-80)**
**Issue**: HTTP -> HTTPS redirect is commented out
**Impact**: HTTP connections not forced to HTTPS
**Fix**: Uncomment and test the redirect logic

---

### 18. **No Input Sanitization (all handlers)**
**Issue**: No sanitization of string inputs (username, passwords passed to bcrypt)
**Impact**: Potential for noisy logs or timing attacks
**Fix**: Add trim/validation on inputs

---

## Performance Issues

### 19. **Inefficient User Lookup in Proxy (proxy/proxy.go:86-90)**
**Issue**: Database query per request in reverse proxy director
```go
user, err := database.DB.GetUserByID(uint(val))
```
**Impact**: 1 extra DB query per proxied request; slows down all traffic
**Fix**: Cache user data with TTL; use embedded claims for basic info

---

### 20. **No Connection Pooling Configuration**
**Issue**: Database connection pool not explicitly configured
**Impact**: Default pool size may be insufficient for load
**Fix**: Set `MaxOpenConns`, `MaxIdleConns` in database initialization

---

## Code Quality Issues

### 21. **Inconsistent Error Messages**
**Issue**: Some endpoints return "error", others return "details"
- Different error response formats across handlers
**Fix**: Create standardized error response struct

---

### 22. **Dead Code in Routes (routes/routes.go:53-65, commented blocks)**
**Issue**: Large blocks of commented code cluttering file
**Fix**: Remove or move to separate branch

---

### 23. **No Structured Logging**
**Issue**: Mix of `log.Println`, `fmt.Println`, `fmt.Printf`
**Impact**: Hard to parse logs; no log levels
**Fix**: Use structured logger (e.g., `slog`, `logrus`)

---

### 24. **Missing Environment Variable Validation**
**Issue**: TLS_PATH validation only checks existence, not readability
```go
TLSPath = os.Getenv("TLS_PATH")
if TLSPath == "" {
    log.Fatal("TLS_PATH is not set")
}
// ❌ Doesn't verify cert.pem and key.pem exist
```
**Fix**: Validate file existence and readability at startup

---

### 25. **No Graceful Shutdown**
**Issue**: Server stops abruptly on signal; no connection drain time
**Impact**: Active requests killed; clients get errors
**Fix**: Implement graceful shutdown with timeout

---

## Recommendations Summary

| Priority | Issue | Fix Time |
|----------|-------|----------|
| 🔴 Critical | Config loading typo | 2 min |
| 🔴 Critical | JWT claims type mismatch | 15 min |
| 🔴 Critical | Race condition in routes map | 20 min |
| 🟠 High | Missing password validation | 10 min |
| 🟠 High | No request timeouts | 10 min |
| 🟠 High | CORS security | 5 min |
| 🟡 Medium | Error handling consistency | 30 min |
| 🟡 Medium | Input validation | 30 min |
| 🟡 Medium | Structured logging | 1 hour |
| 🟡 Medium | Database connection pooling | 15 min |

## Quick Wins (Start Here)
1. Fix config.go line 51 (`if p == ""`)
2. Fix JWT claims extraction in accounting middleware
3. Add sync.RWMutex to deletedCustomRoute
4. Add request timeouts to HTTP clients
5. Change CORS AllowOrigins from `"*"` to environment-configured list
