package handlers

import (
	"net/http"
	"time"

	"auth_service/constants"
	"auth_service/database"
	"auth_service/logger"
	"auth_service/types"

	"github.com/gin-gonic/gin"
)

// UserRolesRequest represents the payload for managing user roles
type UserRolesRequest struct {
	Roles []string `json:"roles" binding:"required,min=1"`
}

// AddUserRolesRequest represents adding roles to a user
type AddUserRolesRequest struct {
	Roles []string `json:"roles" binding:"required,min=1"`
}

// UserRolesResponse represents user roles response
type UserRolesResponse struct {
	UserID uint   `json:"user_id"`
	Roles  []Role `json:"roles"`
}

// Role represents a role with its details
type Role struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// SetUserRoles sets all roles for a user (replaces existing ones)
// @Summary Set user roles
// @Description Replace all roles for a user with the provided ones
// @Tags Users
// @Accept json
// @Produce json
// @Param username path string true "Username"
// @Param request body UserRolesRequest true "Roles to set"
// @Success 200 {object} UserRolesResponse
// @Failure 400 {object} types.APIError "Invalid request"
// @Failure 401 {object} types.APIError "Unauthorized"
// @Failure 404 {object} types.APIError "User not found"
// @Failure 500 {object} types.APIError "Server error"
// @Router /admin/users/{username}/roles [put]
func SetUserRoles(c *gin.Context) {
	requestID, _ := c.Get("request_id")
	log := logger.Get()

	username := c.Param("username")
	if username == "" {
		log.Error("Missing username parameter", "error", "username is required")
		c.JSON(http.StatusBadRequest, types.APIError{
			Code:      constants.ErrorInvalidJSON,
			Message:   "Username is required",
			Timestamp: time.Now().Unix(),
			RequestID: requestID.(string),
		})
		return
	}

	var req UserRolesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error("Invalid request format", "error", err.Error())
		c.JSON(http.StatusBadRequest, types.APIError{
			Code:      constants.ErrorInvalidJSON,
			Message:   "Invalid JSON format",
			Timestamp: time.Now().Unix(),
			RequestID: requestID.(string),
		})
		return
	}

	// Verify user exists
	user, err := database.DB.GetUserByUsername(username)
	if err != nil {
		log.Error("User not found", "username", username, "error", err.Error())
		c.JSON(http.StatusNotFound, types.APIError{
			Code:      constants.ErrorUserNotFound,
			Message:   "User not found",
			Timestamp: time.Now().Unix(),
			RequestID: requestID.(string),
		})
		return
	}

	// Set roles
	if err := database.DB.SetUserRoles(user.ID, req.Roles); err != nil {
		log.Error("Failed to set user roles", "user_id", user.ID, "error", err.Error())
		c.JSON(http.StatusInternalServerError, types.APIError{
			Code:      constants.ErrorInternalServer,
			Message:   "Failed to set user roles",
			Timestamp: time.Now().Unix(),
			RequestID: requestID.(string),
		})
		return
	}

	// Get updated roles
	roles, err := database.DB.GetUserRoles(user.ID)
	if err != nil {
		log.Error("Failed to get user roles", "user_id", user.ID, "error", err.Error())
		c.JSON(http.StatusInternalServerError, types.APIError{
			Code:      constants.ErrorInternalServer,
			Message:   "Failed to get user roles",
			Timestamp: time.Now().Unix(),
			RequestID: requestID.(string),
		})
		return
	}

	// Convert roles to response format
	roleResponses := make([]Role, len(roles))
	for i, r := range roles {
		roleResponses[i] = Role{
			ID:          r.ID,
			Name:        r.Name,
			Description: r.Description,
		}
	}

	log.Info("User roles updated", "user_id", user.ID, "username", user.Username, "roles", req.Roles)

	c.JSON(http.StatusOK, types.APISuccess{
		Data: UserRolesResponse{
			UserID: user.ID,
			Roles:  roleResponses,
		},
		Message:   "User roles updated successfully",
		Timestamp: time.Now().Unix(),
		RequestID: requestID.(string),
	})
}

// AddUserRoles adds roles to a user (doesn't replace existing ones)
// @Summary Add roles to user
// @Description Add one or more roles to a user while keeping existing roles
// @Tags Users
// @Accept json
// @Produce json
// @Param username path string true "Username"
// @Param request body AddUserRolesRequest true "Roles to add"
// @Success 200 {object} UserRolesResponse
// @Failure 400 {object} types.APIError "Invalid request"
// @Failure 401 {object} types.APIError "Unauthorized"
// @Failure 404 {object} types.APIError "User not found"
// @Failure 500 {object} types.APIError "Server error"
// @Router /admin/users/{username}/roles [post]
func AddUserRoles(c *gin.Context) {
	requestID, _ := c.Get("request_id")
	log := logger.Get()

	username := c.Param("username")
	if username == "" {
		log.Error("Missing username parameter", "error", "username is required")
		c.JSON(http.StatusBadRequest, types.APIError{
			Code:      constants.ErrorInvalidJSON,
			Message:   "Username is required",
			Timestamp: time.Now().Unix(),
			RequestID: requestID.(string),
		})
		return
	}

	var req AddUserRolesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error("Invalid request format", "error", err.Error())
		c.JSON(http.StatusBadRequest, types.APIError{
			Code:      constants.ErrorInvalidJSON,
			Message:   "Invalid JSON format",
			Timestamp: time.Now().Unix(),
			RequestID: requestID.(string),
		})
		return
	}

	// Verify user exists
	user, err := database.DB.GetUserByUsername(username)
	if err != nil {
		log.Error("User not found", "username", username, "error", err.Error())
		c.JSON(http.StatusNotFound, types.APIError{
			Code:      constants.ErrorUserNotFound,
			Message:   "User not found",
			Timestamp: time.Now().Unix(),
			RequestID: requestID.(string),
		})
		return
	}

	// Add roles
	if err := database.DB.AddRolesToUser(user.ID, req.Roles); err != nil {
		log.Error("Failed to add roles to user", "user_id", user.ID, "error", err.Error())
		c.JSON(http.StatusInternalServerError, types.APIError{
			Code:      constants.ErrorInternalServer,
			Message:   "Failed to add roles: " + err.Error(),
			Timestamp: time.Now().Unix(),
			RequestID: requestID.(string),
		})
		return
	}

	// Get updated roles
	roles, err := database.DB.GetUserRoles(user.ID)
	if err != nil {
		log.Error("Failed to get user roles", "user_id", user.ID, "error", err.Error())
		c.JSON(http.StatusInternalServerError, types.APIError{
			Code:      constants.ErrorInternalServer,
			Message:   "Failed to get user roles",
			Timestamp: time.Now().Unix(),
			RequestID: requestID.(string),
		})
		return
	}

	// Convert roles to response format
	roleResponses := make([]Role, len(roles))
	for i, r := range roles {
		roleResponses[i] = Role{
			ID:          r.ID,
			Name:        r.Name,
			Description: r.Description,
		}
	}

	log.Info("Roles added to user", "user_id", user.ID, "username", user.Username, "new_roles", req.Roles)

	c.JSON(http.StatusOK, types.APISuccess{
		Data: UserRolesResponse{
			UserID: user.ID,
			Roles:  roleResponses,
		},
		Message:   "Roles added successfully",
		Timestamp: time.Now().Unix(),
		RequestID: requestID.(string),
	})
}

// RemoveUserRoles removes roles from a user
// @Summary Remove roles from user
// @Description Remove one or more roles from a user
// @Tags Users
// @Accept json
// @Produce json
// @Param username path string true "Username"
// @Param request body AddUserRolesRequest true "Roles to remove"
// @Success 200 {object} UserRolesResponse
// @Failure 400 {object} types.APIError "Invalid request"
// @Failure 401 {object} types.APIError "Unauthorized"
// @Failure 404 {object} types.APIError "User not found"
// @Failure 500 {object} types.APIError "Server error"
// @Router /admin/users/{username}/roles [delete]
func RemoveUserRoles(c *gin.Context) {
	requestID, _ := c.Get("request_id")
	log := logger.Get()

	username := c.Param("username")
	if username == "" {
		log.Error("Missing username parameter", "error", "username is required")
		c.JSON(http.StatusBadRequest, types.APIError{
			Code:      constants.ErrorInvalidJSON,
			Message:   "Username is required",
			Timestamp: time.Now().Unix(),
			RequestID: requestID.(string),
		})
		return
	}

	var req AddUserRolesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error("Invalid request format", "error", err.Error())
		c.JSON(http.StatusBadRequest, types.APIError{
			Code:      constants.ErrorInvalidJSON,
			Message:   "Invalid JSON format",
			Timestamp: time.Now().Unix(),
			RequestID: requestID.(string),
		})
		return
	}

	// Verify user exists
	user, err := database.DB.GetUserByUsername(username)
	if err != nil {
		log.Error("User not found", "username", username, "error", err.Error())
		c.JSON(http.StatusNotFound, types.APIError{
			Code:      constants.ErrorUserNotFound,
			Message:   "User not found",
			Timestamp: time.Now().Unix(),
			RequestID: requestID.(string),
		})
		return
	}

	// Remove roles
	if err := database.DB.RemoveRolesFromUser(user.ID, req.Roles); err != nil {
		log.Error("Failed to remove roles from user", "user_id", user.ID, "error", err.Error())
		c.JSON(http.StatusInternalServerError, types.APIError{
			Code:      constants.ErrorInternalServer,
			Message:   "Failed to remove roles: " + err.Error(),
			Timestamp: time.Now().Unix(),
			RequestID: requestID.(string),
		})
		return
	}

	// Get updated roles
	roles, err := database.DB.GetUserRoles(user.ID)
	if err != nil {
		log.Error("Failed to get user roles", "user_id", user.ID, "error", err.Error())
		c.JSON(http.StatusInternalServerError, types.APIError{
			Code:      constants.ErrorInternalServer,
			Message:   "Failed to get user roles",
			Timestamp: time.Now().Unix(),
			RequestID: requestID.(string),
		})
		return
	}

	// Convert roles to response format
	roleResponses := make([]Role, len(roles))
	for i, r := range roles {
		roleResponses[i] = Role{
			ID:          r.ID,
			Name:        r.Name,
			Description: r.Description,
		}
	}

	log.Info("Roles removed from user", "user_id", user.ID, "username", user.Username, "removed_roles", req.Roles)

	c.JSON(http.StatusOK, types.APISuccess{
		Data: UserRolesResponse{
			UserID: user.ID,
			Roles:  roleResponses,
		},
		Message:   "Roles removed successfully",
		Timestamp: time.Now().Unix(),
		RequestID: requestID.(string),
	})
}

// GetUserRoles gets all roles for a user
// @Summary Get user roles
// @Description Get all roles assigned to a user
// @Tags Users
// @Produce json
// @Param username path string true "Username"
// @Success 200 {object} UserRolesResponse
// @Failure 401 {object} types.APIError "Unauthorized"
// @Failure 404 {object} types.APIError "User not found"
// @Failure 500 {object} types.APIError "Server error"
// @Router /admin/users/{username}/roles [get]
func GetUserRoles(c *gin.Context) {
	requestID, _ := c.Get("request_id")
	log := logger.Get()

	username := c.Param("username")
	if username == "" {
		log.Error("Missing username parameter", "error", "username is required")
		c.JSON(http.StatusBadRequest, types.APIError{
			Code:      constants.ErrorInvalidJSON,
			Message:   "Username is required",
			Timestamp: time.Now().Unix(),
			RequestID: requestID.(string),
		})
		return
	}

	// Verify user exists
	user, err := database.DB.GetUserByUsername(username)
	if err != nil {
		log.Error("User not found", "username", username, "error", err.Error())
		c.JSON(http.StatusNotFound, types.APIError{
			Code:      constants.ErrorUserNotFound,
			Message:   "User not found",
			Timestamp: time.Now().Unix(),
			RequestID: requestID.(string),
		})
		return
	}

	// Get roles
	roles, err := database.DB.GetUserRoles(user.ID)
	if err != nil {
		log.Error("Failed to get user roles", "user_id", user.ID, "error", err.Error())
		c.JSON(http.StatusInternalServerError, types.APIError{
			Code:      constants.ErrorInternalServer,
			Message:   "Failed to get user roles",
			Timestamp: time.Now().Unix(),
			RequestID: requestID.(string),
		})
		return
	}

	// Convert roles to response format
	roleResponses := make([]Role, len(roles))
	for i, r := range roles {
		roleResponses[i] = Role{
			ID:          r.ID,
			Name:        r.Name,
			Description: r.Description,
		}
	}

	log.Info("Retrieved user roles", "user_id", user.ID, "username", user.Username)

	c.JSON(http.StatusOK, types.APISuccess{
		Data: UserRolesResponse{
			UserID: user.ID,
			Roles:  roleResponses,
		},
		Message:   "User roles retrieved successfully",
		Timestamp: time.Now().Unix(),
		RequestID: requestID.(string),
	})
}
