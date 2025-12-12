package handlers

import (
	"auth_service/constants"
	"auth_service/database"
	"auth_service/logger"
	"auth_service/types"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// RefreshTokenRequest represents the refresh token request payload
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RefreshAccessTokenHandler generates a new access token from a valid refresh token
// @Summary Refresh access token
// @Description Exchange a valid refresh token for a new access token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body RefreshTokenRequest true "Refresh token"
// @Success 200 {object} types.TokenResponse "New access token generated"
// @Failure 400 {object} types.APIError "Invalid request"
// @Failure 401 {object} types.APIError "Invalid or expired refresh token"
// @Failure 500 {object} types.APIError "Server error"
// @Router /refresh [post]
func RefreshAccessTokenHandler(c *gin.Context) {
	requestID, _ := c.Get("request_id")
	log := logger.Get()

	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error("Invalid refresh token request", "error", err.Error())
		c.JSON(http.StatusBadRequest, types.APIError{
			Code:      constants.ErrorInvalidJSON,
			Message:   "Invalid request format",
			Timestamp: time.Now().Unix(),
			RequestID: requestID.(string),
		})
		return
	}

	// Validate refresh token
	rt, err := database.DB.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		log.Error("Invalid refresh token", "error", err.Error())
		c.JSON(http.StatusUnauthorized, types.APIError{
			Code:      constants.ErrorInvalidRefreshToken,
			Message:   "Invalid or expired refresh token",
			Timestamp: time.Now().Unix(),
			RequestID: requestID.(string),
		})
		return
	}

	// Get user details
	user, err := database.DB.GetUserByID(rt.UserID)
	if err != nil {
		log.Error("User not found for refresh token", "user_id", rt.UserID, "error", err.Error())
		c.JSON(http.StatusNotFound, types.APIError{
			Code:      constants.ErrorUserNotFound,
			Message:   "User not found",
			Timestamp: time.Now().Unix(),
			RequestID: requestID.(string),
		})
		return
	}

	// Generate new access token
	accessToken, err := GenerateAccessToken(user)
	if err != nil {
		log.Error("Failed to generate access token", "user_id", user.ID, "error", err.Error())
		c.JSON(http.StatusInternalServerError, types.APIError{
			Code:      constants.ErrorInternalServer,
			Message:   "Failed to generate access token",
			Timestamp: time.Now().Unix(),
			RequestID: requestID.(string),
		})
		return
	}

	log.Info("Access token refreshed", "user_id", user.ID, "username", user.Username)

	c.JSON(http.StatusOK, types.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: req.RefreshToken, // Can be the same or a new one
		ExpiresIn:    GetTokenExpirationSeconds(),
		TokenType:    "Bearer",
		Timestamp:    time.Now().Unix(),
	})
}

// RevokeRefreshTokenHandler revokes a refresh token (logout)
// @Summary Revoke refresh token
// @Description Logout by revoking the refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body RefreshTokenRequest true "Refresh token to revoke"
// @Success 200 {object} types.APISuccess "Token revoked"
// @Failure 400 {object} types.APIError "Invalid request"
// @Failure 404 {object} types.APIError "Token not found"
// @Failure 500 {object} types.APIError "Server error"
// @Router /logout [post]
func RevokeRefreshTokenHandler(c *gin.Context) {
	requestID, _ := c.Get("request_id")
	log := logger.Get()

	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error("Invalid revoke request", "error", err.Error())
		c.JSON(http.StatusBadRequest, types.APIError{
			Code:      constants.ErrorInvalidJSON,
			Message:   "Invalid request format",
			Timestamp: time.Now().Unix(),
			RequestID: requestID.(string),
		})
		return
	}

	// Get token to find user ID
	rt, err := database.DB.GetRefreshToken(req.RefreshToken)
	if err != nil {
		log.Error("Refresh token not found", "error", err.Error())
		c.JSON(http.StatusNotFound, types.APIError{
			Code:      constants.ErrorInvalidRefreshToken,
			Message:   "Refresh token not found",
			Timestamp: time.Now().Unix(),
			RequestID: requestID.(string),
		})
		return
	}

	// Revoke the token
	if err := database.DB.RevokeRefreshToken(req.RefreshToken); err != nil {
		log.Error("Failed to revoke refresh token", "error", err.Error())
		c.JSON(http.StatusInternalServerError, types.APIError{
			Code:      constants.ErrorInternalServer,
			Message:   "Failed to revoke token",
			Timestamp: time.Now().Unix(),
			RequestID: requestID.(string),
		})
		return
	}

	log.Info("Refresh token revoked", "user_id", rt.UserID)

	c.JSON(http.StatusOK, types.APISuccess{
		Message:   "Token revoked successfully",
		Timestamp: time.Now().Unix(),
		RequestID: requestID.(string),
	})
}

// LogoutAllDevicesHandler revokes all refresh tokens for a user (logout from all devices)
// @Summary Logout from all devices
// @Description Revoke all refresh tokens for the current user
// @Tags Auth
// @Security Bearer
// @Produce json
// @Success 200 {object} types.APISuccess "All tokens revoked"
// @Failure 401 {object} types.APIError "Unauthorized"
// @Failure 500 {object} types.APIError "Server error"
// @Router /logout-all [post]
func LogoutAllDevicesHandler(c *gin.Context) {
	requestID, _ := c.Get("request_id")
	log := logger.Get()

	// Get user ID from JWT claims
	claims := c.MustGet("claims").(jwt.MapClaims)
	userID := uint(claims[constants.ClaimKeyUserID].(float64))

	// Revoke all tokens for this user
	if err := database.DB.RevokeAllUserTokens(userID); err != nil {
		log.Error("Failed to revoke all user tokens", "user_id", userID, "error", err.Error())
		c.JSON(http.StatusInternalServerError, types.APIError{
			Code:      constants.ErrorInternalServer,
			Message:   "Failed to logout from all devices",
			Timestamp: time.Now().Unix(),
			RequestID: requestID.(string),
		})
		return
	}

	log.Info("Logged out from all devices", "user_id", userID)

	c.JSON(http.StatusOK, types.APISuccess{
		Message:   "Logged out from all devices successfully",
		Timestamp: time.Now().Unix(),
		RequestID: requestID.(string),
	})
}

// GenerateRefreshToken creates a new refresh token
func GenerateRefreshToken() (string, error) {
	// Generate a cryptographically secure random token
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// StoreRefreshToken saves a refresh token to the database
func StoreRefreshToken(userID uint) (string, error) {
	token, err := GenerateRefreshToken()
	if err != nil {
		return "", err
	}

	// Get token expiration time (typically longer than access token)
	refreshExpiration := time.Now().Add(7 * 24 * time.Hour).Unix() // 7 days

	rt := &database.RefreshToken{
		UserID:    userID,
		Token:     token,
		ExpiresAt: refreshExpiration,
		Revoked:   false,
	}

	if err := database.DB.CreateRefreshToken(rt); err != nil {
		return "", err
	}

	return token, nil
}
