package proxy

import (
    "net/http"
    "net/http/httputil"
    "net/url"

    "github.com/gin-gonic/gin"
)

// ProxyToFinalService forwards a request to the Final-Service.
func ProxyToFinalService(c *gin.Context) {
    finalServiceURL := "http://localhost:8081" // Change this URL if your Final-Service is hosted elsewhere.
    remote, err := url.Parse(finalServiceURL)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid Final-Service URL"})
        return
    }

    proxy := httputil.NewSingleHostReverseProxy(remote)
    // Update the request's Host header to match the final service.
    c.Request.Host = remote.Host
    proxy.ServeHTTP(c.Writer, c.Request)
}
