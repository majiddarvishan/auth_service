package handlers

import (
    "net/http"

    // "auth_service/database"
    "github.com/gin-gonic/gin"
)

// CreateOrUpdateAccountingRuleHandler allows an admin to define balance-check rules for specific endpoints.
func CreateOrUpdateAccountingRuleHandler(c *gin.Context) {
    type AccountingRuleRequest struct {
        Endpoint string  `json:"endpoint"`
        Charge   float64 `json:"charge"`
    }

    var req AccountingRuleRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
        return
    }

    if req.Endpoint == "" || req.Charge <= 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Valid endpoint and charge amount are required"})
        return
    }

    // Find existing rule or create a new one.
    // var rule database.AccountingRule
    // if err := database.DB.Where("endpoint = ?", req.Endpoint).First(&rule).Error; err != nil {
    //     // If the rule doesn't exist, create a new one.
    //     rule = database.AccountingRule{
    //         Endpoint: req.Endpoint,
    //         Charge:   req.Charge,
    //     }
    //     if err := database.DB.Create(&rule).Error; err != nil {
    //         c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create accounting rule", "details": err.Error()})
    //         return
    //     }
    // } else {
    //     // Update existing rule.
    //     rule.Charge = req.Charge
    //     if err := database.DB.Save(&rule).Error; err != nil {
    //         c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update accounting rule", "details": err.Error()})
    //         return
    //     }
    // }

    // c.JSON(http.StatusOK, gin.H{"message": "Accounting rule saved successfully", "rule": rule})
}
