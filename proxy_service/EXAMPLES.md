# API Examples - Auth Service Proxy

This guide shows practical examples of using the auth service with all the new improvements.

## Prerequisites

```bash
# Setup environment
export DB_HOST=127.0.0.1
export DB_PORT=5432
export DB_USER_NAME=postgres
export DB_PASSWORD=postgres
export DB_NAME=proxy_db
export BASE_API=/api/v1
export SECRET_KEY=f78973efc0c0664995e2bb055bb2cac6779597a5294685f069229c909358f54a
export TOKEN_EXPIRATION_PERIOD=24h
export TLS_PATH=.

# Run with mock database (no PostgreSQL needed)
go run main.go -d mock

# Or run with PostgreSQL
go run main.go -d postgres
```

---

## 1. Health Check (No Auth Required)

### Check Service Health
```bash
curl -X GET http://localhost:8080/health
```

**Response (200 OK):**
```json
{
  "status": "healthy",
  "database": "connected",
  "version": "1.0.0",
  "timestamp": 1702431000
}
```

**Response (503 Service Unavailable - DB down):**
```json
{
  "status": "unhealthy",
  "database": "disconnected",
  "version": "1.0.0",
  "timestamp": 1702431000
}
```

### Get Service Version
```bash
curl -X GET http://localhost:8080/version
```

**Response:**
```json
{
  "version": "1.0.0",
  "timestamp": 1702431000
}
```

---

## 2. User Registration & Login

### Register New User (with validation)

**Valid password** (8+ chars, uppercase, lowercase, digit, special):
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "password": "SecurePass123!"
  }'
```

**Response (201 Created):**
```json
{
  "data": {
    "id": 1,
    "username": "john_doe",
    "role": "user"
  },
  "message": "User registered successfully",
  "timestamp": 1702431000,
  "request_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

### Invalid Password (too short)
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "password": "short"
  }'
```

**Response (400 Bad Request):**
```json
{
  "code": "INVALID_PASSWORD_STRENGTH",
  "message": "password must be at least 8 characters long",
  "details": null,
  "timestamp": 1702431000,
  "request_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

### Invalid Username (special chars)
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john@doe!",
    "password": "SecurePass123!"
  }'
```

**Response (400 Bad Request):**
```json
{
  "code": "INVALID_USERNAME",
  "message": "username can only contain letters, numbers, and underscores",
  "details": null,
  "timestamp": 1702431000,
  "request_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

### Login (Get JWT Token)
```bash
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "password": "SecurePass123!"
  }'
```

**Response (200 OK):**
```json
{
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_in": 86400,
    "token_type": "Bearer"
  },
  "message": "Login successful",
  "timestamp": 1702431000,
  "request_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

### Wrong Password
```bash
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "password": "WrongPassword123!"
  }'
```

**Response (401 Unauthorized):**
```json
{
  "code": "INVALID_PASSWORD",
  "message": "invalid username or password",
  "details": null,
  "timestamp": 1702431000,
  "request_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

---

## 3. Protected Routes (Auth Required)

### Get Admin Dashboard
```bash
curl -X GET http://localhost:8080/api/v1/admin \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Response (200 OK):**
```json
{
  "data": {
    "users": 5,
    "endpoints": 3,
    "roles": 4
  },
  "message": "Admin dashboard data",
  "timestamp": 1702431000,
  "request_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

### Missing Authorization Header
```bash
curl -X GET http://localhost:8080/api/v1/admin
```

**Response (401 Unauthorized):**
```json
{
  "code": "UNAUTHORIZED",
  "message": "authorization header required",
  "details": null,
  "timestamp": 1702431000,
  "request_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

### Invalid Token
```bash
curl -X GET http://localhost:8080/api/v1/admin \
  -H "Authorization: Bearer invalid.token.here"
```

**Response (401 Unauthorized):**
```json
{
  "code": "UNAUTHORIZED",
  "message": "invalid or expired token",
  "details": null,
  "timestamp": 1702431000,
  "request_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

---

## 4. Create Custom Endpoint (Admin Only)

### Valid Custom Endpoint
```bash
curl -X POST http://localhost:8080/api/v1/admin/custom-endpoints \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "path": "/api/users",
    "endpoints": ["https://api.example.com/users"],
    "method": "GET"
  }'
```

**Response (201 Created):**
```json
{
  "data": {
    "id": 1,
    "path": "/api/users",
    "endpoints": ["https://api.example.com/users"],
    "method": "GET",
    "created_at": "2025-12-13T00:50:00Z"
  },
  "message": "Custom endpoint created successfully",
  "timestamp": 1702431000,
  "request_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

### Invalid Path (Path Traversal)
```bash
curl -X POST http://localhost:8080/api/v1/admin/custom-endpoints \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "path": "/api/../../../etc/passwd",
    "endpoints": ["https://api.example.com/users"],
    "method": "GET"
  }'
```

**Response (400 Bad Request):**
```json
{
  "code": "INVALID_PATH",
  "message": "path traversal not allowed",
  "details": null,
  "timestamp": 1702431000,
  "request_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

### Invalid HTTP Method
```bash
curl -X POST http://localhost:8080/api/v1/admin/custom-endpoints \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "path": "/api/users",
    "endpoints": ["https://api.example.com/users"],
    "method": "INVALID"
  }'
```

**Response (400 Bad Request):**
```json
{
  "code": "INVALID_METHOD",
  "message": "invalid HTTP method",
  "details": null,
  "timestamp": 1702431000,
  "request_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

### Invalid Target URL
```bash
curl -X POST http://localhost:8080/api/v1/admin/custom-endpoints \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "path": "/api/users",
    "endpoints": ["not-a-valid-url"],
    "method": "GET"
  }'
```

**Response (400 Bad Request):**
```json
{
  "code": "INVALID_ENDPOINTS",
  "message": "invalid endpoint URL",
  "details": null,
  "timestamp": 1702431000,
  "request_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

---

## 5. Rate Limiting

The service enforces **100 requests per 60 seconds per IP address**.

### Normal Request (Within Limit)
```bash
curl -X GET http://localhost:8080/health
```

**Response (200 OK):**
```json
{
  "status": "healthy",
  "database": "connected",
  "version": "1.0.0",
  "timestamp": 1702431000
}
```

### Exceeding Rate Limit (>100 req/min)
```bash
# After 100+ requests in 60 seconds
curl -X GET http://localhost:8080/health
```

**Response (429 Too Many Requests):**
```json
{
  "code": "RATE_LIMIT_EXCEEDED",
  "message": "rate limit exceeded",
  "details": null,
  "timestamp": 1702431000,
  "request_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

---

## 6. Request Tracing (X-Request-ID)

Every response includes a unique `request_id` for tracing:

```bash
curl -X GET http://localhost:8080/health -v
```

**Response Headers:**
```
X-Request-ID: 550e8400-e29b-41d4-a716-446655440000
Content-Type: application/json
```

**Response Body:**
```json
{
  "status": "healthy",
  "database": "connected",
  "version": "1.0.0",
  "request_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

Use this ID to trace requests through logs:
```bash
# In application logs, search for this ID
grep "550e8400-e29b-41d4-a716-446655440000" /var/log/auth_service.log
```

---

## 7. Security Headers

All responses include industry-standard security headers:

```bash
curl -X GET http://localhost:8080/health -i
```

**Response Headers:**
```
Strict-Transport-Security: max-age=63072000; includeSubDomains; preload
X-Content-Type-Options: nosniff
X-Frame-Options: DENY
X-XSS-Protection: 1; mode=block
Content-Security-Policy: default-src 'self'
Referrer-Policy: strict-origin-when-cross-origin
Permissions-Policy: geolocation=(), microphone=(), camera=()
```

---

## 8. CORS (Cross-Origin Resource Sharing)

### Allowed Origins
Default allowed origins (configurable via `ALLOWED_CORS_ORIGINS`):
- `http://localhost:3000`
- `http://localhost:8080`

### Preflight Request
```bash
curl -X OPTIONS http://localhost:8080/api/v1/login \
  -H "Origin: http://localhost:3000" \
  -H "Access-Control-Request-Method: POST" \
  -i
```

**Response (200 OK):**
```
Access-Control-Allow-Origin: http://localhost:3000
Access-Control-Allow-Methods: GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS
Access-Control-Allow-Headers: Content-Type, Authorization
Access-Control-Max-Age: 86400
```

### Disallowed Origin
```bash
curl -X OPTIONS http://localhost:8080/api/v1/login \
  -H "Origin: http://evil.example.com" \
  -i
```

**Response (200 OK):**
```
No Access-Control headers (origin not allowed)
```

---

## 9. Logging Examples

### Console Output (Structured JSON)
```json
{"time":"2025-12-13T00:50:00.000000000+03:30","level":"INFO","msg":"User registered","username":"john_doe","user_id":1,"request_id":"550e8400-e29b-41d4-a716-446655440000"}
{"time":"2025-12-13T00:50:10.000000000+03:30","level":"INFO","msg":"Login successful","username":"john_doe","request_id":"550e8400-e29b-41d4-a716-446655440001"}
{"time":"2025-12-13T00:50:20.000000000+03:30","level":"ERROR","msg":"Invalid password","username":"john_doe","error":"invalid username or password","request_id":"550e8400-e29b-41d4-a716-446655440002"}
```

### Parse Logs with jq
```bash
# Get all login events
go run main.go -d mock 2>&1 | jq 'select(.msg | contains("Login"))'

# Get all errors
go run main.go -d mock 2>&1 | jq 'select(.level == "ERROR")'

# Get events for a specific request
REQUEST_ID="550e8400-e29b-41d4-a716-446655440000"
go run main.go -d mock 2>&1 | jq "select(.request_id == \"$REQUEST_ID\")"
```

---

## 10. Testing Script

Create `test_api.sh`:

```bash
#!/bin/bash

BASE_URL="http://localhost:8080"
USERNAME="testuser_$(date +%s)"
PASSWORD="TestPass123!"

echo "=== 1. Register User ==="
REGISTER=$(curl -s -X POST $BASE_URL/api/v1/users \
  -H "Content-Type: application/json" \
  -d "{\"username\": \"$USERNAME\", \"password\": \"$PASSWORD\"}")
echo $REGISTER | jq .

echo -e "\n=== 2. Login ==="
LOGIN=$(curl -s -X POST $BASE_URL/api/v1/login \
  -H "Content-Type: application/json" \
  -d "{\"username\": \"$USERNAME\", \"password\": \"$PASSWORD\"}")
TOKEN=$(echo $LOGIN | jq -r '.data.access_token')
echo "Token: $TOKEN"

echo -e "\n=== 3. Access Protected Route ==="
curl -s -X GET $BASE_URL/api/v1/admin \
  -H "Authorization: Bearer $TOKEN" | jq .

echo -e "\n=== 4. Check Health ==="
curl -s -X GET $BASE_URL/health | jq .

echo -e "\n=== 5. Check Version ==="
curl -s -X GET $BASE_URL/version | jq .
```

Run it:
```bash
chmod +x test_api.sh
./test_api.sh
```

---

## Error Code Reference

| Code | HTTP Status | Meaning |
|------|------------|---------|
| `USER_NOT_FOUND` | 404 | User doesn't exist |
| `USER_ALREADY_EXISTS` | 409 | Username taken |
| `INVALID_USERNAME` | 400 | Username format invalid |
| `INVALID_PASSWORD_STRENGTH` | 400 | Password too weak |
| `INVALID_PASSWORD` | 401 | Wrong password |
| `UNAUTHORIZED` | 401 | Missing/invalid token |
| `FORBIDDEN` | 403 | Insufficient permissions |
| `INVALID_PATH` | 400 | Path traversal detected |
| `RATE_LIMIT_EXCEEDED` | 429 | Too many requests |
| `INTERNAL_SERVER_ERROR` | 500 | Server error |

---

## Tips

✅ Always use HTTPS in production (port 8443)
✅ Store tokens securely (httpOnly cookies preferred)
✅ Check `request_id` in responses for debugging
✅ All errors are standardized with `code` + `message` + `request_id`
✅ Logs are JSON-formatted for aggregation systems
✅ Rate limit resets every 60 seconds
✅ Security headers are automatic on all responses
