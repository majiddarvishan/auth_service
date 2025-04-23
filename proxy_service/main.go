package main

import (
	"auth_service/config"
	"auth_service/database"
	"auth_service/routes"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// @title Auth service API
// @version 1.0
// @description A Auth-service gateway.
// @termsOfService https://example.com/terms

// @contact.name Majid Darvishan
// @contact.url https://github.com/shpd
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /
func main() {
	// Load configuration from .env.
	config.LoadConfig()

	// Initialize the database.
	database.InitDB()

	// Setup routes.
	r := routes.SetupRoutes()

	// Run the server on port 8080.
	// r.Run(":8080")

    // To ensure all HTTP requests are redirected to HTTPS, run a separate HTTP server.
    go func() {
        // HTTP server: Listen on port 8080 and redirect all traffic to HTTPS on port 8443.
        httpRouter := gin.Default()
        httpRouter.GET("/", func(c *gin.Context) {
            // Clean host if it includes the HTTP port.
            host := c.Request.Host
            // Remove port if present
            if colonPos := strings.Index(host, ":"); colonPos != -1 {
                host, _, _ = net.SplitHostPort(host)
            }
            c.Redirect(http.StatusMovedPermanently, "https://" + host + ":8443" + c.Request.RequestURI)
        })
        if err := httpRouter.Run(":8080"); err != nil {
            log.Fatal("HTTP redirection server failed:", err)
        }
    }()

    // Run HTTPS server with SSL certificate on port 8443.
    if err := r.RunTLS(":8443", "cert.pem", "key.pem"); err != nil {
        log.Fatal("Failed to start HTTPS server:", err)
    }
}
