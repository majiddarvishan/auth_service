package handlers

import (
	"net/http"

	"auth_service/database"
	"github.com/gin-gonic/gin"
)

// SwaggerRoleRequest represents the payload to create a new role.
// swagger:model RoleRequest
type SwaggerRoleRequest struct {
	// The name of the role
	// required: true
	Name string `json:"name"`

	// A short description of the role
	Description string `json:"description"`
}

// SwaggerRole represents a role in the system.
// swagger:model Role
type SwaggerRole struct {
	// The unique ID of the role
	ID uint `json:"id"`

	// The name of the role
	Name string `json:"name"`

	// The description of the role
	Description string `json:"description"`
}

// SuccessResponse represents a successful response with a single role.
// swagger:model SuccessResponse
type SuccessResponse struct {
	Message string `json:"message"`
	Role    SwaggerRole   `json:"role"`
}

// RolesListResponse represents a list of roles.
// swagger:model RolesListResponse
type RolesListResponse struct {
	Roles []SwaggerRole `json:"roles"`
}

// ErrorResponse represents a standard error structure.
// swagger:model ErrorResponse
type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}

// CreateRoleHandler allows an admin to define a new role in the system.
//
// @Summary      Create a new role
// @Description  Allows an admin to create a new role by specifying name and description.
// @Tags         roles
// @Accept       json
// @Produce      json
// @Param        role  body      SwaggerRoleRequest  true  "Role details"
// @Success      200   {object}  SuccessResponse
// @Failure      400   {object}  ErrorResponse
// @Failure      500   {object}  ErrorResponse
// @Router       /roles [post]
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

// GetRolesHandler returns all roles.
//
// @Summary      List all roles
// @Description  Returns a list of all roles defined in the system.
// @Tags         roles
// @Produce      json
// @Success      200  {object}  RolesListResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /roles [get]
func GetRolesHandler(c *gin.Context) {
    var roles []database.Role
    if err := database.DB.Find(&roles).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch roles"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"roles": roles})
}
