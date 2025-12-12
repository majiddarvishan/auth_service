package handlers

import (
	"net/http"
	"time"

	"auth_service/config"
	"auth_service/constants"
	"auth_service/database"
	"auth_service/logger"
	"auth_service/types"
	"auth_service/validation"

	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// RegisterRequest represents the payload for user registration. godoc
// swagger:model RegisterRequest
// @Description RegisterRequest defines the expected request body for creating a new user.
// @Property username body string true "Username for the new account"
// @Property password body string true "Password for the new account"
// @Property role     body string false "Role for the new user"
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// RegisterHandler handles new user registrations. godoc
// @Summary      Register a new user
// @Description  Create a new user account with username, password, and role
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request  body      RegisterRequest  true  "Registration payload"
// @Success      200      {object}  map[string]string  "User registered successfully"
// @Failure      400      {object}  map[string]string  "Invalid input or missing fields"
// @Failure      500      {object}  map[string]string  "Server error during registration"
// @Router       /users [post]
func RegisterHandler(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		status := http.StatusBadRequest
		requestID, _ := c.Get("request_id")
		c.JSON(status, types.APIError{
			Code:      constants.ErrorInvalidJSON,
			Message:   "Invalid JSON format",
			Timestamp: time.Now().Unix(),
			RequestID: requestID.(string),
		})
		return
	}

	if req.Username == "" || req.Password == "" {
		status := http.StatusBadRequest
		requestID, _ := c.Get("request_id")
		c.JSON(status, types.APIError{
			Code:      constants.ErrorInvalidUsername,
			Message:   "Username and password are required",
			Timestamp: time.Now().Unix(),
			RequestID: requestID.(string),
		})
		return
	}

	// Validate username format
	if err := validation.ValidateUsername(req.Username); err != nil {
		status := http.StatusBadRequest
		requestID, _ := c.Get("request_id")
		c.JSON(status, types.APIError{
			Code:      constants.ErrorInvalidUsername,
			Message:   err.Error(),
			Timestamp: time.Now().Unix(),
			RequestID: requestID.(string),
		})
		return
	}

	// Validate password strength
	if err := validation.ValidatePasswordStrength(req.Password); err != nil {
		status := http.StatusBadRequest
		requestID, _ := c.Get("request_id")
		c.JSON(status, types.APIError{
			Code:      constants.ErrorInvalidPassword2,
			Message:   err.Error(),
			Timestamp: time.Now().Unix(),
			RequestID: requestID.(string),
		})
		return
	}

	// Hash the password.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 14)
	if err != nil {
		status := http.StatusInternalServerError
		requestID, _ := c.Get("request_id")
		c.JSON(status, types.APIError{
			Code:      constants.ErrorInternalServer,
			Message:   "Could not hash password",
			Timestamp: time.Now().Unix(),
			RequestID: requestID.(string),
		})
		return
	}

	// Use provided role or assign a default role.
	roleName := req.Role
	if roleName == "" {
		roleName = constants.DefaultRole
	}

	// Look up the role in the database.
	role, err := database.DB.GetRoleByName(roleName)
	if err != nil {
		status := http.StatusBadRequest
		requestID, _ := c.Get("request_id")
		c.JSON(status, types.APIError{
			Code:      constants.ErrorRoleNotFound,
			Message:   "Role not found",
			Timestamp: time.Now().Unix(),
			RequestID: requestID.(string),
		})
		return
	}

	user := database.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Roles:    []database.Role{*role},
	}

	// // Use provided role or assign a default role.
	// role := req.Role
	// if role == "" {
	// 	role = "user"
	// }

	// user := database.User{
	// 	Username: req.Username,
	// 	Password: string(hashedPassword),
	// 	Role:     role,
	// }

	// Create user in the database.
	if err := database.DB.CreateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

// LoginRequest represents the payload for user login.
// swagger:model LoginRequest
// @Description LoginRequest defines the expected request body for logging in.
// @Property username body string true "Username of the account"
// @Property password body string true "Password of the account"
// @Property captchaId body string true  "ID of the captcha challenge"
// @Property captchaSolution body string true  "Solution to the captcha"
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginHandler authenticates the user and returns a JWT token with refresh token.
// @Summary      Login a user
// @Description  Authenticate user credentials and return a signed JWT with refresh token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request  body      LoginRequest  true  "Login payload"
// @Success      200      {object}  types.TokenResponse  "JWT token with refresh token"
// @Failure      400      {object}  types.APIError "Invalid JSON format"
// @Failure      401      {object}  types.APIError "Unauthorized: invalid credentials"
// @Failure      500      {object}  types.APIError "Server error during token generation"
// @Router       /login [post]
func LoginHandler(c *gin.Context) {
	requestID, _ := c.Get("request_id")
	log := logger.Get()

	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error("Invalid login request", "error", err.Error())
		c.JSON(http.StatusBadRequest, types.APIError{
			Code:      constants.ErrorInvalidJSON,
			Message:   "Invalid JSON format",
			Timestamp: time.Now().Unix(),
			RequestID: requestID.(string),
		})
		return
	}

	// Fetch user and its Role in one go:
	user, err := database.DB.GetUserAndRoleByUsername(req.Username)
	if err != nil {
		log.Error("User not found during login", "username", req.Username, "error", err.Error())
		c.JSON(http.StatusUnauthorized, types.APIError{
			Code:      constants.ErrorInvalidPassword,
			Message:   "Invalid username or password",
			Timestamp: time.Now().Unix(),
			RequestID: requestID.(string),
		})
		return
	}

	// Compare the stored hashed password with the incoming password.
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		log.Error("Invalid password during login", "username", req.Username)
		c.JSON(http.StatusUnauthorized, types.APIError{
			Code:      constants.ErrorInvalidPassword,
			Message:   "Invalid username or password",
			Timestamp: time.Now().Unix(),
			RequestID: requestID.(string),
		})
		return
	}

	accessToken, err := GenerateAccessToken(user)
	if err != nil {
		log.Error("Failed to generate access token", "user_id", user.ID, "error", err.Error())
		c.JSON(http.StatusInternalServerError, types.APIError{
			Code:      constants.ErrorInternalServer,
			Message:   "Failed to generate token",
			Timestamp: time.Now().Unix(),
			RequestID: requestID.(string),
		})
		return
	}

	// Generate refresh token
	refreshToken, err := StoreRefreshToken(user.ID)
	if err != nil {
		log.Error("Failed to generate refresh token", "user_id", user.ID, "error", err.Error())
		c.JSON(http.StatusInternalServerError, types.APIError{
			Code:      constants.ErrorInternalServer,
			Message:   "Failed to generate refresh token",
			Timestamp: time.Now().Unix(),
			RequestID: requestID.(string),
		})
		return
	}

	log.Info("User logged in successfully", "username", user.Username, "user_id", user.ID)

	c.JSON(http.StatusOK, types.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    GetTokenExpirationSeconds(),
		TokenType:    "Bearer",
		Timestamp:    time.Now().Unix(),
	})
}

// SecureL represents the payload for user login.
// swagger:model SecureL
// @Description SecureL defines the expected request body for logging in.
// @Property username body string true "Username of the account"
// @Property password body string true "Password of the account"
// @Property captchaId body string true  "ID of the captcha challenge"
// @Property captchaSolution body string true  "Solution to the captcha"
type SecureLoginRequest struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	CaptchaId       string `json:"captchaId"`
	CaptchaSolution string `json:"captchaSolution"`
}

// LoginHandler authenticates the user and returns a JWT token with refresh token.
// @Summary      Secure Login a user
// @Description  Authenticate user credentials with captcha and return a signed JWT with refresh token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request  body      SecureLoginRequest  true  "Login payload"
// @Success      200      {object}  types.TokenResponse  "JWT token with refresh token"
// @Failure      400      {object}  types.APIError "Invalid JSON format or captcha"
// @Failure      401      {object}  types.APIError "Unauthorized: invalid credentials"
// @Failure      500      {object}  types.APIError "Server error during token generation"
// @Router       /secure-login [post]
func SecureLoginHandler(c *gin.Context) {
	requestID, _ := c.Get("request_id")
	log := logger.Get()

	var req SecureLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error("Invalid secure login request", "error", err.Error())
		c.JSON(http.StatusBadRequest, types.APIError{
			Code:      constants.ErrorInvalidJSON,
			Message:   "Invalid JSON format",
			Timestamp: time.Now().Unix(),
			RequestID: requestID.(string),
		})
		return
	}

	if req.CaptchaId == "" || req.CaptchaSolution == "" {
		log.Error("Missing captcha fields in secure login")
		c.JSON(http.StatusBadRequest, types.APIError{
			Code:      constants.ErrorInvalidJSON,
			Message:   "Captcha is required",
			Timestamp: time.Now().Unix(),
			RequestID: requestID.(string),
		})
		return
	}

	if !captcha.VerifyString(req.CaptchaId, req.CaptchaSolution) {
		log.Error("Captcha verification failed", "username", req.Username)
		c.JSON(http.StatusUnauthorized, types.APIError{
			Code:      constants.ErrorInvalidPassword,
			Message:   "Captcha verification failed",
			Timestamp: time.Now().Unix(),
			RequestID: requestID.(string),
		})
		return
	}

	// Fetch user and its Role in one go:
	user, err := database.DB.GetUserAndRoleByUsername(req.Username)
	if err != nil {
		log.Error("User not found during secure login", "username", req.Username, "error", err.Error())
		c.JSON(http.StatusUnauthorized, types.APIError{
			Code:      constants.ErrorInvalidPassword,
			Message:   "Invalid username or password",
			Timestamp: time.Now().Unix(),
			RequestID: requestID.(string),
		})
		return
	}

	// Compare the stored hashed password with the incoming password.
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		log.Error("Invalid password during secure login", "username", req.Username)
		c.JSON(http.StatusUnauthorized, types.APIError{
			Code:      constants.ErrorInvalidPassword,
			Message:   "Invalid username or password",
			Timestamp: time.Now().Unix(),
			RequestID: requestID.(string),
		})
		return
	}

	accessToken, err := GenerateAccessToken(user)
	if err != nil {
		log.Error("Failed to generate access token", "user_id", user.ID, "error", err.Error())
		c.JSON(http.StatusInternalServerError, types.APIError{
			Code:      constants.ErrorInternalServer,
			Message:   "Failed to generate token",
			Timestamp: time.Now().Unix(),
			RequestID: requestID.(string),
		})
		return
	}

	// Generate refresh token
	refreshToken, err := StoreRefreshToken(user.ID)
	if err != nil {
		log.Error("Failed to generate refresh token", "user_id", user.ID, "error", err.Error())
		c.JSON(http.StatusInternalServerError, types.APIError{
			Code:      constants.ErrorInternalServer,
			Message:   "Failed to generate refresh token",
			Timestamp: time.Now().Unix(),
			RequestID: requestID.(string),
		})
		return
	}

	log.Info("User logged in securely", "username", user.Username, "user_id", user.ID)

	c.JSON(http.StatusOK, types.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    GetTokenExpirationSeconds(),
		TokenType:    "Bearer",
		Timestamp:    time.Now().Unix(),
	})
}

// DeleteUserHandler deletes a user based on the username passed in the URL parameter. godoc
// This endpoint should be accessible only to admins.
// @Summary      Delete a user
// @Description  Delete an existing user account (admin only)
// @Tags         Auth
// @Produce      json
// @Param        username  path      string  true  "Username to delete"
// @Success      200       {object}  map[string]string  "User deleted successfully"
// @Failure      400       {object}  map[string]string  "Username is required"
// @Failure      404       {object}  map[string]string  "User not found"
// @Failure      500       {object}  map[string]string  "Could not delete user"
// @Router       /users/{username} [delete]
func DeleteUserHandler(c *gin.Context) {
	// Get the username from the URL parameter.
	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username is required"})
		return
	}

	// Check if user exists first
	_, err := database.DB.GetUserByUsername(username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	err = database.DB.DeleteUserByUsername(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete user", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// RoleUpdateRequest represents the payload for updating a user's role. godoc
// swagger:model RoleUpdateRequest
// @Description RoleUpdateRequest defines the expected request body for role updates.
// @Property role body string true "New role for the user"
type RoleUpdateRequest struct {
	Role string `json:"role"`
}

// UpdateUserRoleHandler allows an admin to update a user's role. godoc
// @Summary      Update user role
// @Description  Update the role of an existing user (admin only)
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        username  path      string             true  "Username to update"
// @Param        request   body      RoleUpdateRequest  true  "Role update payload"
// @Success      200       {object}  map[string]string  "User role updated successfully"
// @Failure      400       {object}  map[string]string  "Invalid input or missing fields"
// @Failure      404       {object}  map[string]string  "User not found"
// @Failure      500       {object}  map[string]string  "Failed to update user role"
// @Router       /users/{username}/role [put]
func UpdateUserRoleHandler(c *gin.Context) {
	// Get the username from the URL parameter.
	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username is required"})
		return
	}

	var req RoleUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	if req.Role == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role is required"})
		return
	}

	err := database.DB.UpdateUserRoleByUsername(username, req.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user role", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User role updated successfully"})
}

// GenerateAccessToken creates a new JWT access token for a user
func GenerateAccessToken(user *database.User) (string, error) {
	var primaryRole string
	if len(user.Roles) > 0 {
		primaryRole = user.Roles[0].Name
	} else {
		primaryRole = constants.DefaultRole
	}

	expirationTime := time.Now().Add(config.TokenExpirationPeriod)
	claims := jwt.MapClaims{
		constants.ClaimKeyUserID: user.ID,
		constants.ClaimKeyRole:   primaryRole,
		constants.ClaimKeyExpiry: expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.SecretKey))
}

// GetTokenExpirationSeconds returns token expiration in seconds
func GetTokenExpirationSeconds() int {
	return int(config.TokenExpirationPeriod.Seconds())
}
