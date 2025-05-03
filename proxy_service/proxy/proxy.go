package proxy

import (
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
)

// NewMultiTargetReverseProxy creates a reverse proxy that load-balances requests
// among the provided target URLs.
func NewMultiTargetReverseProxy(targets []*url.URL) *httputil.ReverseProxy {
    // Create a new random generator using a custom seed.
    rng := rand.New(rand.NewSource(time.Now().UnixNano()))

    // Custom director function to choose a target per request.
    director := func(req *http.Request) {
        // Pick a target randomly.
        target := targets[rng.Intn(len(targets))]

        // Update the scheme and host of the request to match the selected target.
        req.URL.Scheme = target.Scheme
        req.URL.Host = target.Host

        // Optionally, you can adjust the path if you don't want to append the original request path.
        // For example, to use only the target's path:
        // req.URL.Path = target.Path
        // Or, to preserve the incoming request's path, you can leave it as is.

        // If the target specifies a base path, you might need to join them properly.
    }

    return &httputil.ReverseProxy{Director: director}
}

func ProxyToEndpoint(c *gin.Context, targetEndpoints []string) {
    targets := []*url.URL{}

    for _, target := range targetEndpoints {
        t, err := url.Parse(target)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid target endpoint"})
            return
        }

        targets = append(targets, t)
    }

    proxy := NewMultiTargetReverseProxy(targets)
	proxy.ServeHTTP(c.Writer, c.Request)
}


// func ProxyToEndpoint(c *gin.Context, targetEndpoint string) {
// 	targetURL, err := url.Parse(targetEndpoint)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid target endpoint"})
// 		return
// 	}

// 	proxy := httputil.NewSingleHostReverseProxy(targetURL)
// 	proxy.ServeHTTP(c.Writer, c.Request)
// }
