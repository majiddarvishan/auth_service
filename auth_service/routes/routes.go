package routes

import (
	"auth_service/handlers"
	"auth_service/middleware"
	"auth_service/proxy"
	// "net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRoutes configures and returns the Gin engine.
func SetupRoutes() *gin.Engine {
	r := gin.Default()

	// Enable CORS for frontend requests.
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Allow all origins (or specify "http://localhost:3000")
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type", "Accept"},
		AllowCredentials: true,
	}))

	// r.OPTIONS("/*path", func(c *gin.Context) {
	// 	c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
	// 	c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	// 	c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type, Accept")
	// 	c.Header("Access-Control-Allow-Credentials", "true")
	// 	c.JSON(http.StatusOK, gin.H{"message": "CORS preflight OK"})
	// })

	// PUBLIC ROUTES:
	r.POST("/register", handlers.RegisterHandler)
	r.POST("/login", handlers.LoginHandler)

	r.GET("/admin",
		middleware.AuthMiddleware,          // Ensure user is authenticated.
		middleware.RoleMiddleware("admin"), // Ensure only admins can access.
		handlers.AdminDashboardHandler,
	)

	// DELETE User Endpoint (Admin Only)
	r.DELETE("/user/:username",
		middleware.AuthMiddleware,
		middleware.RoleMiddleware("admin"),
		handlers.DeleteUserHandler,
	)

	// Update User Role (Admin Only)
	r.PUT("/user/:username/role",
		middleware.AuthMiddleware,
		middleware.RoleMiddleware("admin"),
		handlers.UpdateUserRoleHandler,
	)

	// Create New Role (Admin Only)
	r.POST("/roles",
		middleware.AuthMiddleware,
		middleware.RoleMiddleware("admin"),
		handlers.CreateRoleHandler,
	)

	// Create or Update Accounting Rules (Admin Only)
	r.POST("/accounting_rules",
		middleware.AuthMiddleware,
		middleware.RoleMiddleware("admin"),
		handlers.CreateOrUpdateAccountingRuleHandler,
	)

	// SMS endpoints group: any request starting with "/sms/"
	// This group uses the standard authentication and dynamic accounting middleware,
	// and then uses a specialized proxy handler (SMSProxyRequest) to forward the request
	// to the final endpoint.
	r.Any("/sms/*path",
		middleware.AuthMiddleware,
		middleware.DynamicAccountingMiddleware,
		proxy.SMSProxyRequest,
	)

	// Accounting endpoints group: Any request starting with /accounting/* gets redirected to the accounting service.
	r.Any("/accounting/*path",
		middleware.AuthMiddleware,
		proxy.AccountingProxyRequest,
	)

	// Fallback route:
	// Use NoRoute to catch any requests that did not match the above routes.
	// This chain enforces that the request must be authenticated,
	// checked against an accounting rule, and then forwarded to the final component.
	// r.NoRoute(
	//     middleware.AuthMiddleware,              // Ensures a valid JWT is present.
	//     middleware.DynamicAccountingMiddleware, // Checks and deducts balance based on the rule.
	//     proxy.ProxyRequest,                     // Forwards the request to the final component.
	// )

	// (Optional) You can also add specific endpoints that require dynamic billing.
	// For example:
	// r.GET("/premium_data",
	//  middleware.AuthMiddleware,
	//  middleware.DynamicAccountingMiddleware,
	//  func(c *gin.Context) {
	//      c.JSON(http.StatusOK, gin.H{"message": "Premium data accessed"})
	//  },
	// )

	// SMS endpoint: check balance and, if sufficient, forward the request to the final component.
	// Note: We are reusing the generic proxy handler and not hardcoding any SMS logic.
	// r.POST("/sms",
	// 	middleware.AuthMiddleware,
	// 	middleware.ChargeUserMiddleware(5), // Charge $5 for sending an SMS
	// 	proxy.ProxyRequest,                 // Forwards the entire request to config.FinalEndpoint (e.g. http://localhost:8081)
	// )

	// Example Admin-only endpoint.
	// r.GET("/admin",
	// 	middleware.AuthMiddleware,
	// 	middleware.RoleMiddleware("admin"),
	// 	func(c *gin.Context) {
	// 		c.JSON(200, gin.H{"message": "Welcome, Admin!"})
	// 	},
	// )

	// Global Dynamic Accounting Middleware:
	// This middleware checks for the existence of an accounting rule for the incoming path.
	// If a rule exists, it will verify that the user's balance is sufficient and deduct the charge.
	// Otherwise, it will simply let the request pass.
	// r.Use(middleware.DynamicAccountingMiddleware)
	// r.Use(middleware.DynamicAccountingMiddleware)

	// IMPORTANT: The endpoint below must be protected by AuthMiddleware so token claims are set.
	// Here, the request first runs through AuthMiddleware, then DynamicAccountingMiddleware,
	// and finally is forwarded via the generic proxy handler.
	// r.Use(
	//     middleware.AuthMiddleware,               // Decodes the JWT and sets claims in context.
	//     middleware.DynamicAccountingMiddleware,  // Now the token claims are found!
	//     proxy.ProxyRequest,                      // Forwards the request to the final component.
	// )

	// Fallback: all other routes are forwarded to the final component.
	// r.Any("/*path", proxy.ProxyRequest)

	// PROTECTED ROUTES:
	// Create a separate engine for protected routes.
	// protected := gin.New()
	// protected.Use(middleware.AuthMiddleware)
	// // Catch-all route: forward any unmatched requests to the final component.
	// protected.Any("/*path", proxy.ProxyRequest)
	// r.NoRoute(func(c *gin.Context) {
	// 	protected.HandleContext(c)
	// })

	return r
}
