package proxy

import (
	"auth_service/config"
	"auth_service/database"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"

	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func toInt(i interface{}) (int, error) {
	switch v := i.(type) {
	case int:
		return v, nil
	case int8:
		return int(v), nil
	case int16:
		return int(v), nil
	case int32:
		return int(v), nil
	case int64:
		return int(v), nil
	case uint:
		return int(v), nil
	case uint8:
		return int(v), nil
	case uint16:
		return int(v), nil
	case uint32:
		return int(v), nil
	case uint64:
		return int(v), nil
	case float32:
		return int(v), nil
	case float64:
		return int(v), nil
	case string:
		// Try to parse string to int
		n, err := strconv.Atoi(v)
		if err != nil {
			return 0, fmt.Errorf("cannot convert string to int: %v", err)
		}
		return n, nil
	default:
		return 0, fmt.Errorf("unsupported type: %T", v)
	}
}

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
		req.Host = target.Host

		newPath := ""
		if strings.HasPrefix(req.URL.Path, config.BaseApi) {
			newPath = strings.TrimPrefix(req.URL.Path, config.BaseApi)
		}

		req.URL.Path = target.Path + newPath

		// Remove trailing slash if present
		req.URL.Path = strings.TrimSuffix(req.URL.Path, "/")

		// Strip hop-by-hop and forwarding headers
		req.Header.Del("X-Forwarded-For")
		req.Header.Del("X-Real-IP")
		req.Header.Del("Forwarded") // RFC-7239
		req.Header.Del("Via")

		authHeader := req.Header.Get("Authorization")
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		parts := strings.Split(tokenStr, ".")
		if len(parts) != 3 {
			fmt.Println("Invalid JWT format")
			return
		}

		// Decode the payload (second part)
		payload := parts[1]

		// JWT uses base64url encoding, which is slightly different from standard base64
		decoded, err := base64.RawURLEncoding.DecodeString(payload)
		if err != nil {
			fmt.Println("Error decoding payload:", err)
			return
		}

		// Convert JSON payload to a map
		var claims map[string]interface{}
		if err := json.Unmarshal(decoded, &claims); err != nil {
			fmt.Println("Error unmarshaling JSON:", err)
			return
		}

		for key, value := range claims {
			fmt.Printf("  %s: %v\n", key, value)
			if key == "user" {
				val, err := toInt(value)
				if err != nil {
					fmt.Printf("invalid user id %v\n", value)
					return
				}

				user, err := database.DB.GetUserByID(uint(val))
				if err != nil {
					fmt.Printf("Error in token: %v\n", err)
					return
				}

				req.Header.Del("Authorization")
				req.Header.Add("user-id", strconv.Itoa(val))
                req.Header.Add("user-name", user.Username)

				break
			}
		}
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
