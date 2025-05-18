package handlers

import (
	"net/http"
	"time"

	"auth_service/config"
	"auth_service/database"
	// "auth_service/pkg/trie"

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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	if req.Username == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username and password are required"})
		return
	}

	// Hash the password.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 14)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not hash password"})
		return
	}

	// Use provided role or assign a default role.
	roleName := req.Role
	if roleName == "" {
		roleName = "guest"
	}

	// Look up the role in the database.
	role, err := database.DB.GetRoleByName(roleName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role not found"})
		return
	}

	user := database.User{
		Username: req.Username,
		Password: string(hashedPassword),
		RoleID:   role.ID,
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
	Username        string `json:"username"`
	Password        string `json:"password"`
	CaptchaId       string `json:"captchaId"`
	CaptchaSolution string `json:"captchaSolution"`
}

// LoginHandler authenticates the user and returns a JWT token.
// @Summary      Login a user
// @Description  Authenticate user credentials and return a signed JWT
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request  body      LoginRequest  true  "Login payload"
// @Success      200      {object}  map[string]string  "JWT token"
// @Failure      400      {object}  map[string]string  "Invalid JSON format"
// @Failure      401      {object}  map[string]string  "Unauthorized: invalid credentials"
// @Failure      500      {object}  map[string]string  "Server error during token generation"
// @Router       /login [post]
func LoginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	if req.CaptchaId == "" || req.CaptchaSolution == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Captcha is required"})
		return
	}

	if !captcha.VerifyString(req.CaptchaId, req.CaptchaSolution) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Captcha verification failed"})
		return
	}

	// Fetch user and its Role in one go:
	user, err := database.DB.GetUserAndRoleByUsername(req.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Compare the stored hashed password with the incoming password.
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	// Create JWT claims: subject, role, and expiry.
	expirationTime := time.Now().Add(config.TokenExpirationPeriod)
	claims := jwt.MapClaims{
		"user": user.ID,
		"role": user.Role.Name,
		"exp":  expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(config.SecretKey))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenStr})
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
// @Router       /user/{username} [delete]
func DeleteUserHandler(c *gin.Context) {
	// Get the username from the URL parameter.
	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username is required"})
		return
	}

	err := database.DB.DeleteUserByUsername(username)
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
// @Router       /user/{username}/role [put]
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

// @Summary      Get user phone numbers
// @Description  Retrieves all phone numbers associated with the authenticated user.
// @Tags         user
// @Security     ApiKeyAuth
// @Produce      json
// @Success      200  {object}  map[string][]string  "{"phones": ["+1234567890", "+10987654321"]}"
// @Failure      401  {object}  gin.H                "{"error": "user not authenticated"}"
// @Failure      500  {object}  gin.H                "{"error": "could not fetch phones"}"
// @Router       /user/{username}/phones [get]
// func GetUserPhonesHandler(c *gin.Context) {
// 	userName := c.Param("username")
// 	if userName == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Username is required"})
// 		return
// 	}

// 	nums, err := database.DB.GetUserPhones(userName)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"phones": nums})
// }

// @Summary      Add multiple phone numbers for the user
// @Description  Adds one or more new phone numbers associated with the authenticated user.
// @Tags         user
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        payload  body      struct{Numbers []string ` + "`json:\"numbers\" binding:\"required,min=1,dive,required\"`" + `}  true  "List of phone numbers"
// @Success      201      {object}  gin.H "{\"message\": \"N phones added\"}"
// @Failure      400      {object}  gin.H "{\"error\": \"invalid payload\"}"
// @Failure      401      {object}  gin.H "{\"error\": \"user not authenticated\"}"
// @Failure      500      {object}  gin.
// @Router       /user/{username}/phones [post]
// func AddUserPhonesHandler(c *gin.Context) {
// 	type addPhoneRequest struct {
// 		Numbers []string `json:"numbers" binding:"required,min=1,dive,required"`
// 	}

// 	userName := c.Param("username")
// 	if userName == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Username is required"})
// 		return
// 	}

// 	// Bind and validate input
// 	var req addPhoneRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
// 		return
// 	}

// 	if err := database.DB.AddPhoneForUser(userName, req.Numbers); err != nil {
// 		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	for _, num := range req.Numbers {
// 		trie.TrieManagerInstance.Add(userName, num)
// 	}

// 	c.JSON(http.StatusCreated, gin.H{"message": "phone added"})
// }
