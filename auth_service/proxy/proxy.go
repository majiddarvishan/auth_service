package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"auth_service/config"

	"github.com/gin-gonic/gin"
)

// ProxyRequest forwards the incoming request to the Final-Service.
func ProxyRequest(c *gin.Context) {
	finalEndpoint := config.FinalEndpoint

	remote, err := url.Parse(finalEndpoint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid final service URL"})
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)
	// Update the Host header so the final service correctly receives the request.
	c.Request.Host = remote.Host
	proxy.ServeHTTP(c.Writer, c.Request)
}

func SMSProxyRequest(c *gin.Context) {
	finalEndpoint := config.FinalEndpoint

	remote, err := url.Parse(finalEndpoint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid final service URL"})
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)
	// Update the Host header so the final service correctly receives the request.
	c.Request.Host = remote.Host
	proxy.ServeHTTP(c.Writer, c.Request)
}
