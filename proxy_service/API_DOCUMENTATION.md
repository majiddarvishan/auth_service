# Auth Service API Documentation

Comprehensive REST API documentation for the Authentication & Authorization Proxy Service with JWT tokens, role-based access control (RBAC), and multi-role support.

## Table of Contents

1. [Base URL](#base-url)
2. [Authentication](#authentication)
3. [Response Format](#response-format)
4. [Error Handling](#error-handling)
5. [Endpoints](#endpoints)
   - [Health & System](#health--system)
   - [Authentication](#authentication-endpoints)
   - [User Management](#user-management)
   - [Role Management](#role-management)
   - [Refresh Token Management](#refresh-token-management)

---

## Base URL

```
https://localhost:8443/api/v1
```

## Authentication

All endpoints except `/login`, `/secure-login`, `/refresh`, and `/logout` require JWT authentication.

### Bearer Token Format

```
Authorization: Bearer <access_token>
```

### Token Claims

```json
{
  "user": 123,                    // User ID
  "role": "admin",                // Primary role
  "exp": 1234567890,              // Expiration timestamp (Unix)
  "iat": 1234567000               // Issued at timestamp
}
```

### Token Expiration

- **Access Token**: 15 minutes (default)
- **Refresh Token**: 7 days

---

## Response Format

### Success Response (2xx)

```json
{
  "data": {
    // Response data here
  },
  "message": "Operation successful",
  "timestamp": 1234567890,
  "request_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

### Paginated Response

```json
{
  "data": [
    // Array of items
  ],
  "total": 100,
  "page": 1,
  "limit": 10,
  "total_pages": 10,
  "timestamp": 1234567890
}
```

### Token Response

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6",
  "expires_in": 900,            // Seconds
  "token_type": "Bearer",
  "timestamp": 1234567890
}
```

### Error Response

```json
{
  "code": "USER_NOT_FOUND",
  "message": "User does not exist",
  "timestamp": 1234567890,
  "request_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

---

## Error Handling

| Status | Code | Meaning |
|--------|------|---------|
| 200 | OK | Request succeeded |
| 400 | BAD_REQUEST | Invalid parameters |
| 401 | UNAUTHORIZED | Invalid/missing auth or credentials |
| 403 | FORBIDDEN | Permission denied (insufficient role) |
| 404 | NOT_FOUND | Resource not found |
| 409 | CONFLICT | Resource already exists |
| 429 | RATE_LIMIT_EXCEEDED | Too many requests (100/min per IP) |
| 500 | INTERNAL_SERVER_ERROR | Server error |

### Error Codes Reference

```
USER_NOT_FOUND              - User doesn't exist
INVALID_PASSWORD            - Wrong password
USER_ALREADY_EXISTS         - Username taken
INVALID_USERNAME            - Invalid username format
INVALID_PASSWORD_STRENGTH   - Password doesn't meet requirements
UNAUTHORIZED                - Missing/invalid token
FORBIDDEN                   - Insufficient permissions
ROLE_NOT_FOUND              - Role doesn't exist
INVALID_JSON                - Malformed JSON
INVALID_REFRESH_TOKEN       - Invalid or expired refresh token
RATE_LIMIT_EXCEEDED         - Too many requests
```

---

## Endpoints

### Health & System

#### GET /health
Check API health status.

**Response:**
```json
{
  "status": "healthy",
  "database": "connected",
  "version": "1.0.0",
  "timestamp": 1234567890
}
```

#### GET /version
Get API version information.

**Response:**
```json
{
  "version": "1.0.0"
}
```

---

### Authentication Endpoints

#### POST /login
Authenticate user and receive JWT tokens.

**Request:**
```json
{
  "username": "john_doe",
  "password": "SecurePass123!"
}
```

**Response:** `200 OK`
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6",
  "expires_in": 900,
  "token_type": "Bearer",
  "timestamp": 1234567890
}
```

**Errors:**
- `401 UNAUTHORIZED` - Invalid credentials
- `400 INVALID_JSON` - Malformed request

---

#### POST /secure-login
Login with CAPTCHA verification (bot protection).

**Request:**
```json
{
  "username": "john_doe",
  "password": "SecurePass123!",
  "captchaId": "captcha_id_from_/captcha/new",
  "captchaSolution": "solution_string"
}
```

**Response:** `200 OK` (Same as `/login`)

**Errors:**
- `401 UNAUTHORIZED` - Invalid credentials or failed CAPTCHA
- `400 INVALID_JSON` - Malformed request

---

#### POST /refresh
Get new access token using refresh token.

**Request:**
```json
{
  "refresh_token": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6"
}
```

**Response:** `200 OK`
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6",
  "expires_in": 900,
  "token_type": "Bearer",
  "timestamp": 1234567890
}
```

**Errors:**
- `401 UNAUTHORIZED` - Invalid or expired refresh token
- `404 NOT_FOUND` - Refresh token not found

---

#### POST /logout
Revoke a specific refresh token.

**Request:**
```json
{
  "refresh_token": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6"
}
```

**Response:** `200 OK`
```json
{
  "message": "Token revoked successfully",
  "timestamp": 1234567890,
  "request_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

---

#### POST /logout-all
Revoke all refresh tokens for current user (logout from all devices).

**Authorization:** Required (Bearer token)

**Response:** `200 OK`
```json
{
  "message": "Logged out from all devices successfully",
  "timestamp": 1234567890,
  "request_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

---

### User Management

#### POST /users
Create a new user account.

**Authorization:** Required (admin role)

**Request:**
```json
{
  "username": "john_doe",
  "password": "SecurePass123!",
  "role": "user"
}
```

**Response:** `201 Created`
```json
{
  "message": "User registered successfully",
  "timestamp": 1234567890,
  "request_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

**Validation Rules:**
- Username: 3-32 chars, alphanumeric + underscore
- Password: Min 8 chars, uppercase, lowercase, digit, special char

**Errors:**
- `409 CONFLICT` - Username already exists
- `400 INVALID_JSON` - Invalid username or weak password
- `403 FORBIDDEN` - Insufficient permissions

---

#### GET /admin
Get admin dashboard with all users and their roles.

**Authorization:** Required (admin role)

**Response:** `200 OK`
```json
{
  "data": {
    "users": [
      {
        "id": 1,
        "username": "john_doe",
        "roles": [
          {
            "id": 1,
            "name": "admin",
            "description": "Administrator role"
          },
          {
            "id": 2,
            "name": "user",
            "description": "Regular user role"
          }
        ],
        "balance": 1000.50,
        "created_at": "2023-01-15T10:30:00Z"
      }
    ]
  },
  "message": "Admin dashboard",
  "timestamp": 1234567890,
  "request_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

---

#### DELETE /users/{username}
Delete a user account.

**Authorization:** Required (admin role)

**Response:** `200 OK`
```json
{
  "message": "User deleted successfully",
  "timestamp": 1234567890,
  "request_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

**Errors:**
- `404 NOT_FOUND` - User doesn't exist
- `403 FORBIDDEN` - Insufficient permissions

---

#### PUT /users/{username}/role
Update a user's primary role (legacy endpoint).

**Authorization:** Required (admin role)

**Request:**
```json
{
  "role": "admin"
}
```

**Response:** `200 OK`
```json
{
  "message": "User role updated successfully",
  "timestamp": 1234567890,
  "request_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

---

### Multi-Role User Management

#### GET /users/{username}/roles
Get all roles assigned to a user.

**Authorization:** Required (admin role)

**Response:** `200 OK`
```json
{
  "data": {
    "user_id": 1,
    "roles": [
      {
        "id": 1,
        "name": "admin",
        "description": "Administrator role"
      },
      {
        "id": 2,
        "name": "moderator",
        "description": "Moderator role"
      }
    ]
  },
  "message": "User roles retrieved successfully",
  "timestamp": 1234567890,
  "request_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

---

#### POST /users/{username}/roles
Add roles to a user (doesn't remove existing roles).

**Authorization:** Required (admin role)

**Request:**
```json
{
  "roles": ["moderator", "editor"]
}
```

**Response:** `200 OK`
```json
{
  "data": {
    "user_id": 1,
    "roles": [
      {
        "id": 1,
        "name": "admin",
        "description": "Administrator role"
      },
      {
        "id": 2,
        "name": "moderator",
        "description": "Moderator role"
      },
      {
        "id": 3,
        "name": "editor",
        "description": "Content editor role"
      }
    ]
  },
  "message": "Roles added successfully",
  "timestamp": 1234567890,
  "request_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

---

#### PUT /users/{username}/roles
Replace all roles for a user (removes all existing roles first).

**Authorization:** Required (admin role)

**Request:**
```json
{
  "roles": ["user", "contributor"]
}
```

**Response:** `200 OK`
```json
{
  "data": {
    "user_id": 1,
    "roles": [
      {
        "id": 4,
        "name": "user",
        "description": "Regular user role"
      },
      {
        "id": 5,
        "name": "contributor",
        "description": "Content contributor role"
      }
    ]
  },
  "message": "User roles updated successfully",
  "timestamp": 1234567890,
  "request_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

---

#### DELETE /users/{username}/roles
Remove specific roles from a user.

**Authorization:** Required (admin role)

**Request:**
```json
{
  "roles": ["moderator"]
}
```

**Response:** `200 OK`
```json
{
  "data": {
    "user_id": 1,
    "roles": [
      {
        "id": 1,
        "name": "admin",
        "description": "Administrator role"
      }
    ]
  },
  "message": "Roles removed successfully",
  "timestamp": 1234567890,
  "request_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

---

### Role Management

#### POST /roles
Create a new role.

**Authorization:** Required (admin role)

**Request:**
```json
{
  "name": "editor",
  "description": "Content editor with publishing rights"
}
```

**Response:** `201 Created`
```json
{
  "data": {
    "id": 3,
    "name": "editor",
    "description": "Content editor with publishing rights",
    "created_at": "2023-01-15T10:30:00Z"
  },
  "message": "Role created successfully",
  "timestamp": 1234567890,
  "request_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

**Errors:**
- `409 CONFLICT` - Role name already exists
- `403 FORBIDDEN` - Insufficient permissions

---

#### GET /roles
List all roles with pagination support.

**Authorization:** Required (admin role)

**Query Parameters:**
- `page` (optional): Page number, default 1
- `limit` (optional): Items per page, default 10, max 100

**Example:**
```
GET /roles?page=1&limit=20
```

**Response:** `200 OK`
```json
{
  "data": [
    {
      "id": 1,
      "name": "admin",
      "description": "Administrator role"
    },
    {
      "id": 2,
      "name": "user",
      "description": "Regular user role"
    }
  ],
  "total": 5,
  "page": 1,
  "limit": 20,
  "total_pages": 1,
  "timestamp": 1234567890
}
```

---

## Rate Limiting

API enforces rate limiting: **100 requests per minute per IP address**.

**Rate Limit Headers:**
```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1234567890
```

When exceeded:
```
429 Too Many Requests
{
  "code": "RATE_LIMIT_EXCEEDED",
  "message": "Rate limit exceeded"
}
```

---

## Request IDs

Every response includes a unique `X-Request-ID` header for tracing:

```
X-Request-ID: 550e8400-e29b-41d4-a716-446655440000
```

Use this ID when reporting issues or for troubleshooting.

---

## Password Requirements

- **Minimum Length:** 8 characters
- **Must contain:**
  - At least one uppercase letter (A-Z)
  - At least one lowercase letter (a-z)
  - At least one digit (0-9)
  - At least one special character (!@#$%^&*)

**Example:** `MyP@ssw0rd!`

---

## CAPTCHA Endpoints

#### GET /captcha/new
Get a new CAPTCHA challenge ID.

**Response:**
```json
{
  "captcha_id": "captcha_id_string"
}
```

#### GET /captcha/image/{captchaId}
Get CAPTCHA image for display.

**Response:** Binary image data (PNG/JPG)

---

## Security Headers

All responses include security headers:

```
X-Content-Type-Options: nosniff
X-Frame-Options: DENY
X-XSS-Protection: 1; mode=block
Strict-Transport-Security: max-age=31536000; includeSubDomains
```

---

## Example Usage

### Complete Login Flow

```bash
# 1. Get CAPTCHA
curl -s https://localhost:8443/api/v1/captcha/new | jq -r '.captcha_id'

# 2. Login with CAPTCHA
curl -X POST https://localhost:8443/api/v1/secure-login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "password": "MyP@ssw0rd!",
    "captchaId": "captcha_id_here",
    "captchaSolution": "solution_here"
  }' | jq '.access_token'

# 3. Use access token
curl -H "Authorization: Bearer <access_token>" \
  https://localhost:8443/api/v1/admin

# 4. Refresh token when expired
curl -X POST https://localhost:8443/api/v1/refresh \
  -H "Content-Type: application/json" \
  -d '{"refresh_token": "refresh_token_here"}' | jq '.access_token'

# 5. Logout
curl -X POST https://localhost:8443/api/v1/logout \
  -H "Content-Type: application/json" \
  -d '{"refresh_token": "refresh_token_here"}'
```

---

## Changelog

### Version 1.0.0 (Current)

**Features:**
- ✅ JWT authentication with access & refresh tokens
- ✅ Multi-role RBAC (Role-Based Access Control)
- ✅ CAPTCHA verification
- ✅ Rate limiting (100 req/min per IP)
- ✅ Soft deletes for audit trail
- ✅ Paginated role listing
- ✅ Database query optimization
- ✅ Comprehensive error handling
- ✅ Request tracking with unique IDs

---

## Support & Troubleshooting

For issues:
1. Check the error `code` and `request_id`
2. Review password requirements if registration fails
3. Ensure Bearer token format is correct: `Authorization: Bearer <token>`
4. Verify role-based access (admin role required for admin endpoints)
5. Check rate limiting if getting 429 responses

---

**Last Updated:** December 2025
**API Version:** 1.0.0
