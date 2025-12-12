# Multi-Role User Management

The Auth Service now supports assigning **multiple roles to a single user**. This provides fine-grained access control and flexible permission management.

## Overview

### Before (Single Role)
- Each user had exactly ONE role
- Role was overwritten when updated
- Limited flexibility for complex permission hierarchies

### Now (Multiple Roles)
- Each user can have MULTIPLE roles
- Roles are managed independently (add/remove/set)
- Better support for complex permission structures
- Backward compatible with existing single-role data

## Database Schema

The system uses a **many-to-many** relationship between users and roles:

```
Users Table ────────┐
                    │
                 user_roles
              (junction table)
                    │
Roles Table ────────┘
```

This allows:
- ✅ Users to have multiple roles
- ✅ Roles to be assigned to multiple users
- ✅ Easy querying of all roles for a user
- ✅ Easy querying of all users with a role

## API Endpoints

### 1. Get User's Roles

Get all roles assigned to a user.

**Endpoint:** `GET /api/v1/users/{userID}/roles`

**Auth Required:** Yes (Admin only)

**Request:**
```bash
curl -X GET http://localhost:8080/api/v1/users/1/roles \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Response (200 OK):**
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
        "id": 3,
        "name": "moderator",
        "description": "Moderator role"
      }
    ]
  },
  "message": "User roles retrieved successfully",
  "timestamp": 1702431000,
  "request_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

---

### 2. Set User's Roles (Replace All)

Replace ALL existing roles with new ones.

**Endpoint:** `PUT /api/v1/users/{userID}/roles`

**Auth Required:** Yes (Admin only)

**Request:**
```bash
curl -X PUT http://localhost:8080/api/v1/users/1/roles \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "roles": ["admin", "moderator"]
  }'
```

**Response (200 OK):**
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
        "id": 3,
        "name": "moderator",
        "description": "Moderator role"
      }
    ]
  },
  "message": "User roles updated successfully",
  "timestamp": 1702431000,
  "request_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

---

### 3. Add Roles to User (Keep Existing)

Add one or more roles to a user WITHOUT removing existing roles.

**Endpoint:** `POST /api/v1/users/{userID}/roles`

**Auth Required:** Yes (Admin only)

**Request:**
```bash
curl -X POST http://localhost:8080/api/v1/users/1/roles \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "roles": ["editor"]
  }'
```

**Response (200 OK):**
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
        "id": 3,
        "name": "moderator",
        "description": "Moderator role"
      },
      {
        "id": 4,
        "name": "editor",
        "description": "Editor role"
      }
    ]
  },
  "message": "Roles added successfully",
  "timestamp": 1702431000,
  "request_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

---

### 4. Remove Roles from User

Remove one or more roles from a user.

**Endpoint:** `DELETE /api/v1/users/{userID}/roles`

**Auth Required:** Yes (Admin only)

**Request:**
```bash
curl -X DELETE http://localhost:8080/api/v1/users/1/roles \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "roles": ["editor"]
  }'
```

**Response (200 OK):**
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
        "id": 3,
        "name": "moderator",
        "description": "Moderator role"
      }
    ]
  },
  "message": "Roles removed successfully",
  "timestamp": 1702431000,
  "request_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

---

## Use Cases

### Use Case 1: User with Multiple Permissions
```bash
# Create user
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"username": "john_admin", "password": "SecurePass123!"}'

# Set as both admin and moderator
curl -X PUT http://localhost:8080/api/v1/users/2/roles \
  -H "Authorization: Bearer ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"roles": ["admin", "moderator"]}'

# Result: john_admin has both admin and moderator permissions
```

### Use Case 2: Gradually Add Permissions
```bash
# User starts as guest
curl -X POST http://localhost:8080/api/v1/users/3/roles \
  -H "Authorization: Bearer ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"roles": ["user"]}'

# Promote to editor
curl -X POST http://localhost:8080/api/v1/users/3/roles \
  -H "Authorization: Bearer ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"roles": ["editor"]}'

# Result: user now has both user and editor roles
```

### Use Case 3: Temporary Permission Removal
```bash
# User has multiple roles
# Current: admin, moderator, editor

# Temporarily remove moderator role
curl -X DELETE http://localhost:8080/api/v1/users/1/roles \
  -H "Authorization: Bearer ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"roles": ["moderator"]}'

# Result: admin and editor roles remain
```

### Use Case 4: Complete Role Replacement
```bash
# User currently has: admin, moderator, editor

# Replace all with just user role
curl -X PUT http://localhost:8080/api/v1/users/1/roles \
  -H "Authorization: Bearer ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"roles": ["user"]}'

# Result: Only user role remains (all others removed)
```

---

## Role Checking in Code

### Check if user has specific role

```go
// Get user with roles
user, _ := database.DB.GetUserByID(userID)

// Check if user has a specific role
hasAdminRole := false
for _, role := range user.Roles {
    if role.Name == "admin" {
        hasAdminRole = true
        break
    }
}

if hasAdminRole {
    // User is admin
}
```

### Check if user has any of multiple roles

```go
func UserHasAnyRole(user *User, roleNames []string) bool {
    for _, userRole := range user.Roles {
        for _, requiredRole := range roleNames {
            if userRole.Name == requiredRole {
                return true
            }
        }
    }
    return false
}

// Usage
if UserHasAnyRole(user, []string{"admin", "moderator"}) {
    // User has admin OR moderator role
}
```

### Check if user has all required roles

```go
func UserHasAllRoles(user *User, roleNames []string) bool {
    userRoleMap := make(map[string]bool)
    for _, role := range user.Roles {
        userRoleMap[role.Name] = true
    }

    for _, requiredRole := range roleNames {
        if !userRoleMap[requiredRole] {
            return false
        }
    }
    return true
}

// Usage
if UserHasAllRoles(user, []string{"admin", "moderator"}) {
    // User has BOTH admin AND moderator roles
}
```

---

## Middleware: Role-Based Access Control

The existing role middleware checks for at least one role:

```go
// Checks if user has the "admin" role
middleware.RoleMiddleware("admin")

// In the future, you could extend this for multiple roles:
// middleware.RoleMiddleware("admin|moderator")  // Has admin OR moderator
// middleware.RoleMiddleware("admin&moderator")  // Has admin AND moderator
```

---

## JWT Token Claims

When a token is issued, it includes roles information:

```json
{
  "user": 1,
  "role": "admin",
  "exp": 1702517400,
  "iat": 1702431000
}
```

### Future Enhancement: Multiple Roles in JWT

You could enhance this to include all roles:

```json
{
  "user": 1,
  "roles": ["admin", "moderator", "editor"],
  "exp": 1702517400,
  "iat": 1702431000
}
```

---

## Error Handling

### Invalid User ID
```bash
curl -X GET http://localhost:8080/api/v1/users/abc/roles \
  -H "Authorization: Bearer TOKEN"
```

**Response (400 Bad Request):**
```json
{
  "code": "INVALID_JSON",
  "message": "Invalid user ID format",
  "timestamp": 1702431000,
  "request_id": "..."
}
```

### User Not Found
```bash
curl -X GET http://localhost:8080/api/v1/users/9999/roles \
  -H "Authorization: Bearer TOKEN"
```

**Response (404 Not Found):**
```json
{
  "code": "USER_NOT_FOUND",
  "message": "User not found",
  "timestamp": 1702431000,
  "request_id": "..."
}
```

### Invalid Role Name
```bash
curl -X POST http://localhost:8080/api/v1/users/1/roles \
  -H "Authorization: Bearer TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"roles": ["invalid_role"]}'
```

**Response (500 Internal Server Error):**
```json
{
  "code": "INTERNAL_SERVER_ERROR",
  "message": "Failed to add roles: role not found: invalid_role",
  "timestamp": 1702431000,
  "request_id": "..."
}
```

---

## Data Migration

If migrating from single-role to multi-role system:

### Before (Single Role)
```go
type User struct {
    ID     uint
    Username string
    RoleID uint
    Role   Role
}
```

### After (Multiple Roles)
```go
type User struct {
    ID    uint
    Username string
    Roles []Role `gorm:"many2many:user_roles;"`
}
```

GORM handles the migration automatically:
1. Creates `user_roles` junction table
2. Migrates existing `RoleID` values to the junction table
3. Both systems work during transition

---

## Testing

### Test Script

```bash
#!/bin/bash

USER_ID=1
ADMIN_TOKEN="YOUR_JWT_TOKEN"

# Get current roles
echo "=== Get User Roles ==="
curl -s -X GET http://localhost:8080/api/v1/users/$USER_ID/roles \
  -H "Authorization: Bearer $ADMIN_TOKEN" | jq .

# Set roles
echo -e "\n=== Set Roles (admin, moderator) ==="
curl -s -X PUT http://localhost:8080/api/v1/users/$USER_ID/roles \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"roles": ["admin", "moderator"]}' | jq .

# Add role
echo -e "\n=== Add Role (editor) ==="
curl -s -X POST http://localhost:8080/api/v1/users/$USER_ID/roles \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"roles": ["editor"]}' | jq .

# Remove role
echo -e "\n=== Remove Role (moderator) ==="
curl -s -X DELETE http://localhost:8080/api/v1/users/$USER_ID/roles \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"roles": ["moderator"]}' | jq .

# Get final roles
echo -e "\n=== Get Final Roles ==="
curl -s -X GET http://localhost:8080/api/v1/users/$USER_ID/roles \
  -H "Authorization: Bearer $ADMIN_TOKEN" | jq .
```

---

## Performance Considerations

### Database Queries
- `GetUserRoles(userID)` - Single query with Preload
- `AddRolesToUser()` - Finds user and roles, then appends
- `RemoveRolesFromUser()` - Finds user and roles, then deletes
- `SetUserRoles()` - Finds user and roles, then replaces

### Optimization Tips
1. **Batch Operations**: If assigning roles to many users, consider batch processing
2. **Caching**: Cache role assignments for frequently accessed users
3. **Eager Loading**: Always use `Preload("Roles")` when fetching users
4. **Indexing**: Junction table is indexed on both user_id and role_id

---

## Future Enhancements

- ✅ Multi-role support (DONE)
- ⏳ JWT token with all roles
- ⏳ Role-based route middleware for multiple roles
- ⏳ Permission-based access control (PBAC)
- ⏳ Role hierarchy (parent/child roles)
- ⏳ Temporal role assignment (roles expire after X days)
- ⏳ Audit trail for role changes

---

## Summary

The multi-role system provides:

| Feature | Benefit |
|---------|---------|
| **Multiple Roles per User** | Flexible permission management |
| **Add/Remove/Set Operations** | Fine-grained control |
| **Backward Compatible** | Works with existing single-role data |
| **Clean API** | Intuitive REST endpoints |
| **Error Handling** | Clear error messages |
| **Logging** | All changes logged with request IDs |
