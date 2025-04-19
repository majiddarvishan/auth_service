package main

import (
    "accounting_service/config"
    "accounting_service/database"
    "accounting_service/routes"
    "log"
)

func main() {
    config.LoadConfig()
    database.InitDB()
    r := routes.SetupRoutes()

    log.Println("Accounting service running on port " + config.Port)
    r.Run(":" + config.Port)
}
