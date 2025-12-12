# Auth Service - API Gateway & Proxy

A high-performance, production-ready Go-based authentication and authorization gateway with dynamic reverse proxy capabilities. Provides user registration/login, JWT token validation, RBAC (Role-Based Access Control), custom endpoint proxying, and API accounting/quotas.

## 🎯 Features

### 🔐 Authentication & Authorization
- **User Registration** with strong password validation (8+ chars, uppercase, lowercase, digit, special char)
- **JWT Token Generation** using golang-jwt/v5
- **Role-Based Access Control** (Admin, User, Guest)
- **Secure Password Hashing** with bcrypt
- **Token Expiration** - Configurable token lifetime

### 🛡️ Security
- **HTTPS/TLS Support** with configurable certificates
- **Security Headers** - HSTS, CSP, X-Frame-Options, X-Content-Type-Options
- **Rate Limiting** - 100 requests/minute per IP (configurable)
- **Input Validation** - Password strength, username format, path traversal prevention
- **CORS Protection** - Configurable allowed origins

### 📊 Observability
- **Structured Logging** - JSON-formatted logs with request tracing
- **Request ID Tracking** - UUID-based end-to-end request tracing
- **Health Checks** - `/health` and `/version` endpoints for load balancers
- **Error Standardization** - Consistent error responses with error codes

### 🔄 Proxy Features
- **Dynamic Reverse Proxy** - Create custom endpoints at runtime
- **Load Balancing** - Random load balancing across multiple targets
- **Custom Routes** - Admin-controlled API endpoint forwarding
- **Accounting Integration** - Optional quota/accounting service integration

### 🗄️ Database
- **PostgreSQL Support** with GORM ORM
- **Mock Database** for testing (no PostgreSQL needed)
- **Connection Pooling** with configurable pool size
- **Auto Migrations** - Schema auto-created from Go structs

## 🚀 Quick Start

### Prerequisites
- Go 1.24+
- PostgreSQL 12+ (optional, can use mock database)
- OpenSSL (for TLS certificate generation)

### 1. Clone & Setup

```bash
# Clone repository
git clone https://github.com/majiddarvishan/auth_service.git
cd auth_service/proxy_service

# Install dependencies
go mod download

# Create .env file
cat > .env << EOF
BASE_API=/api/v1
TLS_PATH=.
SECRET_KEY=$(openssl rand -hex 32)
TOKEN_EXPIRATION_PERIOD=24h
DB_HOST=127.0.0.1
DB_PORT=5432
DB_USER_NAME=postgres
DB_PASSWORD=postgres
DB_NAME=proxy_db
ACCOUNTING_ENDPOINT=http://localhost:8082
EOF
```

### 2. Generate TLS Certificate

**Option A: Self-Signed (Development)**
```bash
openssl req -x509 -newkey rsa:2048 -nodes -keyout localhost-key.pem -out localhost.pem -days 365
```

**Option B: Using mkcert (Recommended for Development)**
```bash
# Install mkcert
sudo apt install libnss3-tools mkcert

# Create certificate
mkcert -install
mkcert localhost 127.0.0.1 ::1

# Copy to project root
cp localhost.pem . && cp localhost-key.pem .
```

### 3. Run Service

**With Mock Database (No PostgreSQL needed):**
```bash
go run main.go -d mock
```

**With PostgreSQL:**
```bash
# Ensure PostgreSQL is running and database exists
createdb proxy_db

# Run service
go run main.go -d postgres
```

**Build for Production:**
```bash
go build -o auth_service main.go
./auth_service -d postgres
```

### 4. Test the Service

```bash
# Health check
curl http://localhost:8080/health

# Register user
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"username": "john", "password": "SecurePass123!"}'

# Login
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"username": "john", "password": "SecurePass123!"}'
```

See [EXAMPLES.md](EXAMPLES.md) for comprehensive API examples.

## 📁 Project Structure

```
auth_service/proxy_service/
├── main.go                 # Entry point, graceful shutdown, config validation
├── config/                 # Configuration & environment validation
│   ├── config.go          # Load .env variables
│   └── validator.go       # Validate required variables & formats
├── database/              # Data layer with Store interface pattern
│   ├── store.go           # Store interface definition
│   ├── pgstore.go         # PostgreSQL implementation
│   ├── mockstore.go       # In-memory mock implementation
│   ├── models.go          # GORM data models
│   └── database.go_       # Database initialization
├── handlers/              # HTTP endpoint handlers
│   ├── user.go            # Register, login, user management
│   ├── admin.go           # Admin dashboard & controls
│   ├── roles.go           # Role management
│   ├── custom_endpoints.go # Dynamic endpoint creation
│   ├── captcha.go         # CAPTCHA endpoints
│   ├── health.go          # Health check endpoints
│   └── user_test.go       # Test fixtures
├── middleware/            # Request processing middleware
│   ├── auth.go            # JWT token validation
│   ├── role.go            # Role-based access control
│   ├── request_id.go      # Request tracing with UUID
│   ├── security_headers.go # Security headers injection
│   ├── rate_limit.go      # Rate limiting (100 req/min)
│   ├── accounting.go      # Quota/accounting integration
│   └── cors.go            # CORS configuration
├── routes/                # Gin router setup
│   └── routes.go          # All routes & middleware registration
├── proxy/                 # Reverse proxy implementation
│   └── proxy.go           # Dynamic proxy with load balancing
├── logger/                # Structured logging
│   └── logger.go          # JSON logger with slog
├── constants/             # Centralized constants
│   └── constants.go       # Roles, error codes, claim keys
├── validation/            # Input validators
│   └── validation.go      # Password, username, path validators
├── types/                 # Shared data types
│   └── types.go           # APIError, APISuccess, responses
├── deployments/           # Docker & deployment scripts
│   ├── Dockerfile         # Production image
│   ├── docker-compose.yml # Docker compose setup
│   ├── init_db.sh         # Database initialization
│   └── add_dynamic_route.sh # Add custom endpoint script
├── test/                  # API testing files
│   ├── test.http          # HTTP client test requests
│   └── rest-client.environmentVariables.json # Test env vars
└── docs/                  # API documentation
    ├── swagger.json       # OpenAPI spec
    └── swagger.yaml       # OpenAPI spec (YAML)
```

## ⚙️ Configuration

### Environment Variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `BASE_API` | ✅ | `/api/v1` | API route prefix |
| `TLS_PATH` | ✅ | `.` | Path to certificate files |
| `SECRET_KEY` | ✅ | - | JWT signing key (32+ chars) |
| `TOKEN_EXPIRATION_PERIOD` | ✅ | `24h` | Token lifetime (Go duration) |
| `DB_HOST` | ⚠️ | - | PostgreSQL host (only if using postgres) |
| `DB_PORT` | ⚠️ | - | PostgreSQL port (only if using postgres) |
| `DB_USER_NAME` | ⚠️ | - | PostgreSQL user (only if using postgres) |
| `DB_PASSWORD` | ⚠️ | - | PostgreSQL password (only if using postgres) |
| `DB_NAME` | ⚠️ | - | PostgreSQL database (only if using postgres) |
| `ACCOUNTING_ENDPOINT` | ❌ | `http://localhost:8082` | Accounting service URL |
| `ALLOWED_CORS_ORIGINS` | ❌ | `localhost:3000,localhost:8080` | CORS allowed origins |
| `DB_MAX_OPEN_CONNS` | ❌ | `25` | Connection pool max size |
| `DB_MAX_IDLE_CONNS` | ❌ | `5` | Connection pool idle size |

**Note:** ⚠️ Only required when using `-d postgres`

### Generate SECRET_KEY

```bash
openssl rand -hex 32
# Output: f78973efc0c0664995e2bb055bb2cac6779597a5294685f069229c909358f54a
```

## 📡 API Endpoints

### Health & Info (No Auth)
- `GET /health` - Service health status & DB connectivity
- `GET /version` - Service version info

### Authentication (No Auth)
- `POST /api/v1/users` - Register new user
- `POST /api/v1/login` - Login (get JWT token)

### CAPTCHA (No Auth, CORS-enabled)
- `GET /api/v1/captcha/new` - Get new CAPTCHA ID
- `GET /api/v1/captcha/image/:id` - Get CAPTCHA image

### Protected Routes (JWT Auth)
- `GET /api/v1/admin` - Admin dashboard
- `GET /api/v1/admin/roles` - List roles
- `POST /api/v1/admin/roles` - Create role
- `GET /api/v1/admin/custom-endpoints` - List custom endpoints
- `POST /api/v1/admin/custom-endpoints` - Create custom endpoint
- `DELETE /api/v1/admin/custom-endpoints` - Delete custom endpoint

### Swagger UI
- `GET /swagger/index.html` - Interactive API documentation

See [EXAMPLES.md](EXAMPLES.md) for detailed request/response examples.

## 📚 Usage Examples

### 1. Register User

```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "password": "SecurePass123!"
  }'
```

**Response:**
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

### 2. Login & Get Token

```bash
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "password": "SecurePass123!"
  }'
```

**Response:**
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

### 3. Access Protected Route

```bash
curl -X GET http://localhost:8080/api/v1/admin \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Response:**
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

### 4. Create Custom Endpoint (Admin Only)

```bash
curl -X POST http://localhost:8080/api/v1/admin/custom-endpoints \
  -H "Authorization: Bearer ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "path": "/api/users",
    "endpoints": ["https://api.example.com/users"],
    "method": "GET"
  }'
```

**Response:**
```json
{
  "data": {
    "id": 1,
    "path": "/api/users",
    "endpoints": ["https://api.example.com/users"],
    "method": "GET"
  },
  "message": "Custom endpoint created successfully",
  "timestamp": 1702431000,
  "request_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

### 5. Check Health

```bash
curl http://localhost:8080/health
```

**Response:**
```json
{
  "status": "healthy",
  "database": "connected",
  "version": "1.0.0",
  "timestamp": 1702431000
}
```

**More examples:** See [EXAMPLES.md](EXAMPLES.md)

## 🔑 Password Requirements

Passwords must contain:
- ✅ At least 8 characters
- ✅ One uppercase letter (A-Z)
- ✅ One lowercase letter (a-z)
- ✅ One digit (0-9)
- ✅ One special character (!@#$%^&*)

**Example:** `SecurePass123!`

## 👤 Username Requirements

Usernames must:
- ✅ Be 3-32 characters long
- ✅ Contain only letters, numbers, and underscores
- ✅ Start with a letter

**Valid examples:** `john_doe`, `user123`, `Admin_User`
**Invalid examples:** `jd` (too short), `john@doe` (special char), `123user` (starts with number)

## 🗄️ Database Operations

### Initialize PostgreSQL

```bash
# Create database
createdb proxy_db

# Run with service (auto-migrations)
go run main.go -d postgres
```

### Backup Database
```bash
pg_dump -U postgres proxy_db > backup.sql
```

### Restore Database
```bash
psql -U postgres proxy_db < backup.sql
```

## ✅ Testing

### Unit Tests
```bash
# Run all tests (uses mock database)
go test ./...

# Run with coverage
go test -cover ./...

# Run specific test
go test -v ./handlers -run TestRegisterHandler
```

### Integration Tests
```bash
# Run with PostgreSQL
go test ./... -args -db postgres
```

### Manual Testing
```bash
# Use the testing script
chmod +x test_api.sh
./test_api.sh

# Or use curl/REST client
# See EXAMPLES.md for detailed examples
```

## 🐳 Docker Deployment

### Build Image
```bash
docker build -f deployments/Dockerfile -t auth-service:latest .
```

### Run Container
```bash
docker run -p 8080:8080 -p 8443:8443 \
  -e BASE_API=/api/v1 \
  -e SECRET_KEY=$(openssl rand -hex 32) \
  -e DB_HOST=postgres \
  -e DB_USER_NAME=postgres \
  -e DB_PASSWORD=postgres \
  -e DB_NAME=proxy_db \
  auth-service:latest
```

### Docker Compose
```bash
# Start all services (includes PostgreSQL)
docker-compose -f deployments/docker-compose.yml up

# Stop services
docker-compose -f deployments/docker-compose.yml down
```

## ⚡ Performance Considerations

- **Connection Pooling**: Configurable via `DB_MAX_OPEN_CONNS` (default: 25)
- **Rate Limiting**: 100 req/min per IP (configurable per route)
- **Request Timeout**: 30s graceful shutdown window
- **Reverse Proxy**: Random load balancing across targets
- **Logging**: Structured JSON (efficient for parsing & aggregation)
- **Memory**: ~50MB base, scales with connection pool

## 📝 Logging

All logs are structured JSON format for easy parsing and aggregation:

```json
{"time":"2025-12-13T00:50:00Z","level":"INFO","msg":"User registered","username":"john","user_id":1,"request_id":"uuid"}
```

### Parse Logs

```bash
# Show all login events
go run main.go -d mock 2>&1 | jq 'select(.msg | contains("Login"))'

# Show all errors
go run main.go -d mock 2>&1 | jq 'select(.level == "ERROR")'

# Filter by request ID
REQUEST_ID="550e8400-e29b-41d4-a716-446655440000"
go run main.go -d mock 2>&1 | jq "select(.request_id == \"$REQUEST_ID\")"

# Count events by level
go run main.go -d mock 2>&1 | jq -s 'group_by(.level) | map({level: .[0].level, count: length})'
```

## 🔒 Security Best Practices

✅ **Always use HTTPS** in production (configure proper certificates)
✅ **Store JWT tokens** in httpOnly cookies, not localStorage
✅ **Keep SECRET_KEY secure** and rotate periodically
✅ **Use strong database passwords** (20+ chars with mixed case/numbers/symbols)
✅ **Enable CORS only** for trusted origins
✅ **Monitor rate limit** metrics and adjust if needed
✅ **Log and audit** all authentication events
✅ **Implement password reset** with secure token (future feature)
✅ **Use CAPTCHA** on registration endpoint
✅ **Implement account lockout** after failed login attempts (future)
✅ **Enable WAF** (Web Application Firewall) in production

## 📋 API Response Format

All responses follow a consistent format:

### Success Response (200, 201)
```json
{
  "data": { /* endpoint-specific data */ },
  "message": "Operation successful",
  "timestamp": 1702431000,
  "request_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

### Error Response (400, 401, 403, 429, 500)
```json
{
  "code": "ERROR_CODE",
  "message": "Human-readable message",
  "details": null,
  "timestamp": 1702431000,
  "request_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

**Common Error Codes:**
- `USER_NOT_FOUND` - User doesn't exist (404)
- `USER_ALREADY_EXISTS` - Username taken (409)
- `INVALID_USERNAME` - Invalid username format (400)
- `INVALID_PASSWORD_STRENGTH` - Password too weak (400)
- `INVALID_PASSWORD` - Wrong password (401)
- `UNAUTHORIZED` - Missing/invalid token (401)
- `FORBIDDEN` - Insufficient permissions (403)
- `INVALID_PATH` - Path traversal detected (400)
- `RATE_LIMIT_EXCEEDED` - Too many requests (429)
- `INTERNAL_SERVER_ERROR` - Server error (500)

## 🐛 Troubleshooting

### Service won't start
```bash
# Check .env file exists and is valid
cat .env

# Verify TLS certificates exist
ls -la localhost.pem localhost-key.pem

# Check database connection (if using postgres)
psql -h 127.0.0.1 -U postgres -d proxy_db
```

### JWT token errors
```bash
# Token expired?
# Check TOKEN_EXPIRATION_PERIOD in .env (default: 24h)

# Invalid signature?
# Ensure SECRET_KEY matches (must be same 32+ hex chars)
```

### Database connection errors
```bash
# Reset database
dropdb proxy_db
createdb proxy_db
go run main.go -d postgres
```

### Rate limit issues
```bash
# Rate limiter resets every 60 seconds per IP
# Wait 60 seconds or restart service
# For testing, use -d mock to avoid DB issues
```

### Certificate errors
```bash
# Regenerate self-signed certificate
rm localhost.pem localhost-key.pem
openssl req -x509 -newkey rsa:2048 -nodes -keyout localhost-key.pem -out localhost.pem -days 365
```

## 📚 Documentation

- **[EXAMPLES.md](EXAMPLES.md)** - Complete API request/response examples
- **[QUICK_REFERENCE.md](QUICK_REFERENCE.md)** - Quick lookup for new features
- **[IMPLEMENTATION_SUMMARY.md](IMPLEMENTATION_SUMMARY.md)** - What was implemented
- **[.github/copilot-instructions.md](.github/copilot-instructions.md)** - Detailed architecture guide
- **API Docs** - Available at `/swagger/index.html` when service is running

## 🤝 Contributing

1. Fork repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open Pull Request

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details

## 📞 Support & Contact

- **Email:** support@example.com
- **Issues:** [GitHub Issues](https://github.com/majiddarvishan/auth_service/issues)
- **Security:** Report vulnerabilities to support@example.com

## 🙏 Acknowledgments

Built with:
- [Gin](https://github.com/gin-gonic/gin) - HTTP web framework
- [GORM](https://github.com/go-gorm/gorm) - ORM library
- [golang-jwt](https://github.com/golang-jwt/jwt) - JWT implementation
- [bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt) - Password hashing

---

**Made with ❤️ by Majid Darvishan**
