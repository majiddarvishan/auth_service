   # SHPD Router Web Service API Documentation
   ## 📋 Table of Contents - Complete Documentation

   - [Overview](#overview)
   - [Quick Start](#quick-start)
   - [Authentication](#authentication)
     - [JWT Token Usage](#jwt-token-usage)
     - [Token Expiration & Refresh](#token-expiration--refresh)
   - [Rate Limiting](#rate-limiting)
   - [API Versioning](#api-versioning)
   - [Common Error Codes](#common-error-codes)
   - [Endpoints](#endpoints)
     - [🔐 Authentication Endpoints](#authentication-endpoints)
     - [👤 User Management Endpoints](#user-management-endpoints)
     - [🎭 Role Management Endpoints](#role-management-endpoints)
     - [⚙️ Admin Endpoints](#admin-endpoints)
     - [🤖 CAPTCHA Endpoints](#captcha-endpoints)
     - [📞 MSISDN Mapping Endpoints](#msisdn-mapping-endpoints)
     - [📱 SMS Endpoints](#sms-endpoints)
     - [🏷️ Prefix Management Endpoints](#prefix-management-endpoints)
   - [Code Examples](#code-examples)
     - [Python Examples](#python-examples)
     - [JavaScript Examples](#javascript-examples)
   - [Sequence Diagrams](#sequence-diagrams)
   - [Best Practices](#best-practices)
   - [Try It Out](#try-it-out)
   - [Postman Collection](#postman-collection)
   - [Support & License](#Support-License)

   ## Overview

   This is a comprehensive API documentation covering three main services:

   1. **Authentication Service** - User management, role-based access control, and secure login with CAPTCHA
   2. **MSISDN Mapping Service** - Phone number mapping and replacement management
   3. **SMS Service** - SMS sending, delivery tracking, and prefix management

   **Base URL:** `http://localhost:8080/api/v1`
      
   **Version:** 1.0

   **Contact:** Majid Darvishan (m.darvishan@shpdco.ir)

---

   ## Quick Start

   ```bash
   # 1. Register a user
   curl -X POST http://localhost:8080/v1/api/users \
     -H "Content-Type: application/json" \
     -d '{"username": "testuser", "password": "Pass123!", "role": "user"}'
   
   # 2. Login to get JWT token
   curl -X POST http://localhost:8080/v1/api/login \
     -H "Content-Type: application/json" \
     -d '{"username": "testuser", "password": "Pass123!"}'
   
   # 3. Use token for authenticated requests
   curl -X GET http://localhost:8080/v1/api/admin/dashboard \
     -H "Authorization: Bearer YOUR_JWT_TOKEN"
   ```

---

   ## Authentication

   ### JWT Token Usage

   Most endpoints require JWT (JSON Web Token) authentication. After successful login, you'll receive a token that must be included in subsequent requests.

   **Header Format:**
   ```
   Authorization: Bearer <your-jwt-token>
   ```

   **Token Structure:**
   ```
   eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.payload.signature
   ```

   ### Token Expiration & Refresh

   - **Default Expiration:** Tokens typically expire after 24 hours (check your server configuration)
   - **Refresh Strategy:** Re-authenticate using `/login` or `/secure-login` endpoints when token expires
   - **Best Practice:** Implement automatic token refresh in your client application before expiration

   **Handling Expired Tokens:**
   ```javascript
   // JavaScript example
   if (response.status === 401) {
     // Token expired, re-authenticate
     await refreshToken();
   }
   ```

---

   ## Rate Limiting

   To ensure fair usage and system stability, the following rate limits apply:

   | Endpoint Category | Rate Limit | Window |
   |------------------|------------|--------|
   | Authentication (login) | 5 requests | 1 minute |
   | SMS Sending | 100 requests | 1 minute |
   | Status Queries | 200 requests | 1 minute |
   | Other Endpoints | 60 requests | 1 minute |

   **Rate Limit Headers:**
   ```
   X-RateLimit-Limit: 100
   X-RateLimit-Remaining: 95
   X-RateLimit-Reset: 1705315800
   ```

   **429 Response (Rate Limited):**
   ```json
   {
     "error": "Rate limit exceeded",
     "retry_after": 60
   }
   ```

---

   ## API Versioning

   **Current Version:** v1

   **Version Strategy:** URL-based versioning
   - Auth/MSISDN/SMS: `/api/v1/...`

   **Backward Compatibility:**
   - Minor changes (bug fixes, new optional fields) are backward compatible
   - Breaking changes will be introduced in new versions (v2, v3, etc.)
   - Deprecated endpoints will be supported for at least 6 months with advance notice

   **Checking API Version:**
   ```bash
   curl -X GET http://localhost:8080/api/v1/version
   ```

---

   ## Common Error Codes

   ### HTTP Status Codes

   | Code | Meaning | Description |
   |------|---------|-------------|
   | 200 | OK | Request successful |
   | 201 | Created | Resource created successfully |
   | 400 | Bad Request | Invalid input or malformed request |
   | 401 | Unauthorized | Missing or invalid authentication |
   | 403 | Forbidden | Authenticated but insufficient permissions |
   | 404 | Not Found | Resource not found |
   | 409 | Conflict | Resource already exists or conflict detected |
   | 429 | Too Many Requests | Rate limit exceeded |
   | 500 | Internal Server Error | Server-side error |

   ### Application Error Codes

   | Code | Category | Description | Solution |
   |------|----------|-------------|----------|
   | 1001 | Authentication | Invalid credentials | Check username/password |
   | 1002 | Authentication | Token expired | Re-authenticate |
   | 1003 | Authentication | Invalid token | Obtain new token |
   | 2001 | Validation | Missing required field | Check request body |
   | 2002 | Validation | Invalid format | Verify data format |
   | 2003 | Validation | Value out of range | Check allowed values |
   | 3001 | SMS | Invalid phone number | Verify number format |
   | 3002 | SMS | Message too long | Split or shorten message |
   | 3003 | SMS | Delivery failed | Check number status |
   | 4001 | Prefix | Duplicate prefix | Use different prefix |
   | 4002 | Prefix | Invalid prefix format | Check prefix rules |
   | 5001 | System | Database error | Contact support |
   | 5002 | System | Service unavailable | Retry later |

   ### Error Response Format

   **Standard Error Response:**
   ```json
   {
     "error": "Error message",
     "details": "Detailed explanation",
     "code": 2001,
     "timestamp": "2026-01-15T10:30:00Z"
   }
   ```

   **SMS Service Error Response:**
   ```json
   {
     "code": 3001,
     "message": "Invalid phone number format"
   }
   ```

---

   ## Best Practices

   ### 🔒 Security Best Practices

   1. **Token Management**
      - Store JWT tokens securely (use httpOnly cookies in browsers)
      - Never expose tokens in URLs or logs
      - Implement token refresh before expiration
      - Clear tokens on logout
   
   2. **Password Security**
      - Enforce strong password requirements (min 8 chars, uppercase, lowercase, numbers, special chars)
      - Never log or display passwords
      - Use secure login endpoint for public-facing applications
      - Implement account lockout after failed attempts
   
   3. **API Key Protection**
      - Never commit API keys or tokens to version control
      - Use environment variables for sensitive configuration
      - Rotate credentials regularly
      - Use different credentials for development/production
   
   4. **Input Validation**
      - Always validate and sanitize user input
      - Use parameterized queries to prevent injection
      - Validate phone numbers before SMS sending
      - Check message length limits (160 chars for standard SMS)

   ### 📱 SMS Best Practices

   1. **Message Optimization**
      - Keep messages under 160 characters when possible
      - Use clear, concise language
      - Include opt-out instructions for marketing messages
      - Test messages before bulk sending
   
   2. **Delivery Tracking**
      - Always use track_ids for important messages
      - Implement retry logic for failed deliveries
      - Monitor delivery rates and adjust strategy
      - Store group_id for bulk sends to track batches
   
   3. **Prefix Management**
      - Register all sender numbers/prefixes before use
      - Use descriptive prefixes for easy identification
      - Regularly audit and clean up unused prefixes
      - Test with prefixes before production use
   
   4. **Rate Management**
      - Respect rate limits (100 SMS/minute)
      - Implement exponential backoff for rate limit errors
      - Use bulk endpoints for multiple messages
      - Spread large batches over time

   ### ⚡ Performance Best Practices

   1. **API Usage**
      - Use bulk endpoints when sending multiple SMS
      - Implement connection pooling for high-volume applications
      - Cache frequently accessed data (roles, prefixes)
      - Use pagination for large result sets
   
   2. **Error Handling**
      - Always check response status codes
      - Implement proper error handling and logging
      - Use exponential backoff for retries
      - Log failed requests for debugging
   
   3. **Network Optimization**
      - Use compression for large payloads
      - Implement timeouts for all requests (30s recommended)
      - Reuse HTTP connections
      - Handle network failures gracefully

   ### 📊 Monitoring and Logging

   1. **Application Monitoring**
      - Track API response times
      - Monitor error rates by endpoint
      - Set up alerts for rate limit warnings
      - Log all authentication failures
   
   2. **SMS Monitoring**
      - Track delivery success rates
      - Monitor pending message queues
      - Alert on high failure rates
      - Track costs per user/campaign
   
   3. **Security Monitoring**
      - Log all admin actions
      - Monitor for unusual authentication patterns
      - Track failed login attempts
      - Alert on potential security incidents

   ### 🔄 Integration Best Practices

   1. **Development Workflow**
      - Use separate environments (dev, staging, production)
      - Test thoroughly in staging before production
      - Implement automated testing for critical flows
      - Use feature flags for gradual rollouts
   
   2. **Error Recovery**
      - Implement idempotent operations where possible
      - Store message metadata for replay capability
      - Handle partial failures in bulk operations
      - Provide clear error messages to users
   
   3. **Documentation**
      - Keep internal API documentation updated
      - Document custom workflows and integrations
      - Maintain changelog for API changes
      - Share best practices with your team

---

   ## Endpoints

   ### 🔐 Authentication Endpoints

   #### 1. Register a New User

   `POST /users`

   Creates a new user account with specified username, password, and role.

   **Request Model:**
   ```json
   {
     "username": "string",
     "password": "string",
     "role": "string"
   }
   ```

   **Request Example:**
   ```json
   {
     "username": "john_doe",
     "password": "SecurePass123!",
     "role": "user"
   }
   ```

   **Response Model (200 OK):**
   ```json
   {
     "message": "string"
   }
   ```

   **Response Example:**
   ```json
   {
     "message": "User registered successfully"
   }
   ```

   **Error Responses:**
   - `400`: Invalid input or missing fields
   - `500`: Server error during registration

---

   #### 2. Login (Standard)

   `POST /login`

   Authenticate with username and password to receive a JWT token.

   **Request Model:**
   ```json
   {
     "username": "string",
     "password": "string"
   }
   ```

   **Request Example:**
   ```json
   {
     "username": "john_doe",
     "password": "SecurePass123!"
   }
   ```

   **Response Model (200 OK):**
   ```json
   {
     "token": "string"
   }
   ```

   **Response Example:**
   ```json
   {
     "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
   }
   ```

   **Error Responses:**
   - `400`: Invalid JSON format
   - `401`: Unauthorized - invalid credentials
   - `500`: Server error during token generation

---

   #### 3. Secure Login (with CAPTCHA)

   `POST /secure-login`

   Enhanced login endpoint that requires CAPTCHA verification for additional security.

   **Request Model:**
   ```json
   {
     "username": "string",
     "password": "string",
     "captchaId": "string",
     "captchaSolution": "string"
   }
   ```

   **Request Example:**
   ```json
   {
     "username": "john_doe",
     "password": "SecurePass123!",
     "captchaId": "a1b2c3",
     "captchaSolution": "XYZABC"
   }
   ```

   **Response Model (200 OK):**
   ```json
   {
     "token": "string"
   }
   ```

   **Response Example:**
   ```json
   {
     "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
   }
   ```

   **Error Responses:**
   - `400`: Invalid JSON format or incorrect CAPTCHA
   - `401`: Unauthorized - invalid credentials
   - `500`: Server error during token generation

---

   ### 👤 User Management Endpoints

   #### 4. Update User Role

   `PUT /users/{username}/role`

   Update the role of an existing user. Admin access required.

   **Path Parameters:**
   - `username` (string, required): The username to update

   **Request Model:**
   ```json
   {
     "role": "string"
   }
   ```

   **Request Example:**
   ```json
   {
     "role": "admin"
   }
   ```

   **Response Model (200 OK):**
   ```json
   {
     "message": "string"
   }
   ```

   **Response Example:**
   ```json
   {
     "message": "User role updated successfully"
   }
   ```

   **Error Responses:**
   - `400`: Invalid input or missing fields
   - `404`: User not found
   - `500`: Failed to update user role

---

   #### 5. Delete User

   `DELETE /users/{username}`

   Delete an existing user account. Admin access required.

   **Path Parameters:**
   - `username` (string, required): The username to delete

   **Response Model (200 OK):**
   ```json
   {
     "message": "string"
   }
   ```

   **Response Example:**
   ```json
   {
     "message": "User deleted successfully"
   }
   ```

   **Error Responses:**
   - `400`: Username is required
   - `404`: User not found
   - `500`: Could not delete user

---

   ### 🎭 Role Management Endpoints

   #### 6. List All Roles

   `GET /roles`

   Retrieves a list of all roles defined in the system.

   **Response Model (200 OK):**
   ```json
   {
     "roles": [
       {
         "id": "integer",
         "name": "string",
         "description": "string"
       }
     ]
   }
   ```

   **Response Example:**
   ```json
   {
     "roles": [
       {
         "id": 1,
         "name": "admin",
         "description": "Administrator with full access"
       },
       {
         "id": 2,
         "name": "user",
         "description": "Standard user with limited access"
       }
     ]
   }
   ```

   **Error Responses:**
   - `500`: Internal server error

---

   #### 7. Create New Role

   `POST /roles`

   Create a new role with specified name and description. Admin access required.

   **Request Model:**
   ```json
   {
     "name": "string",
     "description": "string"
   }
   ```

   **Request Example:**
   ```json
   {
     "name": "moderator",
     "description": "Moderator with content management permissions"
   }
   ```

   **Response Model (200 OK):**
   ```json
   {
     "message": "string",
     "role": {
       "id": "integer",
       "name": "string",
       "description": "string"
     }
   }
   ```

   **Response Example:**
   ```json
   {
     "message": "Role created successfully",
     "role": {
       "id": 3,
       "name": "moderator",
       "description": "Moderator with content management permissions"
     }
   }
   ```

   **Error Responses:**
   - `400`: Bad request - invalid input
   - `500`: Internal server error

---

   ### ⚙️ Admin Endpoints

   #### 8. Admin Dashboard

   `GET /admin/dashboard`

   Retrieves comprehensive dashboard data including all users and accounting rules. Admin access required.

   **Headers:**
   - `Authorization: Bearer <jwt-token>`

   **Response Model (200 OK):**
   ```json
   {
     "message": "string",
     "users": [
       {
         "username": "string",
         "role": "string"
       }
     ],
     "rules": [
       {
         "id": "integer",
         "endpoint": "string",
         "charge": "number",
         "created_at": "string",
         "updated_at": "string"
       }
     ]
   }
   ```

   **Response Example:**
   ```json
   {
     "message": "Dashboard data retrieved successfully",
     "users": [
       {
         "username": "john_doe",
         "role": "user"
       },
       {
         "username": "admin_user",
         "role": "admin"
       }
     ],
     "rules": [
       {
         "id": 1,
         "endpoint": "/api/premium-feature",
         "charge": 0.50,
         "created_at": "2024-01-15T10:30:00Z",
         "updated_at": "2024-01-15T10:30:00Z"
       }
     ]
   }
   ```

   **Error Responses:**
   - `401`: Token claims missing or invalid
   - `403`: Not an admin
   - `500`: Database error

---

   #### 9. Create Custom Endpoint

   `POST /admin/custom-endpoints`

   Create custom endpoints to redirect requests to other endpoints. Admin access required.

   **Request Model:**
   ```json
   {
     "path": "string",
     "method": "string",
     "endpoints": ["string"],
     "needAccounting": "boolean"
   }
   ```

   **Request Example:**
   ```json
   {
     "path": "/custom/feature",
     "method": "GET",
     "endpoints": [
       "/internal/service1",
       "/internal/service2"
     ],
     "needAccounting": true
   }
   ```

   **Response Model (200 OK):**
   ```json
   {
     "message": "string"
   }
   ```

   **Response Example:**
   ```json
   {
     "message": "Custom endpoint created successfully"
   }
   ```

   **Error Responses:**
   - `400`: Invalid JSON format
   - `401`: Unauthorized
   - `500`: Server error

---

   ### 🤖 CAPTCHA Endpoints

   #### 10. Generate New CAPTCHA

   `GET /captcha/new`

   Generates a new CAPTCHA ID for verification purposes.

   **Response Model (200 OK):**
   ```json
   {
     "captchaId": "string"
   }
   ```

   **Response Example:**
   ```json
   {
     "captchaId": "a1b2c3"
   }
   ```

---

   #### 11. Get CAPTCHA Image

   `GET /captcha/image/{captchaId}`

   Retrieves the CAPTCHA image in PNG format for the given CAPTCHA ID.

   **Path Parameters:**
   - `captchaId` (string, required): The CAPTCHA ID to retrieve

   **Example:** `GET /captcha/image/a1b2c3`

   **Response (200 OK):**
   - Content-Type: `image/png`
   - Binary image data

   **Error Responses:**
   - `500`: Internal server error

---

   ### 📞 MSISDN Mapping Endpoints

   MSISDN (Mobile Station International Subscriber Directory Number) mapping allows you to replace phone numbers in the system.

   #### 12. List All MSISDN Mappings

   `GET /msisdn-mappings`

   Retrieves a list of all MSISDN mappings with optional pagination support.

   **Query Parameters:**
   - `begin` (integer, optional, default: 0): Start index for pagination
   - `limit` (integer, optional, default: 100): Maximum number of results to return

   **Response Model (200 OK):**
   ```json
   [
     {
       "msisdn": "string",
       "new_msisdn": "string",
       "created_at": "string"
     }
   ]
   ```

   **Response Example:**
   ```json
   [
     {
       "msisdn": "1234567890",
       "new_msisdn": "0987654321",
       "created_at": "2024-01-15T10:30:00Z"
     }
   ]
   ```

   **Error Responses:**
   - `400`: Bad request - invalid pagination parameters

---

   #### 13. Get Specific MSISDN Mapping

   `GET /msisdn-mappings/{msisdn}`

   Retrieves details for a specific MSISDN mapping.

   **Path Parameters:**
   - `msisdn` (string, required): The original MSISDN to look up

   **Response Model (200 OK):**
   ```json
   {
     "msisdn": "string",
     "new_msisdn": "string",
     "created_at": "string"
   }
   ```

   **Response Example:**
   ```json
   {
     "msisdn": "1234567890",
     "new_msisdn": "0987654321",
     "created_at": "2024-01-15T10:30:00Z"
   }
   ```

   **Error Responses:**
   - `404`: MSISDN mapping not found

---

   #### 14. Add New MSISDN Mapping

   `POST /msisdn-mappings`

   Creates a new MSISDN mapping to replace an old phone number with a new one.

   **Request Model:**
   ```json
   {
     "msisdn": "string",
     "new_msisdn": "string",
     "created_at": "string"
   }
   ```

   **Request Example:**
   ```json
   {
     "msisdn": "1234567890",
     "new_msisdn": "0987654321",
     "created_at": "2024-01-15T10:30:00Z"
   }
   ```

   **Response Model (201 Created):**
   ```json
   {
     "msisdn": "string",
     "new_msisdn": "string",
     "created_at": "string"
   }
   ```

   **Response Example:**
   ```json
   {
     "msisdn": "1234567890",
     "new_msisdn": "0987654321",
     "created_at": "2024-01-15T10:30:00Z"
   }
   ```

   **Error Responses:**
   - `400`: Bad request - invalid input data
   - `409`: Conflict - MSISDN mapping already exists

---

   #### 15. Remove MSISDN Mapping

   `DELETE /msisdn-mappings/{msisdn}`

   Removes an existing MSISDN mapping so the number will no longer be replaced.

   **Path Parameters:**
   - `msisdn` (string, required): The original MSISDN to remove

   **Response Model (200 OK):**
   ```json
   {
     "message": "string"
   }
   ```

   **Response Example:**
   ```json
   {
     "message": "MSISDN mapping removed successfully"
   }
   ```

   **Error Responses:**
   - `404`: MSISDN mapping not found

---

   ### 📱 SMS Endpoints

   The SMS service allows you to send text messages, track delivery status, and manage sender prefixes. All SMS endpoints require a `user-name` header for authentication.

   #### 16. Send Simple SMS

   `POST /sms/send/simple`

   Sends a single SMS message to a single receiver.

   **Headers:**
   - `user-name` (string, required): Authenticated username

   **Request Model:**
   ```json
   {
     "sender": "string",
     "receiver": "string",
     "message": "string"
   }
   ```

   **Request Example:**
   ```json
   {
     "sender": "1234567890",
     "receiver": "0987654321",
     "message": "Hello, this is a test message!"
   }
   ```

   **Response Model (200 OK):**
   ```json
   {
     "group_id": "string",
     "status": [
       {
         "message_id": "string",
         "track_id": "string",
         "receiver": "string",
         "sender": "string",
         "received_time": "string",
         "total_parts": "integer",
         "error": {
           "code": "integer",
           "message": "string"
         }
       }
     ]
   }
   ```

   **Response Example:**
   ```json
   {
     "group_id": "grp_abc123",
     "status": [
       {
         "message_id": "msg_12345",
         "track_id": "trk_67890",
         "receiver": "0987654321",
         "sender": "1234567890",
         "received_time": "2024-01-15T10:30:00Z",
         "total_parts": 1,
         "error": null
       }
     ]
   }
   ```

   **Error Responses:**
   - `400`: Bad request - invalid input
   - `500`: Internal server error

---

   #### 17. Send Bulk SMS

   `POST /sms/send/bulk`

   Sends multiple SMS messages to multiple receivers with support for different senders and tracking IDs.

   **Headers:**
   - `user-name` (string, required): Authenticated username

   **Request Model:**
   ```json
   {
     "sender": ["string"],
     "receiver": ["string"],
     "message": ["string"],
     "track_id": ["string"]
   }
   ```

   **Request Example:**
   ```json
   {
     "sender": ["1234567890", "1234567890"],
     "receiver": ["0987654321", "5551234567"],
     "message": ["Hello User 1!", "Hello User 2!"],
     "track_id": ["trk_001", "trk_002"]
   }
   ```

   **Response Model (200 OK):**
   ```json
   {
     "group_id": "string",
     "status": [
       {
         "message_id": "string",
         "track_id": "string",
         "receiver": "string",
         "sender": "string",
         "received_time": "string",
         "total_parts": "integer",
         "error": {
           "code": "integer",
           "message": "string"
         }
       }
     ]
   }
   ```

   **Response Example:**
   ```json
   {
     "group_id": "grp_abc456",
     "status": [
       {
         "message_id": "msg_12345",
         "track_id": "trk_001",
         "receiver": "0987654321",
         "sender": "1234567890",
         "received_time": "2024-01-15T10:30:00Z",
         "total_parts": 1,
         "error": null
       },
       {
         "message_id": "msg_12346",
         "track_id": "trk_002",
         "receiver": "5551234567",
         "sender": "1234567890",
         "received_time": "2024-01-15T10:30:05Z",
         "total_parts": 1,
         "error": null
       }
     ]
   }
   ```

   **Error Responses:**
   - `400`: Bad request - invalid input
   - `500`: Internal server error

---

   #### 18. Get SMS Status

   `POST /sms/status`

   Retrieves delivery status for specific SMS messages by message ID or track ID.

   **Headers:**
   - `user-name` (string, required): Authenticated username

   **Request Model:**
   ```json
   [
     {
       "message_id": "string",
       "track_id": "string"
     }
   ]
   ```

   **Request Example:**
   ```json
   [
     {
       "message_id": "msg_12345",
       "track_id": ""
     },
     {
       "message_id": "",
       "track_id": "trk_001"
     }
   ]
   ```

   **Response Model (200 OK):**
   ```json
   [
     {
       "message_id": "string",
       "track_id": "string",
       "overall_status": "PENDING | DELIVRD | FAILED | INVALID | INTERNAL_ERROR"
     }
   ]
   ```

   **Response Example:**
   ```json
   [
     {
       "message_id": "msg_12345",
       "track_id": "trk_001",
       "overall_status": "DELIVRD"
     }
   ]
   ```

   **SMS Status Values:**
   - `PENDING`: Message is queued or being sent
   - `DELIVRD`: Message delivered successfully
   - `FAILED`: Message delivery failed
   - `INVALID`: Invalid message format or parameters
   - `INTERNAL_ERROR`: System error occurred

   **Error Responses:**
   - `400`: Bad request - invalid input
   - `404`: Message not found

---

   #### 19. Get Bulk SMS Status

   `GET /sms/status/bulk`

   Retrieves delivery status for a batch of messages using a group ID.

   **Headers:**
   - `user-name` (string, required): Authenticated username

   **Query Parameters:**
   - `models.QueryParamGroupID` (string, required): The group ID returned from bulk send

   **Response Model (200 OK):**
   ```json
   {
     "status": [
       {
         "message_id": "string",
         "track_id": "string",
         "overall_status": "PENDING | DELIVRD | FAILED | INVALID | INTERNAL_ERROR"
       }
     ]
   }
   ```

   **Response Example:**
   ```json
   {
     "status": [
       {
         "message_id": "msg_12345",
         "track_id": "trk_001",
         "overall_status": "DELIVRD"
       }
     ]
   }
   ```

   **Error Responses:**
   - `400`: Bad request - invalid group ID
   - `500`: Internal server error

---

   #### 20. Get Received Messages

   `POST /sms/received-message`

   Retrieves all received SMS messages for specified receivers within a time range.

   **Headers:**
   - `user-name` (string, required): Authenticated username

   **Request Model:**
   ```json
   {
     "receiver": ["string"],
     "epoch_time": "string"
   }
   ```

   **Request Example:**
   ```json
   {
     "receiver": ["1234567890", "0987654321"],
     "epoch_time": "1705315800"
   }
   ```

   **Response Model (200 OK):**
   ```json
   [
     {
       "sender": "string",
       "receiver": "string",
       "message": "string",
       "received_time": "string"
     }
   ]
   ```

   **Response Example:**
   ```json
   [
     {
       "sender": "5551234567",
       "receiver": "1234567890",
       "message": "Thank you for your message!",
       "received_time": "2024-01-15T10:35:00Z"
     }
   ]
   ```

   **Error Responses:**
   - `400`: Bad request - invalid input
   - `500`: Internal server error

---

   ### 🏷️ Prefix Management Endpoints

   Prefix management allows you to configure sender number prefixes for users.

   #### 21. Add Prefixes

   `POST /prefix/add`

   Adds new sender prefixes for a user.

   **Headers:**
   - `user-name` (string, required): Authenticated username

   **Request Model:**
   ```json
   {
     "user": "string",
     "numbers": [
       {
         "number": "string",
         "is_prefix": "boolean"
       }
     ]
   }
   ```

   **Request Example:**
   ```json
   {
     "user": "john_doe",
     "numbers": [
       {
         "number": "1234",
         "is_prefix": true
       }
     ]
   }
   ```

   **Response Model (200 OK):**
   ```json
   [
     {
       "prefix": "string",
       "error": {
         "code": "integer",
         "message": "string"
       }
     }
   ]
   ```

   **Response Example:**
   ```json
   [
     {
       "prefix": "1234",
       "error": null
     }
   ]
   ```

   **Error Responses:**
   - `400`: Bad request - invalid input format

---

   #### 22. Remove Prefixes

   `DELETE /prefix/remove`

   Removes one or more prefixes for a user.

   **Headers:**
   - `user-name` (string, required): Authenticated username

   **Request Model:**
   ```json
   {
     "user": "string",
     "numbers": ["string"]
   }
   ```

   **Request Example:**
   ```json
   {
     "user": "john_doe",
     "numbers": ["1234", "5678"]
   }
   ```

   **Response Model (200 OK):**
   ```json
   [
     {
       "prefix": "string",
       "error": null
     }
   ]
   ```

   **Error Responses:**
   - `400`: Bad request - invalid input

---

   #### 23. Search Prefixes

   `GET /prefix/search`

   Searches for prefixes by username, prefix value, or returns all prefixes with pagination.

   **Headers:**
   - `user-name` (string, required): Authenticated username

   **Query Parameters:**
   - `username` (string, optional): Comma-separated usernames to filter
   - `prefix` (string, optional): Comma-separated prefixes to filter
   - `page` (integer, optional, default: 1): Page number
   - `page_size` (integer, optional, default: 50, max: 100): Results per page

   **Response Model (200 OK):**
   ```json
   [
     {
       "user": "string",
       "prefix": "string",
       "is_prefix": "boolean",
       "create_at": "string"
     }
   ]
   ```

   **Response Example:**
   ```json
   [
     {
       "user": "john_doe",
       "prefix": "1234",
       "is_prefix": true,
       "create_at": "2024-01-15T10:30:00Z"
     }
   ]
   ```

   **Error Responses:**
   - `400`: Bad request - invalid parameters
   - `401`: Unauthorized
   - `500`: Internal server error

---

   ## Python Examples

   ### Authentication and User Registration

   ```python
   import requests
   import json
   
   BASE_URL = "http://localhost:8080/v1/api"
   
   # Register a new user
   def register_user(username, password, role="user"):
       url = f"{BASE_URL}/users"
       payload = {
           "username": username,
           "password": password,
           "role": role
       }
       response = requests.post(url, json=payload)
       return response.json()
   
   # Login and get JWT token
   def login(username, password):
       url = f"{BASE_URL}/login"
       payload = {
           "username": username,
           "password": password
       }
       response = requests.post(url, json=payload)
       if response.status_code == 200:
           return response.json()["token"]
       else:
           raise Exception(f"Login failed: {response.text}")
   
   # Secure login with CAPTCHA
   def secure_login(username, password, captcha_id, captcha_solution):
       url = f"{BASE_URL}/secure-login"
       payload = {
           "username": username,
           "password": password,
           "captchaId": captcha_id,
           "captchaSolution": captcha_solution
       }
       response = requests.post(url, json=payload)
       return response.json()
   
   # Get CAPTCHA
   def get_captcha():
       url = f"{BASE_URL}/captcha/new"
       response = requests.get(url)
       return response.json()["captchaId"]
   
   # Example usage
   try:
       # Register
       result = register_user("john_doe", "SecurePass123!")
       print(f"Registration: {result}")
       
       # Login
       token = login("john_doe", "SecurePass123!")
       print(f"Token: {token}")
       
       # Use token for authenticated request
       headers = {"Authorization": f"Bearer {token}"}
       dashboard = requests.get(f"{BASE_URL}/admin/dashboard", headers=headers)
       print(f"Dashboard: {dashboard.json()}")
       
   except Exception as e:
       print(f"Error: {e}")
   ```

---

   ### SMS Sending with Python

   ```python
   import requests
   
   SMS_BASE_URL = "http://10.10.10.21:8080/api/v1"
   
   class SMSClient:
       def __init__(self, username):
           self.username = username
           self.headers = {"user-name": username, "Content-Type": "application/json"}
       
       def send_simple_sms(self, sender, receiver, message):
           """Send a single SMS"""
           url = f"{SMS_BASE_URL}/sms/send/simple"
           payload = {
               "sender": sender,
               "receiver": receiver,
               "message": message
           }
           response = requests.post(url, json=payload, headers=self.headers)
           return response.json()
       
       def send_bulk_sms(self, senders, receivers, messages, track_ids=None):
           """Send multiple SMS messages"""
           url = f"{SMS_BASE_URL}/sms/send/bulk"
           payload = {
               "sender": senders,
               "receiver": receivers,
               "message": messages,
               "track_id": track_ids or [f"trk_{i}" for i in range(len(messages))]
           }
           response = requests.post(url, json=payload, headers=self.headers)
           return response.json()
       
       def check_sms_status(self, message_ids=None, track_ids=None):
           """Check status of SMS messages"""
           url = f"{SMS_BASE_URL}/sms/status"
           queries = []
           
           if message_ids:
               queries.extend([{"message_id": mid, "track_id": ""} for mid in message_ids])
           if track_ids:
               queries.extend([{"message_id": "", "track_id": tid} for tid in track_ids])
           
           response = requests.post(url, json=queries, headers=self.headers)
           return response.json()
       
       def get_bulk_status(self, group_id):
           """Get status of bulk SMS by group ID"""
           url = f"{SMS_BASE_URL}/sms/status/bulk"
           params = {"models.QueryParamGroupID": group_id}
           response = requests.get(url, params=params, headers=self.headers)
           return response.json()
       
       def get_received_messages(self, receivers, epoch_time):
           """Get received messages for specific receivers"""
           url = f"{SMS_BASE_URL}/sms/received-message"
           payload = {
               "receiver": receivers,
               "epoch_time": epoch_time
           }
           response = requests.post(url, json=payload, headers=self.headers)
           return response.json()
   
   # Example usage
   client = SMSClient("john_doe")
   
   # Send simple SMS
   result = client.send_simple_sms(
       sender="1234567890",
       receiver="0987654321",
       message="Hello from Python!"
   )
   print(f"SMS sent: {result}")
   
   # Send bulk SMS
   bulk_result = client.send_bulk_sms(
       senders=["1234567890", "1234567890"],
       receivers=["0987654321", "5551234567"],
       messages=["Hello User 1!", "Hello User 2!"],
       track_ids=["trk_001", "trk_002"]
   )
   print(f"Bulk SMS sent: {bulk_result}")
   
   # Check status
   group_id = bulk_result.get("group_id")
   if group_id:
       status = client.get_bulk_status(group_id)
       print(f"Status: {status}")
   
   # Get received messages
   import time
   epoch_time = str(int(time.time()) - 86400)  # Last 24 hours
   received = client.get_received_messages(
       receivers=["1234567890"],
       epoch_time=epoch_time
   )
   print(f"Received messages: {received}")
   ```

---

   ### Prefix Management with Python

   ```python
   import requests
   
   SMS_BASE_URL = "http://10.10.10.21:8080/api/v1"
   
   class PrefixManager:
       def __init__(self, username):
           self.username = username
           self.headers = {"user-name": username, "Content-Type": "application/json"}
       
       def add_prefixes(self, numbers):
           """Add prefixes for user"""
           url = f"{SMS_BASE_URL}/prefix/add"
           payload = {
               "user": self.username,
               "numbers": numbers
           }
           response = requests.post(url, json=payload, headers=self.headers)
           return response.json()
       
       def remove_prefixes(self, prefixes):
           """Remove prefixes for user"""
           url = f"{SMS_BASE_URL}/prefix/remove"
           payload = {
               "user": self.username,
               "numbers": prefixes
           }
           response = requests.delete(url, json=payload, headers=self.headers)
           return response.json()
       
       def search_prefixes(self, username=None, prefix=None, page=1, page_size=50):
           """Search prefixes with filters"""
           url = f"{SMS_BASE_URL}/prefix/search"
           params = {"page": page, "page_size": page_size}
           
           if username:
               params["username"] = username
           if prefix:
               params["prefix"] = prefix
           
           response = requests.get(url, params=params, headers=self.headers)
           return response.json()
   
   # Example usage
   manager = PrefixManager("john_doe")
   
   # Add prefixes
   result = manager.add_prefixes([
       {"number": "1234", "is_prefix": True},
       {"number": "5678", "is_prefix": False}
   ])
   print(f"Add result: {result}")
   
   # Search prefixes
   prefixes = manager.search_prefixes(username="john_doe")
   print(f"Prefixes: {prefixes}")
   
   # Remove prefixes
   remove_result = manager.remove_prefixes(["1234", "5678"])
   print(f"Remove result: {remove_result}")
   ```

---

   ### MSISDN Mapping with Python

   ```python
   import requests
   
   BASE_URL = "http://localhost:8080/v1/api"
   
   class MSISDNManager:
       def __init__(self):
           self.base_url = f"{BASE_URL}/msisdn-mappings"
       
       def list_mappings(self, begin=0, limit=100):
           """List all MSISDN mappings"""
           params = {"begin": begin, "limit": limit}
           response = requests.get(self.base_url, params=params)
           return response.json()
       
       def get_mapping(self, msisdn):
           """Get specific MSISDN mapping"""
           url = f"{self.base_url}/{msisdn}"
           response = requests.get(url)
           return response.json()
       
       def add_mapping(self, msisdn, new_msisdn, created_at=None):
           """Add new MSISDN mapping"""
           from datetime import datetime
           payload = {
               "msisdn": msisdn,
               "new_msisdn": new_msisdn,
               "created_at": created_at or datetime.now().isoformat()
           }
           response = requests.post(self.base_url, json=payload)
           return response.json()
       
       def remove_mapping(self, msisdn):
           """Remove MSISDN mapping"""
           url = f"{self.base_url}/{msisdn}"
           response = requests.delete(url)
           return response.json()
   
   # Example usage
   msisdn_manager = MSISDNManager()
   
   # Add mapping
   result = msisdn_manager.add_mapping("1234567890", "0987654321")
   print(f"Mapping added: {result}")
   
   # Get mapping
   mapping = msisdn_manager.get_mapping("1234567890")
   print(f"Mapping: {mapping}")
   
   # List all mappings
   all_mappings = msisdn_manager.list_mappings(begin=0, limit=50)
   print(f"All mappings: {all_mappings}")
   
   # Remove mapping
   remove_result = msisdn_manager.remove_mapping("1234567890")
   print(f"Remove result: {remove_result}")
   ```

---

   ### Complete Python Example with Error Handling

   ```python
   import requests
   import logging
   from typing import Optional, Dict, Any
   
   # Configure logging
   logging.basicConfig(level=logging.INFO)
   logger = logging.getLogger(__name__)
   
   class APIClient:
       def __init__(self, base_url: str, username: str, password: str):
           self.base_url = base_url
           self.username = username
           self.password = password
           self.token: Optional[str] = None
           self.session = requests.Session()
       
       def login(self) -> bool:
           """Login and store token"""
           try:
               url = f"{self.base_url}/login"
               response = self.session.post(url, json={
                   "username": self.username,
                   "password": self.password
               })
               response.raise_for_status()
               self.token = response.json().get("token")
               logger.info("Login successful")
               return True
           except requests.exceptions.RequestException as e:
               logger.error(f"Login failed: {e}")
               return False
       
       def _get_headers(self) -> Dict[str, str]:
           """Get headers with authentication"""
           headers = {"Content-Type": "application/json"}
           if self.token:
               headers["Authorization"] = f"Bearer {self.token}"
           return headers
       
       def make_request(self, method: str, endpoint: str, **kwargs) -> Optional[Dict[str, Any]]:
           """Make authenticated request with error handling"""
           url = f"{self.base_url}{endpoint}"
           headers = self._get_headers()
           
           try:
               response = self.session.request(
                   method, url, headers=headers, **kwargs
               )
               response.raise_for_status()
               return response.json()
           except requests.exceptions.HTTPError as e:
               if e.response.status_code == 401:
                   logger.warning("Token expired, re-authenticating...")
                   if self.login():
                       return self.make_request(method, endpoint, **kwargs)
               logger.error(f"HTTP error: {e}")
               return None
           except requests.exceptions.RequestException as e:
               logger.error(f"Request failed: {e}")
               return None
       
       def get_dashboard(self):
           """Get admin dashboard"""
           return self.make_request("GET", "/admin/dashboard")
       
       def create_user(self, username: str, password: str, role: str = "user"):
           """Create new user"""
           return self.make_request("POST", "/users", json={
               "username": username,
               "password": password,
               "role": role
           })
   
   # Example usage
   client = APIClient(
       base_url="http://localhost:8080/v1/api",
       username="admin",
       password="admin_password"
   )
   
   if client.login():
       # Create user
       new_user = client.create_user("test_user", "password123")
       print(f"User created: {new_user}")
       
       # Get dashboard
       dashboard = client.get_dashboard()
       print(f"Dashboard: {dashboard}")
   ```

---

   ## JavaScript Examples

   ### Authentication and User Management

   ```javascript
   const BASE_URL = "http://localhost:8080/v1/api";
   
   // Register a new user
   async function registerUser(username, password, role = "user") {
     const response = await fetch(`${BASE_URL}/users`, {
       method: "POST",
       headers: { "Content-Type": "application/json" },
       body: JSON.stringify({ username, password, role })
     });
     return await response.json();
   }
   
   // Login and get JWT token
   async function login(username, password) {
     const response = await fetch(`${BASE_URL}/login`, {
       method: "POST",
       headers: { "Content-Type": "application/json" },
       body: JSON.stringify({ username, password })
     });
     
     if (!response.ok) {
       throw new Error(`Login failed: ${response.statusText}`);
     }
     
     const data = await response.json();
     return data.token;
   }
   
   // Secure login with CAPTCHA
   async function secureLogin(username, password, captchaId, captchaSolution) {
     const response = await fetch(`${BASE_URL}/secure-login`, {
       method: "POST",
       headers: { "Content-Type": "application/json" },
       body: JSON.stringify({
         username,
         password,
         captchaId,
         captchaSolution
       })
     });
     return await response.json();
   }
   
   // Get CAPTCHA
   async function getCaptcha() {
     const response = await fetch(`${BASE_URL}/captcha/new`);
     const data = await response.json();
     return data.captchaId;
   }
   
   // Authenticated request helper
   async function authenticatedRequest(url, token, options = {}) {
     const response = await fetch(url, {
       ...options,
       headers: {
         ...options.headers,
         "Authorization": `Bearer ${token}`,
         "Content-Type": "application/json"
       }
     });
     
     if (response.status === 401) {
       throw new Error("Token expired. Please login again.");
     }
     
     return await response.json();
   }
   
   // Example usage
   (async () => {
     try {
       // Register
       const registerResult = await registerUser("john_doe", "SecurePass123!");
       console.log("Registration:", registerResult);
       
       // Login
       const token = await login("john_doe", "SecurePass123!");
       console.log("Token:", token);
       
       // Get admin dashboard
       const dashboard = await authenticatedRequest(
         `${BASE_URL}/admin/dashboard`,
         token
       );
       console.log("Dashboard:", dashboard);
       
       // Update user role
       const updateResult = await authenticatedRequest(
         `${BASE_URL}/users/john_doe/role`,
         token,
         {
           method: "PUT",
           body: JSON.stringify({ role: "admin" })
         }
       );
       console.log("Role updated:", updateResult);
       
     } catch (error) {
       console.error("Error:", error.message);
     }
   })();
   ```

---

   ### SMS Sending with JavaScript

   ```javascript
   const SMS_BASE_URL = "http://10.10.10.21:8080/api/v1";
   
   class SMSClient {
     constructor(username) {
       this.username = username;
       this.headers = {
         "user-name": username,
         "Content-Type": "application/json"
       };
     }
     
     async sendSimpleSMS(sender, receiver, message) {
       const response = await fetch(`${SMS_BASE_URL}/sms/send/simple`, {
         method: "POST",
         headers: this.headers,
         body: JSON.stringify({ sender, receiver, message })
       });
       return await response.json();
     }
     
     async sendBulkSMS(senders, receivers, messages, trackIds = null) {
       const track_id = trackIds || receivers.map((_, i) => `trk_${i}`);
       
       const response = await fetch(`${SMS_BASE_URL}/sms/send/bulk`, {
         method: "POST",
         headers: this.headers,
         body: JSON.stringify({
           sender: senders,
           receiver: receivers,
           message: messages,
           track_id
         })
       });
       return await response.json();
     }
     
     async checkSMSStatus(messageIds = [], trackIds = []) {
       const queries = [
         ...messageIds.map(id => ({ message_id: id, track_id: "" })),
         ...trackIds.map(id => ({ message_id: "", track_id: id }))
       ];
       
       const response = await fetch(`${SMS_BASE_URL}/sms/status`, {
         method: "POST",
         headers: this.headers,
         body: JSON.stringify(queries)
       });
       return await response.json();
     }
     
     async getBulkStatus(groupId) {
       const response = await fetch(
         `${SMS_BASE_URL}/sms/status/bulk?models.QueryParamGroupID=${groupId}`,
         { headers: this.headers }
       );
       return await response.json();
     }
     
     async getReceivedMessages(receivers, epochTime) {
       const response = await fetch(`${SMS_BASE_URL}/sms/received-message`, {
         method: "POST",
         headers: this.headers,
         body: JSON.stringify({ receiver: receivers, epoch_time: epochTime })
       });
       return await response.json();
     }
   }
   
   // Example usage with async/await
   (async () => {
     const client = new SMSClient("john_doe");
     
     try {
       // Send simple SMS
       const result = await client.sendSimpleSMS(
         "1234567890",
         "0987654321",
         "Hello from JavaScript!"
       );
       console.log("SMS sent:", result);
       
       // Send bulk SMS
       const bulkResult = await client.sendBulkSMS(
         ["1234567890", "1234567890"],
         ["0987654321", "5551234567"],
         ["Hello User 1!", "Hello User 2!"],
         ["trk_001", "trk_002"]
       );
       console.log("Bulk SMS sent:", bulkResult);
       
       // Check status
       const groupId = bulkResult.group_id;
       if (groupId) {
         const status = await client.getBulkStatus(groupId);
         console.log("Status:", status);
       }
       
       // Check received messages
       const epochTime = String(Math.floor(Date.now() / 1000) - 86400);
       const received = await client.getReceivedMessages(
         ["1234567890"],
         epochTime
       );
       console.log("Received messages:", received);
       
     } catch (error) {
       console.error("Error:", error);
     }
   })();
   ```

---

   ### React Integration Example

   ```jsx
   import React, { useState, useEffect } from 'react';
   
   function SMSComponent() {
     const [username] = useState("john_doe");
     const [smsData, setSmsData] = useState({
       sender: "",
       receiver: "",
       message: ""
     });
     const [status, setStatus] = useState(null);
     const [loading, setLoading] = useState(false);
     const [sentMessages, setSentMessages] = useState([]);
   
     const sendSMS = async () => {
       setLoading(true);
       setStatus(null);
       
       try {
         const response = await fetch("http://10.10.10.21:8080/api/v1/sms/send/simple", {
           method: "POST",
           headers: {
             "user-name": username,
             "Content-Type": "application/json"
           },
           body: JSON.stringify(smsData)
         });
         
         if (!response.ok) {
           throw new Error(`HTTP error! status: ${response.status}`);
         }
         
         const result = await response.json();
         setStatus({ success: true, data: result });
         setSentMessages(prev => [...prev, {
           ...smsData,
           groupId: result.group_id,
           timestamp: new Date().toISOString()
         }]);
         
         // Clear form
         setSmsData({ sender: "", receiver: "", message: "" });
         
       } catch (error) {
         setStatus({ success: false, error: error.message });
       } finally {
         setLoading(false);
       }
     };
   
     const checkStatus = async (groupId) => {
       try {
         const response = await fetch(
           `http://10.10.10.21:8080/api/v1/sms/status/bulk?models.QueryParamGroupID=${groupId}`,
           {
             headers: { "user-name": username }
           }
         );
         const result = await response.json();
         console.log("Status:", result);
         alert(`Status: ${JSON.stringify(result, null, 2)}`);
       } catch (error) {
         console.error("Error checking status:", error);
       }
     };
   
     return (
       <div style={{ maxWidth: "600px", margin: "20px auto", padding: "20px" }}>
         <h2>Send SMS</h2>
         
         <div style={{ marginBottom: "10px" }}>
           <input
             type="text"
             placeholder="Sender"
             value={smsData.sender}
             onChange={(e) => setSmsData({...smsData, sender: e.target.value})}
             style={{ width: "100%", padding: "8px", marginBottom: "10px" }}
           />
           <input
             type="text"
             placeholder="Receiver"
             value={smsData.receiver}
             onChange={(e) => setSmsData({...smsData, receiver: e.target.value})}
             style={{ width: "100%", padding: "8px", marginBottom: "10px" }}
           />
           <textarea
             placeholder="Message"
             value={smsData.message}
             onChange={(e) => setSmsData({...smsData, message: e.target.value})}
             style={{ width: "100%", padding: "8px", marginBottom: "10px", minHeight: "80px" }}
           />
           <button 
             onClick={sendSMS} 
             disabled={loading || !smsData.sender || !smsData.receiver || !smsData.message}
             style={{ 
               width: "100%", 
               padding: "10px", 
               backgroundColor: "#007bff", 
               color: "white", 
               border: "none", 
               cursor: "pointer" 
             }}
           >
             {loading ? "Sending..." : "Send SMS"}
           </button>
         </div>
         
         {status && (
           <div style={{ 
             padding: "10px", 
             marginTop: "10px", 
             backgroundColor: status.success ? "#d4edda" : "#f8d7da",
             border: `1px solid ${status.success ? "#c3e6cb" : "#f5c6cb"}`,
             borderRadius: "4px"
           }}>
             {status.success ? (
               <div>
                 <p style={{ margin: "0 0 10px 0" }}>
                   ✓ SMS sent successfully! Group ID: {status.data.group_id}
                 </p>
                 {status.data.status && status.data.status.length > 0 && (
                   <p style={{ margin: 0 }}>
                     Message ID: {status.data.status[0].message_id}
                   </p>
                 )}
               </div>
             ) : (
               <p style={{ margin: 0 }}>✗ Error: {status.error}</p>
             )}
           </div>
         )}
         
         {sentMessages.length > 0 && (
           <div style={{ marginTop: "20px" }}>
             <h3>Sent Messages</h3>
             {sentMessages.map((msg, index) => (
               <div key={index} style={{ 
                 padding: "10px", 
                 marginBottom: "10px", 
                 backgroundColor: "#f8f9fa",
                 border: "1px solid #dee2e6",
                 borderRadius: "4px"
               }}>
                 <p><strong>From:</strong> {msg.sender} → {msg.receiver}</p>
                 <p><strong>Message:</strong> {msg.message}</p>
                 <p><strong>Group ID:</strong> {msg.groupId}</p>
                 <button 
                   onClick={() => checkStatus(msg.groupId)}
                   style={{ 
                     padding: "5px 10px", 
                     backgroundColor: "#28a745", 
                     color: "white", 
                     border: "none", 
                     borderRadius: "3px",
                     cursor: "pointer"
                   }}
                 >
                   Check Status
                 </button>
               </div>
             ))}
           </div>
         )}
       </div>
     );
   }
   
   export default SMSComponent;
   ```

---

   ### Node.js Express Integration

   ```javascript
   const express = require('express');
   const axios = require('axios');
   const app = express();
   
   app.use(express.json());
   
   const AUTH_BASE_URL = "http://localhost:8080/v1/api";
   const SMS_BASE_URL = "http://10.10.10.21:8080/api/v1";
   
   // Middleware to verify JWT
   async function verifyToken(req, res, next) {
     const token = req.headers.authorization?.split(' ')[1];
     
     if (!token) {
       return res.status(401).json({ error: "No token provided" });
     }
     
     try {
       // Verify token by making a request to protected endpoint
       await axios.get(`${AUTH_BASE_URL}/admin/dashboard`, {
         headers: { Authorization: `Bearer ${token}` }
       });
       req.token = token;
       next();
     } catch (error) {
       res.status(401).json({ error: "Invalid token" });
     }
   }
   
   // Login endpoint
   app.post('/api/login', async (req, res) => {
     try {
       const { username, password } = req.body;
       const response = await axios.post(`${AUTH_BASE_URL}/login`, {
         username,
         password
       });
       res.json(response.data);
     } catch (error) {
       res.status(error.response?.status || 500).json({
         error: error.response?.data || "Login failed"
       });
     }
   });
   
   // Send SMS endpoint
   app.post('/api/sms/send', verifyToken, async (req, res) => {
     try {
       const { sender, receiver, message } = req.body;
       const response = await axios.post(
         `${SMS_BASE_URL}/sms/send/simple`,
         { sender, receiver, message },
         { headers: { "user-name": req.body.username } }
       );
       res.json(response.data);
     } catch (error) {
       res.status(error.response?.status || 500).json({
         error: error.response?.data || "SMS send failed"
       });
     }
   });
   
   // Check SMS status
   app.get('/api/sms/status/:groupId', verifyToken, async (req, res) => {
     try {
       const response = await axios.get(
         `${SMS_BASE_URL}/sms/status/bulk`,
         {
           params: { "models.QueryParamGroupID": req.params.groupId },
           headers: { "user-name": req.query.username }
         }
       );
       res.json(response.data);
     } catch (error) {
       res.status(error.response?.status || 500).json({
         error: error.response?.data || "Status check failed"
       });
     }
   });
   
   app.listen(3000, () => {
     console.log('Server running on port 3000');
   });
   ```

---

   ## Sequence Diagrams

   ### User Registration and Authentication Flow

   ```
   ┌──────┐          ┌────────┐          ┌──────────┐
   │Client│          │API     │          │Database  │
   └──┬───┘          └───┬────┘          └────┬─────┘
      │                  │                     │
      │ POST /users      │                     │
      │─────────────────>│                     │
      │  {username,pwd}  │                     │
      │                  │  Validate input     │
      │                  │  Hash password      │
      │                  │  Create user        │
      │                  │────────────────────>│
      │                  │                     │
      │                  │  User created       │
      │                  │<────────────────────│
      │  200 OK          │                     │
      │<─────────────────│                     │
      │                  │                     │
      │ POST /login      │                     │
      │─────────────────>│                     │
      │  {username,pwd}  │                     │
      │                  │  Verify credentials │
      │                  │────────────────────>│
      │                  │                     │
      │                  │  User found         │
      │                  │<────────────────────│
      │                  │  Generate JWT       │
      │  {token}         │                     │
      │<─────────────────│                     │
      │                  │                     │
      │ GET /dashboard   │                     │
      │─────────────────>│                     │
      │  Auth: Bearer    │                     │
      │                  │  Verify JWT         │
      │                  │  Check permissions  │
      │                  │  Query data         │
      │                  │────────────────────>│
      │                  │                     │
      │                  │  Dashboard data     │
      │                  │<────────────────────│
      │  200 OK          │                     │
      │<─────────────────│                     │
      │                  │                     │
   ```

---

   ### Secure Login with CAPTCHA Flow

   ```
   ┌──────┐     ┌────────┐     ┌─────────┐     ┌──────────┐
   │Client│     │API     │     │CAPTCHA  │     │Database  │
   └──┬───┘     └───┬────┘     └────┬────┘     └────┬─────┘
      │             │               │               │
      │GET /captcha/new             │               │
      │───────────>│                │               │
      │            │  Generate      │               │
      │            │──────────────>│               │
      │            │                │               │
      │            │  {captchaId}   │               │
      │            │<──────────────│               │
      │{captchaId} │                │               │
      │<───────────│                │               │
      │            │                │               │
      │GET /captcha/image/{id}      │               │
      │───────────>│                │               │
      │            │  Get image     │               │
      │            │──────────────>│               │
      │            │                │               │
      │            │  PNG image     │               │
      │            │<──────────────│               │
      │ PNG image  │                │               │
      │<───────────│                │               │
      │            │                │               │
      │[User solves CAPTCHA]        │               │
      │            │                │               │
      │POST /secure-login           │               │
      │───────────>│                │               │
      │{user,pwd,  │                │               │
      │captchaId,  │  Verify        │               │
      │solution}   │  solution      │               │
      │            │──────────────>│               │
      │            │                │               │
      │            │  Valid         │               │
      │            │<──────────────│               │
      │            │  Verify credentials            │
      │            │──────────────────────────────>│
      │            │                │               │
      │            │  User found    │               │
      │            │<──────────────────────────────│
      │            │  Generate JWT  │               │
      │  {token}   │                │               │
      │<───────────│                │               │
   ```

---

   ### SMS Sending and Status Tracking Flow

   ```
   ┌──────┐     ┌────────┐     ┌─────────┐     ┌──────────┐
   │Client│     │SMS API │     │SMS      │     │Database  │
   │      │     │        │     │Gateway  │     │          │
   └──┬───┘     └───┬────┘     └────┬────┘     └────┬─────┘
      │             │               │               │
      │POST /sms/send/simple        │               │
      │───────────>│                │               │
      │{sender,    │  Validate      │               │
      │receiver,   │  Store SMS     │               │
      │message}    │──────────────────────────────>│
      │            │                │               │
      │            │  SMS record    │               │
      │            │<──────────────────────────────│
      │            │  Send to gateway               │
      │            │──────────────>│               │
      │            │                │               │
      │            │  Queued        │               │
      │            │<──────────────│               │
      │{group_id,  │                │               │
      │message_id, │                │               │
      │status}     │                │               │
      │<───────────│                │               │
      │            │                │               │
      │[Wait for delivery]          │               │
      │            │                │               │
      │            │  Delivery report               │
      │            │<──────────────│               │
      │            │  Update status │               │
      │            │──────────────────────────────>│
      │            │                │               │
      │POST /sms/status             │               │
      │───────────>│                │               │
      │{message_id}│  Query status  │               │
      │            │──────────────────────────────>│
      │            │                │               │
      │            │  Status info   │               │
      │            │<──────────────────────────────│
      │{status:    │                │               │
      │DELIVRD}    │                │               │
      │<───────────│                │               │
   ```

---

   ### Bulk SMS Sending Flow

   ```
   ┌──────┐     ┌────────┐     ┌─────────┐     ┌──────────┐
   │Client│     │SMS API │     │SMS      │     │Database  │
   │      │     │        │     │Gateway  │     │          │
   └──┬───┘     └───┬────┘     └────┬────┘     └────┬─────┘
      │             │               │               │
      │POST /sms/send/bulk          │               │
      │───────────>│                │               │
      │{senders[], │  Validate all  │               │
      │receivers[],│  Create group  │               │
      │messages[], │  Store batch   │               │
      │track_ids[]}│──────────────────────────────>│
      │            │                │               │
      │            │  Group created │               │
      │            │<──────────────────────────────│
      │            │                │               │
      │            │  ┌─────────────────────┐      │
      │            │  │For each SMS message │      │
      │            │  └──────────┬──────────┘      │
      │            │             │                  │
      │            │  Send SMS 1 │                  │
      │            │──────────────>                 │
      │            │  Send SMS 2 │                  │
      │            │──────────────>                 │
      │            │  Send SMS N │                  │
      │            │──────────────>                 │
      │            │             │                  │
      │            │  Batch queued│                 │
      │            │<──────────────                 │
      │{group_id,  │             │                  │
      │status[]}   │             │                  │
      │<───────────│             │                  │
      │            │             │                  │
      │GET /sms/status/bulk     │                  │
      │───────────>│             │                  │
      │{group_id}  │  Query all  │                  │
      │            │  in group   │                  │
      │            │──────────────────────────────>│
      │            │             │                  │
      │            │  All statuses                  │
      │            │<──────────────────────────────│
      │{status[]}  │             │                  │
      │<───────────│             │                  │
   ```

---

   ### MSISDN Mapping Flow

   ```
   ┌──────┐          ┌────────┐          ┌──────────┐
   │Client│          │API     │          │Database  │
   └──┬───┘          └───┬────┘          └────┬─────┘
      │                  │                     │
      │POST /msisdn-mappings                   │
      │─────────────────>│                     │
      │{msisdn,          │  Validate           │
      │new_msisdn}       │  Check duplicate    │
      │                  │────────────────────>│
      │                  │                     │
      │                  │  No conflict        │
      │                  │<────────────────────│
      │                  │  Store mapping      │
      │                  │────────────────────>│
      │                  │                     │
      │                  │  Mapping created    │
      │                  │<────────────────────│
      │  201 Created     │                     │
      │<─────────────────│                     │
      │                  │                     │
      │GET /msisdn-mappings/{msisdn}           │
      │─────────────────>│                     │
      │                  │  Query mapping      │
      │                  │────────────────────>│
      │                  │                     │
      │                  │  Mapping data       │
      │                  │<────────────────────│
      │  200 OK          │                     │
      │<─────────────────│                     │
      │                  │                     │
      │DELETE /msisdn-mappings/{msisdn}        │
      │─────────────────>│                     │
      │                  │  Remove mapping     │
      │                  │────────────────────>│
      │                  │                     │
      │                  │  Deleted            │
      │                  │<────────────────────│
      │  200 OK          │                     │
      │<─────────────────│                     │
   ```

---

   ## Try It Out

   ### 🛠️ Interactive Testing Tools

   1. **Swagger UI**
      ```
      http://localhost:8080/swagger/
      ```
      Provides interactive API documentation where you can test endpoints directly in your browser.
      
   2. **Postman / Insomnia**
      Download the Postman collection (see below) or use Insomnia for a rich API testing experience with:
      - Request history
      - Environment variables
      - Pre-request scripts
      - Test assertions
   
   3. **cURL**
      All examples in this documentation use cURL, which is available on most systems. Great for:
      - Quick testing
      - Scripting and automation
      - CI/CD integration

---

   ### 🔧 Testing Checklist

   Before deploying to production, test these critical flows:

   - [ ] User registration and login
   - [ ] JWT token authentication and expiration
   - [ ] Secure login with CAPTCHA
   - [ ] Simple SMS sending
   - [ ] Bulk SMS sending and status tracking
   - [ ] Prefix management (add/remove/search)
   - [ ] MSISDN mapping operations
   - [ ] Error handling for invalid inputs
   - [ ] Rate limit behavior
   - [ ] Admin dashboard access control

---

   ### 🐛 Debugging Tips

   1. **Enable detailed logging**
      ```bash
      export DEBUG=true
      export LOG_LEVEL=debug
      ```
   
   2. **Check response headers**
      ```bash
      curl -v http://localhost:8080/v1/api/users
      ```
   
   3. **Validate JSON payloads**
      ```bash
      echo '{"username":"test"}' | jq .
      ```
   
   4. **Test authentication separately**
      ```bash
      # Get token
      TOKEN=$(curl -s -X POST http://localhost:8080/v1/api/login \
        -H "Content-Type: application/json" \
        -d '{"username":"test","password":"pass"}' | jq -r '.token')
      
      # Use token
      curl -H "Authorization: Bearer $TOKEN" \
        http://localhost:8080/v1/api/admin/dashboard
      ```
   
   5. **Test rate limiting**
      ```bash
      # Send multiple requests quickly
      for i in {1..10}; do
        curl -X POST http://10.10.10.21:8080/api/v1/sms/send/simple \
          -H "user-name: test" \
          -H "Content-Type: application/json" \
          -d '{"sender":"123","receiver":"456","message":"test"}' &
      done
      wait
      ```
   
   6. **Monitor API responses**
      ```bash
      # Watch logs in real-time
      tail -f /var/log/api/access.log
      
      # Filter for errors
      tail -f /var/log/api/error.log | grep "ERROR"
      ```

---

   ## Postman Collection

   ### 📦 Import Instructions

   You can import this API into Postman for easy testing:

   **Option 1: Import from URL**
   ```
   https://your-domain.com/api/postman-collection.json
   ```

   **Option 2: Import manually**
   1. Open Postman
   2. Click "Import" button
   3. Paste the collection JSON below
   4. Configure environment variables

---

   ### 🔧 Environment Variables

   Set these variables in Postman:

   | Variable | Description | Example |
   |----------|-------------|---------|
   | `auth_base_url` | Auth service base URL | `http://localhost:8080/v1/api` |
   | `sms_base_url` | SMS service base URL | `http://10.10.10.21:8080/api/v1` |
   | `jwt_token` | JWT authentication token | Auto-set after login |
   | `username` | Your username | `john_doe` |
   | `password` | Your password | `SecurePass123!` |
   | `test_phone` | Test phone number | `1234567890` |

---

   ### 📝 Postman Collection Structure

   ```
   Auth Service API
   ├── Authentication
   │   ├── Register User
   │   ├── Login
   │   ├── Secure Login with CAPTCHA
   │   └── Get CAPTCHA
   ├── User Management
   │   ├── Update User Role
   │   └── Delete User
   ├── Roles
   │   ├── List Roles
   │   └── Create Role
   ├── Admin
   │   ├── Dashboard
   │   └── Create Custom Endpoint
   ├── MSISDN Mapping
   │   ├── List Mappings
   │   ├── Get Mapping
   │   ├── Add Mapping
   │   └── Remove Mapping
   ├── SMS
   │   ├── Send Simple SMS
   │   ├── Send Bulk SMS
   │   ├── Get SMS Status
   │   ├── Get Bulk Status
   │   └── Get Received Messages
   └── Prefix Management
       ├── Add Prefixes
       ├── Remove Prefixes
       └── Search Prefixes
   ```

---

   ### 🔄 Pre-request Scripts

   Add this to Postman collection pre-request script to auto-refresh tokens:

   ```javascript
   // Check if token exists and is not expired
   const token = pm.environment.get("jwt_token");
   const tokenExpiry = pm.environment.get("token_expiry");
   
   if (!token || (tokenExpiry && Date.now() > tokenExpiry)) {
       // Token expired or missing, refresh it
       pm.sendRequest({
           url: pm.environment.get("auth_base_url") + "/login",
           method: 'POST',
           header: {
               'Content-Type': 'application/json'
           },
           body: {
               mode: 'raw',
               raw: JSON.stringify({
                   username: pm.environment.get("username"),
                   password: pm.environment.get("password")
               })
           }
       }, function (err, response) {
           if (!err && response.code === 200) {
               const newToken = response.json().token;
               pm.environment.set("jwt_token", newToken);
               // Set expiry to 23 hours from now
               pm.environment.set("token_expiry", Date.now() + (23 * 60 * 60 * 1000));
           }
       });
   }
   ```

---

   ### 📥 Complete Postman Collection JSON

   ```json
   {
     "info": {
       "name": "Auth & SMS Service API",
       "description": "Complete API collection for Auth, MSISDN, and SMS services",
       "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json",
       "version": "1.0.0"
     },
     "variable": [
       {
         "key": "auth_base_url",
         "value": "http://localhost:8080/v1/api",
         "type": "string"
       },
       {
         "key": "sms_base_url",
         "value": "http://10.10.10.21:8080/api/v1",
         "type": "string"
       }
     ],
     "item": [
       {
         "name": "Authentication",
         "item": [
           {
             "name": "Register User",
             "request": {
               "method": "POST",
               "header": [
                 {
                   "key": "Content-Type",
                   "value": "application/json"
                 }
               ],
               "body": {
                 "mode": "raw",
                 "raw": "{\n  \"username\": \"{{username}}\",\n  \"password\": \"{{password}}\",\n  \"role\": \"user\"\n}"
               },
               "url": {
                 "raw": "{{auth_base_url}}/users",
                 "host": ["{{auth_base_url}}"],
                 "path": ["users"]
               }
             }
           },
           {
             "name": "Login",
             "event": [
               {
                 "listen": "test",
                 "script": {
                   "exec": [
                     "if (pm.response.code === 200) {",
                     "    const response = pm.response.json();",
                     "    pm.environment.set('jwt_token', response.token);",
                     "    pm.environment.set('token_expiry', Date.now() + (23 * 60 * 60 * 1000));",
                     "}"
                   ],
                   "type": "text/javascript"
                 }
               }
             ],
             "request": {
               "method": "POST",
               "header": [
                 {
                   "key": "Content-Type",
                   "value": "application/json"
                 }
               ],
               "body": {
                 "mode": "raw",
                 "raw": "{\n  \"username\": \"{{username}}\",\n  \"password\": \"{{password}}\"\n}"
               },
               "url": {
                 "raw": "{{auth_base_url}}/login",
                 "host": ["{{auth_base_url}}"],
                 "path": ["login"]
               }
             }
           },
           {
             "name": "Get CAPTCHA",
             "event": [
               {
                 "listen": "test",
                 "script": {
                   "exec": [
                     "if (pm.response.code === 200) {",
                     "    const response = pm.response.json();",
                     "    pm.environment.set('captcha_id', response.captchaId);",
                     "}"
                   ],
                   "type": "text/javascript"
                 }
               }
             ],
             "request": {
               "method": "GET",
               "header": [],
               "url": {
                 "raw": "{{auth_base_url}}/captcha/new",
                 "host": ["{{auth_base_url}}"],
                 "path": ["captcha", "new"]
               }
             }
           },
           {
             "name": "Secure Login with CAPTCHA",
             "request": {
               "method": "POST",
               "header": [
                 {
                   "key": "Content-Type",
                   "value": "application/json"
                 }
               ],
               "body": {
                 "mode": "raw",
                 "raw": "{\n  \"username\": \"{{username}}\",\n  \"password\": \"{{password}}\",\n  \"captchaId\": \"{{captcha_id}}\",\n  \"captchaSolution\": \"SOLUTION_HERE\"\n}"
               },
               "url": {
                 "raw": "{{auth_base_url}}/secure-login",
                 "host": ["{{auth_base_url}}"],
                 "path": ["secure-login"]
               }
             }
           }
         ]
       },
       {
         "name": "SMS",
         "item": [
           {
             "name": "Send Simple SMS",
             "event": [
               {
                 "listen": "test",
                 "script": {
                   "exec": [
                     "if (pm.response.code === 200) {",
                     "    const response = pm.response.json();",
                     "    pm.environment.set('last_group_id', response.group_id);",
                     "    if (response.status && response.status.length > 0) {",
                     "        pm.environment.set('last_message_id', response.status[0].message_id);",
                     "        pm.environment.set('last_track_id', response.status[0].track_id);",
                     "    }",
                     "}"
                   ],
                   "type": "text/javascript"
                 }
               }
             ],
             "request": {
               "method": "POST",
               "header": [
                 {
                   "key": "Content-Type",
                   "value": "application/json"
                 },
                 {
                   "key": "user-name",
                   "value": "{{username}}"
                 }
               ],
               "body": {
                 "mode": "raw",
                 "raw": "{\n  \"sender\": \"{{test_phone}}\",\n  \"receiver\": \"0987654321\",\n  \"message\": \"Test message from Postman\"\n}"
               },
               "url": {
                 "raw": "{{sms_base_url}}/sms/send/simple",
                 "host": ["{{sms_base_url}}"],
                 "path": ["sms", "send", "simple"]
               }
             }
           },
           {
             "name": "Send Bulk SMS",
             "request": {
               "method": "POST",
               "header": [
                 {
                   "key": "Content-Type",
                   "value": "application/json"
                 },
                 {
                   "key": "user-name",
                   "value": "{{username}}"
                 }
               ],
               "body": {
                 "mode": "raw",
                 "raw": "{\n  \"sender\": [\"{{test_phone}}\", \"{{test_phone}}\"],\n  \"receiver\": [\"0987654321\", \"5551234567\"],\n  \"message\": [\"Hello User 1!\", \"Hello User 2!\"],\n  \"track_id\": [\"trk_001\", \"trk_002\"]\n}"
               },
               "url": {
                 "raw": "{{sms_base_url}}/sms/send/bulk",
                 "host": ["{{sms_base_url}}"],
                 "path": ["sms", "send", "bulk"]
               }
             }
           },
           {
             "name": "Get SMS Status",
             "request": {
               "method": "POST",
               "header": [
                 {
                   "key": "Content-Type",
                   "value": "application/json"
                 },
                 {
                   "key": "user-name",
                   "value": "{{username}}"
                 }
               ],
               "body": {
                 "mode": "raw",
                 "raw": "[\n  {\n    \"message_id\": \"{{last_message_id}}\",\n    \"track_id\": \"\"\n  }\n]"
               },
               "url": {
                 "raw": "{{sms_base_url}}/sms/status",
                 "host": ["{{sms_base_url}}"],
                 "path": ["sms", "status"]
               }
             }
           },
           {
             "name": "Get Bulk SMS Status",
             "request": {
               "method": "GET",
               "header": [
                 {
                   "key": "user-name",
                   "value": "{{username}}"
                 }
               ],
               "url": {
                 "raw": "{{sms_base_url}}/sms/status/bulk?models.QueryParamGroupID={{last_group_id}}",
                 "host": ["{{sms_base_url}}"],
                 "path": ["sms", "status", "bulk"],
                 "query": [
                   {
                     "key": "models.QueryParamGroupID",
                     "value": "{{last_group_id}}"
                   }
                 ]
               }
             }
           }
         ]
       },
       {
         "name": "Admin",
         "item": [
           {
             "name": "Dashboard",
             "request": {
               "method": "GET",
               "header": [
                 {
                   "key": "Authorization",
                   "value": "Bearer {{jwt_token}}"
                 }
               ],
               "url": {
                 "raw": "{{auth_base_url}}/admin/dashboard",
                 "host": ["{{auth_base_url}}"],
                 "path": ["admin", "dashboard"]
               }
             }
           }
         ]
       }
     ]
   }
   ```

---

   ## Security Considerations

   1. **Password Requirements**: Implement strong password policies on the client side
   2. **JWT Expiration**: Tokens should have reasonable expiration times
   3. **HTTPS**: Always use HTTPS in production environments
   4. **Rate Limiting**: Monitor and respect rate limits
   5. **CAPTCHA**: Use secure login endpoint for public-facing applications
   6. **Role Validation**: Ensure proper role-based access control
   7. **Input Sanitization**: Validate and sanitize all user inputs
   8. **SQL Injection**: Use parameterized queries
   9. **XSS Protection**: Encode output data
   10. **CORS**: Configure CORS appropriately for your domains

---

   ## Support and Contact

   **Technical Support:** support@example.com

   **Developer Portal:** https://docs.example.com

   **GitHub:** https://github.com/shpd

   **Issue Tracker:** https://github.com/shpd/auth-service/issues

---

   ## License

   This API is licensed under the MIT License. See https://opensource.org/licenses/MIT for details.

---

   ## Changelog

   ### Version 1.0.0 (Current)
   - Initial release
   - Complete authentication system with JWT
   - MSISDN mapping service
   - SMS sending and tracking
   - Prefix management system
   - CAPTCHA integration
   - Admin dashboard

   ### Upcoming Features
   - Webhook support for SMS delivery notifications
   - Two-factor authentication (2FA)
   - API usage analytics dashboard
   - Batch operations for user management
   - Advanced filtering for SMS history
   - WebSocket support for real-time notifications

---

   ## End of Documentation

   Thank you for using our API! For questions or support, please contact us at [support@shpdco.ir](mailto:support@shpdco.ir)

