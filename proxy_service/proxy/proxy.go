package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"auth_service/config"

	"github.com/gin-gonic/gin"
)

// AccountingProxyRequest forwards requests for accounting endpoints to the accounting service.
func AccountingProxyRequest(c *gin.Context) {
    // Use the accounting service endpoint from config.
    accountingEndpoint := config.AccountingEndpoint // e.g., "http://localhost:8082"
    remote, err := url.Parse(accountingEndpoint)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid accounting service URL"})
        return
    }

    // Remove the "/accounting" prefix so that the accounting service gets the appropriate path.
    c.Request.URL.Path = strings.TrimPrefix(c.Request.URL.Path, "/accounting")
    if c.Request.URL.Path == "" {
        c.Request.URL.Path = "/"
    }
    c.Request.Host = remote.Host

    // Create and use the reverse proxy.
    proxy := httputil.NewSingleHostReverseProxy(remote)
    proxy.ServeHTTP(c.Writer, c.Request)
}

func SMSProxyRequest(c *gin.Context) {
	smsEndpoint := config.SmsEndpoint

	remote, err := url.Parse(smsEndpoint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid final service URL"})
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)
	// Update the Host header so the final service correctly receives the request.
	c.Request.Host = remote.Host
	proxy.ServeHTTP(c.Writer, c.Request)
}
