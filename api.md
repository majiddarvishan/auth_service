# REST Web Service Documentation

## 1. Overview

Provide a high-level description of the REST web service.

* **Service Name:**
* **Version:**
* **Base URL:**
* **Protocol:** HTTP / HTTPS
* **Content Type:** application/json
* **Audience:** Developers / Integrators

---

## Quick Start

A minimal guide for new clients to start using the API.

### Step 1: Obtain Credentials

* Create an account
* Generate an API key or obtain JWT credentials

### Step 2: Send Your First Request

Example using API Key authentication:

```
POST /api/v1/sms/send
X-API-Key: YOUR_API_KEY
Content-Type: application/json
```

```json
{
  "to": "+989121234567",
  "message": "Hello world"
}
```

### Step 3: Check Message Status

```
GET /api/v1/sms/status/msg_123456
X-API-Key: YOUR_API_KEY
```

### Step 4: Handle Errors & Retries

* Retry on `5xx` errors
* Do not retry on `4xx` errors except `429`

---

## 2. Architecture Overview

Briefly explain the system architecture.

* High-level diagram (optional)
* Major components
* External dependencies (DB, cache, message broker, SMPP, etc.)

---

## 3. Authentication & Authorization

Describe how clients authenticate and are authorized.

### 3.1 Authentication Method

* API Key / JWT / OAuth2 / Basic Auth
* Header or parameter used

**Example:**

```
Authorization: Bearer <JWT_TOKEN>
```

### 3.2 Roles & Permissions

| Role  | Description    |
| ----- | -------------- |
| admin | Full access    |
| user  | Limited access |

---

## 4. Common Request Headers

| Header        | Required | Description      |
| ------------- | -------- | ---------------- |
| Content-Type  | Yes      | application/json |
| Authorization | Yes      | Auth token       |

---

## 5. Common Response Format

Standard response structure used across APIs.

```json
{
  "success": true,
  "data": {},
  "error": null
}
```

---

## 6. Error Handling

List error formats and codes.

### 6.1 Error Response Format

```json
{
  "success": false,
  "error": {
    "code": "INVALID_REQUEST",
    "message": "Description of error"
  }
}
```

### 6.2 HTTP Status Codes

| Code | Meaning               |
| ---- | --------------------- |
| 200  | Success               |
| 400  | Bad Request           |
| 401  | Unauthorized          |
| 403  | Forbidden             |
| 404  | Not Found             |
| 409  | Idempotency conflict  |
| 429  | Rate limit exceeded   |
| 500  | Internal Server Error |

---

## 6.3 Idempotency Keys

Idempotency keys allow clients to safely retry requests without creating duplicate operations.

### When to Use

Idempotency is **required** for endpoints that create resources or trigger actions, such as:

* Sending SMS
* Bulk SMS

### How It Works

* Client generates a unique key per request
* Key is sent via HTTP header
* Server stores the result for a limited time
* Repeated requests with the same key return the **original response**

### Header Usage

```
Idempotency-Key: <UNIQUE_KEY>
```

### Rules

* Keys must be unique per operation
* Recommended format: UUID v4
* Keys expire after a configurable time window (e.g. 24 hours)

### Error Behavior

| Scenario                    | Response                   |
| --------------------------- | -------------------------- |
| Same key, same payload      | `200 OK` (cached response) |
| Same key, different payload | `409 Conflict`             |

---

## 6.4 API-Key–Based Authentication

Some endpoints support API-key authentication instead of JWT.

### API Key Usage

* API key must be sent via HTTP header

```
X-API-Key: <YOUR_API_KEY>
```

### API Key Scopes

API keys are created with one or more **scopes** that define what actions are allowed.

| Scope        | Description                                     |
| ------------ | ----------------------------------------------- |
| sms:send     | Send single or bulk SMS                         |
| sms:status   | Read message delivery status                    |
| account:read | Read account balance and profile                |
| reports:read | Access reports and analytics                    |
| admin:manage | Manage users, API keys, and limits (admin only) |

### Scope Enforcement

* Requests using an API key without the required scope will be rejected
* Scope violations return HTTP `403 Forbidden`

---

### Scope-to-Endpoint Mapping

The table below shows which scopes are required for each endpoint.

| Endpoint                          | Method | Required Scope |
| --------------------------------- | ------ | -------------- |
| `/api/v1/sms/send`                | POST   | `sms:send`     |
| `/api/v1/sms/bulk`                | POST   | `sms:send`     |
| `/api/v1/sms/status/{message_id}` | GET    | `sms:status`   |
| `/api/v1/account/balance`         | GET    | `account:read` |
| `/api/v1/reports/messages`        | GET    | `reports:read` |
| `/api/v1/admin/api-keys`          | POST   | `admin:manage` |
| `/api/v1/admin/users`             | GET    | `admin:manage` |

---

--|--------|
| 200 | Success |
| 400 | Bad Request |
| 401 | Unauthorized |
| 403 | Forbidden |
| 404 | Not Found |
| 500 | Internal Server Error |

---

## 7. API Endpoints

### API Separation

The API is divided into two main groups:

* **Public API**: Used by client applications to send messages and query data
* **Admin API**: Used by administrators to manage users, API keys, and system settings

All admin endpoints are prefixed with `/api/v1/admin` and require an **Admin JWT token**.

---

## 7A. Public API Endpoints

Below are example endpoints to demonstrate documentation style and level of detail.

---

### 7.1 Authentication – Login

Authenticate a user and return a JWT token.

* **URL:** `/api/v1/auth/login`
* **Method:** POST
* **Auth Required:** No

#### Request

```json
{
  "username": "user@example.com",
  "password": "strongPassword123"
}
```

#### Response

```json
{
  "success": true,
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
    "expires_in": 3600
  }
}
```

---

### 7.2 Send SMS (API Key Auth)

Send a single SMS using API-key authentication.

* **URL:** `/api/v1/sms/send`
* **Method:** POST
* **Auth Required:** Yes (API Key)

#### Request Headers

```
X-API-Key: <YOUR_API_KEY>
```

#### Request Body

```json
{
  "to": "+989121234567",
  "message": "Hello from API-key auth"
}
```

#### Response

```json
{
  "success": true,
  "data": {
    "message_id": "msg_123456",
    "status": "QUEUED"
  }
}
```

---

## 7B. Admin API Endpoints

### 7.3 Admin – Create API Key

Create a new API key for a user with specific scopes.

* **URL:** `/api/v1/admin/api-keys`
* **Method:** POST
* **Auth Required:** Yes (Admin JWT)

#### Request

```json
{
  "user_id": "user_123",
  "name": "billing-service",
  "scopes": ["sms:send", "sms:status", "account:read"]
}
```

#### Response

```json
{
  "user_id": "user_123",
  "name": "billing-service"
}
```

#### Response

```json
{
  "success": true,
  "data": {
    "api_key": "sk_live_xxxxxxxxx",
    "created_at": "2025-07-12T09:30:00Z"
  }
}
```

---

### 7.4 Admin – List Users

Retrieve all registered users.

* **URL:** `/api/v1/admin/users`
* **Method:** GET
* **Auth Required:** Yes (Admin JWT)

#### Response

```json
{
  "success": true,
  "data": [
    {
      "id": "user_123",
      "email": "user@example.com",
      "role": "user",
      "status": "active"
    }
  ]
}
```

---

## 8. Pagination

(If Applicable)
Explain pagination strategy.

* Query parameters: `page`, `limit`

```json
{
  "page": 1,
  "limit": 20,
  "total": 120
}
```

---

## 9. Rate Limiting

Describe rate limits.

* Requests per minute
* Headers returned

Example:

```
X-RateLimit-Limit: 1000
X-RateLimit-Remaining: 980
```

---

## 10. Webhooks (If Any)

Describe webhook events.

* Event types
* Payload format
* Retry policy

---

## 11. Versioning Strategy

Explain how API versions are handled.

* URL-based: `/v1/`
* Header-based

---

## 11. CLI / cURL Examples

Practical command-line examples for quick testing and automation.

---

### Login and Get JWT Token

```bash
curl -X POST https://api.example.com/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "user@example.com",
    "password": "strongPassword123"
  }'
```

---

### Send SMS Using API Key (Scoped)

This request requires the `sms:send` scope and supports idempotency.

````bash
curl -X POST https://api.example.com/api/v1/sms/send \
  -H "Content-Type: application/json" \
  -H "X-API-Key: YOUR_API_KEY" \
  -H "Idempotency-Key: 550e8400-e29b-41d4-a716-446655440000" \
  -d '{
    "to": "+989121234567",
    "message": "Hello with idempotency"
  }'
```bash
curl -X POST https://api.example.com/api/v1/sms/send \
  -H "Content-Type: application/json" \
  -H "X-API-Key: YOUR_API_KEY" \
  -d '{
    "to": "+989121234567",
    "message": "Hello from curl"
  }'
````

---

### Check Message Status

```bash
curl -X GET https://api.example.com/api/v1/sms/status/msg_123456 \
  -H "X-API-Key: YOUR_API_KEY"
```

---

### Admin: Create API Key

```bash
curl -X POST https://api.example.com/api/v1/admin/api-keys \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user_123",
    "name": "crm-system"
  }'
```

---

## 12. Environment Configuration

| Environment | Base URL                                                   |
| ----------- | ---------------------------------------------------------- |
| Dev         | [https://dev.api.example.com](https://dev.api.example.com) |
| Prod        | [https://api.example.com](https://api.example.com)         |

---

## 13. Changelog

Track changes between versions.

| Version | Date       | Description     |
| ------- | ---------- | --------------- |
| 1.0.0   | YYYY-MM-DD | Initial release |

---

## 14. Contact & Support

* Maintainer:
* Email:
* Issue Tracker:

---

## 15. Appendix

Any additional notes, examples, or diagrams.
