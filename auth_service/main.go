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
)

var secretKey string

// generateToken creates a JWT token with a 24-hour expiration.
func generateToken(username string) (string, error) {
    claims := jwt.MapClaims{
        "sub":  username,
        "exp":  time.Now().Add(time.Hour * 24).Unix(),
        "role": "admin", // Hardcoded role for example purposes
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(secretKey))
}

// loginHandler handles POST /login and returns a JWT when credentials are correct.
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

    // Basic credentials validation (replace with your own logic)
    if req.Username == "admin" && req.Password == "password" {
        token, err := generateToken(req.Username)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
            return
        }
        c.JSON(http.StatusOK, gin.H{"token": token})
        return
    }

    c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
}

// AuthMiddleware checks for a valid JWT token.
func AuthMiddleware(c *gin.Context) {
    authHeader := c.GetHeader("Authorization")
    if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        c.Abort()
        return
    }

    tokenString := strings.TrimPrefix(authHeader, "Bearer ")
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        // Validate signing algorithm
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

// proxyToFinalService forwards the incoming request to the Final-Service.
func proxyToFinalService(c *gin.Context) {
    finalServiceURL := "http://localhost:8081" // URL of your Final-Service
    remote, err := url.Parse(finalServiceURL)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid Final-Service URL"})
        return
    }

    proxy := httputil.NewSingleHostReverseProxy(remote)
    // Update the Request Host to the final service host.
    c.Request.Host = remote.Host
    proxy.ServeHTTP(c.Writer, c.Request)
}

func main() {
    // Load environment variables from .env file.
    if err := godotenv.Load(); err != nil {
        log.Fatal("Error loading .env file")
    }

    secretKey = os.Getenv("SECRET_KEY")
    if secretKey == "" {
        log.Fatal("SECRET_KEY is not set in .env file")
    }

    // Create the main Gin router.
    r := gin.Default()

    // Public routes (no authentication needed).
    r.POST("/login", loginHandler)

    // Create a separate Gin engine for protected routes.
    // This engine is not registered directly on r, so we avoid conflicts.
    protected := gin.New()
    // Apply our authentication middleware on the protected engine.
    protected.Use(AuthMiddleware)
    // Catch-all (wildcard) route to proxy requests to the Final-Service.
    protected.Any("/*path", proxyToFinalService)

    // In the main engine, use a NoRoute handler to delegate unmatched paths (i.e. protected routes)
    // to our protected engine.
    r.NoRoute(func(c *gin.Context) {
        // Delegate the current request to the protected engine.
        protected.HandleContext(c)
    })

    // Run the Auth-Service on port 8080.
    r.Run(":8080")
}
