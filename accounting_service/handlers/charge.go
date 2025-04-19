package handlers

import (
    "net/http"

    "accounting_service/database"
    "github.com/gin-gonic/gin"
)

// ChargeRequest is the payload expected by the accounting service.
type ChargeRequest struct {
    Username string `json:"username"`
    Endpoint string `json:"endpoint"`
}

func ChargeHandler(c *gin.Context) {
    var req ChargeRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
        return
    }

    // Look up the accounting rule for the endpoint.
    var rule database.AccountingRule
    if err := database.DB.Where("endpoint = ?", req.Endpoint).First(&rule).Error; err != nil {
        c.JSON(http.StatusForbidden, gin.H{"error": "No accounting rule defined for this endpoint"})
        return
    }

    // Find the user.
    var user database.User
    if err := database.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    // Check and deduct.
    if user.Balance < rule.Charge {
        c.JSON(http.StatusPaymentRequired, gin.H{"error": "Insufficient balance"})
        return
    }

    user.Balance -= rule.Charge
    if err := database.DB.Save(&user).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update balance"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Charge deducted", "new_balance": user.Balance})
}

func UpdateUserChargeHandler(c *gin.Context) {
    username := c.Param("username")
    var req struct {
        Charge float64 `json:"charge"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
        return
    }
    var user database.User
    if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }
    user.Balance = req.Charge  // Make sure the User model has a Charge field; or if you meant Balance, update accordingly.
    if err := database.DB.Save(&user).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user charge"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "User charge updated successfully"})
}

