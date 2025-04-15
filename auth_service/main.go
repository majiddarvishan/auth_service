package main

import (
    "log"
    "net/http"
    "net/http/httputil"
    "net/url"
    "os"
    "strings"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v4"
    "github.com/joho/godotenv"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

var secretKey string
var db *gorm.DB

// User model stores users in the PostgreSQL database.
type User struct {
    gorm.Model
    Username string `gorm:"uniqueIndex"`
    Password string
    Role     string
}

// initDB initializes the database connection and migrates the schema.
// The DSN is read from the environment variable DATABASE_URL.
// Example DATABASE_URL:
//   "host=localhost user=postgres password=postgres dbname=mydb port=5432 sslmode=disable TimeZone=Asia/Shanghai"
func initDB() {
    dsn := os.Getenv("DATABASE_URL")
    if dsn == "" {
        log.Fatal("DATABASE_URL is not set in .env file")
    }

    var err error
    db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to the database: ", err)
    }

    if err := db.AutoMigrate(&User{}); err != nil {
        log.Fatal("Failed to auto migrate user schema: ", err)
    }
}

// hashPassword hashes a plaintext password using bcrypt.
func hashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    if err != nil {
        return "", err
    }
    return string(bytes), nil
}

// checkPasswordHash compares a plaintext password with a bcrypt hash.
func checkPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}

// generateToken creates a JWT token with a 24-hour expiration.
func generateToken(username string) (string, error) {
    claims := jwt.MapClaims{
        "sub":  username,
        "exp":  time.Now().Add(time.Hour * 24).Unix(),
        "role": "admin", // For this example, all users are given the "admin" role.
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(secretKey))
}

// registerHandler allows new users to register. It hashes their password
// and then stores the user record in the database.
func registerHandler(c *gin.Context) {
    type RegisterRequest struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }

    var req RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
        return
    }

    if req.Username == "" || req.Password == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Username and password cannot be empty"})
        return
    }

    hashedPassword, err := hashPassword(req.Password)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not hash password"})
        return
    }

    user := User{
        Username: req.Username,
        Password: hashedPassword,
        Role:     "admin",
    }

    if err := db.Create(&user).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user", "details": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

// loginHandler authenticates the user by checking the database record,
// comparing the provided password with the hashed password, and returns a JWT.
func loginHandler(c *gin.Context) {
    type LoginRequest struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }

    var req LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
        return
    }

    var user User
    if err := db.Where("username = ?", req.Username).First(&user).Error; err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
        return
    }

    if !checkPasswordHash(req.Password, user.Password) {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
        return
    }

    token, err := generateToken(user.Username)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"token": token})
}

// AuthMiddleware validates the JWT token for protected routes.
func AuthMiddleware(c *gin.Context) {
    authHeader := c.GetHeader("Authorization")
    if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        c.Abort()
        return
    }

    tokenString := strings.TrimPrefix(authHeader, "Bearer ")
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, jwt.ErrInvalidKey
        }
        return []byte(secretKey), nil
    })

    if err != nil || !token.Valid {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
        c.Abort()
        return
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok || claims["role"] != "admin" {
        c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
        c.Abort()
        return
    }

    c.Next()
}

// proxyToFinalService forwards incoming (protected) requests to the Final-Service.
func proxyToFinalService(c *gin.Context) {
    finalServiceURL := "http://localhost:8081" // Change if your Final-Service is hosted elsewhere.
    remote, err := url.Parse(finalServiceURL)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid Final-Service URL"})
        return
    }

    proxy := httputil.NewSingleHostReverseProxy(remote)
    c.Request.Host = remote.Host
    proxy.ServeHTTP(c.Writer, c.Request)
}

func main() {
    // Load environment variables from .env.
    if err := godotenv.Load(); err != nil {
        log.Fatal("Error loading .env file")
    }

    secretKey = os.Getenv("SECRET_KEY")
    if secretKey == "" {
        log.Fatal("SECRET_KEY is not set in .env file")
    }

    initDB()

    r := gin.Default()

    // PUBLIC ROUTES:
    r.POST("/register", registerHandler)
    r.POST("/login", loginHandler)

    // PROTECTED ROUTES:
    // Create a separate Gin engine for protected routes.
    protected := gin.New()
    protected.Use(AuthMiddleware)
    // Use a catch-all route to forward any unmatched requests.
    protected.Any("/*path", proxyToFinalService)

    // Delegate unmatched routes from the main router to the protected engine.
    r.NoRoute(func(c *gin.Context) {
        protected.HandleContext(c)
    })

    // Run the Auth-Service on port 8080.
    r.Run(":8080")
}
