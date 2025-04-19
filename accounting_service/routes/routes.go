package routes

import (
	"accounting_service/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	r := gin.Default()

	// an endpoint that checks and deducts a charge
	r.POST("/charge", handlers.ChargeHandler)

	r.PUT("/users/:username/charge", handlers.UpdateUserChargeHandler)

	return r
}
