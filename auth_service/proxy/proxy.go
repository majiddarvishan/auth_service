package proxy

import (
    "net/http"
    "net/http/httputil"
    "net/url"

    "github.com/gin-gonic/gin"
)

// ProxyRequest forwards the incoming request to the Final-Service.
func ProxyRequest(c *gin.Context) {
    finalServiceURL := "http://localhost:8081" // Adjust if needed.
    remote, err := url.Parse(finalServiceURL)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid final service URL"})
        return
    }

    proxy := httputil.NewSingleHostReverseProxy(remote)
    // Update the Host header so the final service correctly receives the request.
    c.Request.Host = remote.Host
    proxy.ServeHTTP(c.Writer, c.Request)
}
