package handlers

import (
    "net/http"

    "auth_service/database"
    "github.com/gin-gonic/gin"
)

// CreateRoleHandler allows an admin to define a new role in the system.
func CreateRoleHandler(c *gin.Context) {
    // Define the expected request structure.
    type RoleRequest struct {
        Name        string `json:"name"`
        Description string `json:"description"`
    }

    var req RoleRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
        return
    }
    if req.Name == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Role name is required"})
        return
    }

    // Create a new Role record.
    role := database.Role{
        Name:        req.Name,
        Description: req.Description,
    }
    if err := database.DB.Create(&role).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create role", "details": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Role created successfully", "role": role})
}
