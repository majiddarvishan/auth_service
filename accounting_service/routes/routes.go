package routes

import (
	"accounting_service/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	r := gin.Default()

	r.POST("/accounting/charge", handlers.ChargeHandler)
	r.PUT("/accounting/user/:username/charge", handlers.UpdateUserChargeHandler)

	return r
}
