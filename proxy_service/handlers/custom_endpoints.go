package handlers

import (
	"net/http"
	"time"

	"auth_service/constants"
	"auth_service/database"
	"auth_service/logger"
	"auth_service/middleware"
	"auth_service/proxy"
	"auth_service/types"
	"auth_service/validation"

	"github.com/gin-gonic/gin"
)

func registerCustomEndpointDynamic(r *gin.RouterGroup, ep *database.CustomEndpoint) {
	// Wrap the handler with the endpoint parameter.
	wrappedHandler := func(c *gin.Context) {
		proxy.ProxyToEndpoint(c, ep)
	}

	// Build the handler chain for the dynamic route.
	handlersChain := []gin.HandlerFunc{middleware.AuthMiddleware}
	if ep.NeedAccounting {
		handlersChain = append(handlersChain, middleware.DynamicAccountingMiddleware)
	}
	handlersChain = append(handlersChain, wrappedHandler)

	// Register based on the HTTP method.
	switch ep.Method {
	case "GET":
		r.GET(ep.Path, handlersChain...)
	case "POST":
		r.POST(ep.Path, handlersChain...)
	case "PUT":
		r.PUT(ep.Path, handlersChain...)
	case "DELETE":
		r.DELETE(ep.Path, handlersChain...)
	default:
		r.Any(ep.Path, handlersChain...)
	}
	logger.Info("Registered dynamic route", "path", ep.Path, "method", ep.Method, "target", ep.Endpoints[0])
}

// func RegisterCustomEndpoints(r *gin.Engine) {
func RegisterCustomEndpoints(routerGroup *gin.RouterGroup) {
	endpoints, err := database.DB.GetAllCustomEndpoints()
	if err != nil {
		logger.Error("Error fetching custom endpoints", "error", err.Error())
		return
	}

	for _, endpoint := range endpoints {
		registerCustomEndpointDynamic(routerGroup, &endpoint)
	}
}

// SwaggerCustomEndpoint represents the payload for Custom Endpoint.
// swagger:model SwaggerCustomEndpoint
// @Description SwaggerCustomEndpoint defines the expected request body for custom endpoint.
// @Property path body string true "Prefix of URI"
// @Property method body string true "method of request"
// @Property endpoints body []string true  "ID of the captcha challenge"
// @Property needAccounting body string true  "Needs check accouting before redirect it"
type SwaggerCustomEndpoint struct {
	Path           string
	Method         string
	Endpoints      []string
	NeedAccounting bool
}

// CreateCustomEndpointHandler create custom endpoint.
// @Summary      Create Custom Endpoint
// @Description  Create Custom Endpoint to redirect its requests to another endpoints
// @Tags         CustomEndpoints
// @Accept       json
// @Produce      json
// @Param        request  body      SwaggerCustomEndpoint  true  "CustomEndpoint payload"
// @Success      200      {object}  map[string]string  "JWT token"
// @Failure      400      {object}  map[string]string  "Invalid JSON format"
// @Failure      401      {object}  map[string]string  "Unauthorized: invalid credentials"
// @Failure      500      {object}  map[string]string  "Server error during token generation"
// @Router       /admin/custom-endpoints [post]
func CreateCustomEndpointHandler(dynamicGroup *gin.RouterGroup) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req database.CustomEndpoint
		requestID, _ := c.Get("request_id")

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, types.APIError{
				Code:      constants.ErrorInvalidJSON,
				Message:   "Invalid JSON payload",
				Timestamp: time.Now().Unix(),
				RequestID: requestID.(string),
			})
			return
		}

		// Validate path
		if err := validation.ValidateEndpointPath(req.Path); err != nil {
			c.JSON(http.StatusBadRequest, types.APIError{
				Code:      constants.ErrorInvalidPathTravers,
				Message:   err.Error(),
				Timestamp: time.Now().Unix(),
				RequestID: requestID.(string),
			})
			return
		}

		// Validate targets
		if err := validation.ValidateEndpointTargets(req.Endpoints); err != nil {
			c.JSON(http.StatusBadRequest, types.APIError{
				Code:      "INVALID_ENDPOINTS",
				Message:   err.Error(),
				Timestamp: time.Now().Unix(),
				RequestID: requestID.(string),
			})
			return
		}

		// Validate HTTP method if provided
		if req.Method != "" && req.Method != "ANY" {
			if err := validation.ValidateHTTPMethod(req.Method); err != nil {
				c.JSON(http.StatusBadRequest, types.APIError{
					Code:      "INVALID_METHOD",
					Message:   err.Error(),
					Timestamp: time.Now().Unix(),
					RequestID: requestID.(string),
				})
				return
			}

			if req.Method == "" {
				req.Method = "ANY"
			}

			req.Path += "/*path"
			req.Enabled = true

			if err := database.DB.CreateCustomEndpoint(&req); err != nil {
				c.JSON(http.StatusInternalServerError, types.APIError{
					Code:      constants.ErrorInternalServer,
					Message:   "Failed to create custom endpoint",
					Timestamp: time.Now().Unix(),
					RequestID: requestID.(string),
				})
				return
			}

			c.JSON(http.StatusOK, types.APISuccess{
				Message:   "Custom endpoint created successfully",
				Data:      req,
				Timestamp: time.Now().Unix(),
				RequestID: requestID.(string),
			})

			registerCustomEndpointDynamic(dynamicGroup, &req)

			c.Next()
		}
	}
}

// DeleteCustomEndpointHandler deletes a custom endpoint dynamically.
//
// @Summary      Delete a custom endpoint
// @Description  Deletes a previously registered custom endpoint by path and restores a 404 handler for it.
// @Tags         CustomEndpoints
// @Accept       json
// @Produce      json
// @Param        request body SwaggerCustomEndpoint true "Custom endpoint object (must include path)"
// @Success      200 {object} map[string]interface{} "Custom endpoint deleted successfully"
// @Failure      400 {object} map[string]interface{} "Invalid JSON payload"
// @Failure      500 {object} map[string]interface{} "Failed to delete or find custom endpoint"
// @Router       /admin/custom-endpoints [delete]
func DeleteCustomEndpointHandler(dynamicGroup *gin.RouterGroup) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req database.CustomEndpoint
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
			return
		}

		/* Since Delete is a write-only operation, it does not fetch or populate the fields of the deleted record.
		 * So we explicitly query the record first to retrive data.
		 */
		endPoint, err := database.DB.GetCustomEndpointByPath(req.Path + "/*path")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find custom endpoint", "path": req.Path})
			return
		}

		err = database.DB.DeleteCustomEndpointByPath(req.Path + "/*path")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete custom endpoint", "path": req.Path})
			return
		}

		// notFoundHandler := func(c *gin.Context) {
		// 	c.JSON(404, gin.H{"error": "Not found"})
		// 	c.Abort() // Important: stop the chain
		// }

		// switch endPoint.Method {
		// case "GET":
		// 	dynamicGroup.GET(endPoint.Path, notFoundHandler)
		// case "POST":
		// 	dynamicGroup.POST(endPoint.Path, notFoundHandler)
		// case "PUT":
		// 	dynamicGroup.PUT(endPoint.Path, notFoundHandler)
		// case "DELETE":
		// 	dynamicGroup.DELETE(endPoint.Path, notFoundHandler)
		// default:
		// 	dynamicGroup.Any(endPoint.Path, notFoundHandler)
		// }
		proxy.DeleteEndpoint(endPoint)

		c.JSON(http.StatusOK, gin.H{"message": "Custom endpoint deleted successfully", "path": req.Path})

		c.Next()
	}
}
