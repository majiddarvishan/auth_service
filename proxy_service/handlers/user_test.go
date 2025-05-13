package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"yourapp/handlers"
	"yourapp/models"
)

// setupRouter initializes Gin with the GET and POST /user/phones routes for testing.
func setupRouter(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// Mock auth middleware: sets userID=1
	r.Use(func(c *gin.Context) {
		c.Set("userID", uint(1))
		c.Next()
	})

	protected := r.Group("/user")
	protected.GET("/phones", handlers.GetUserPhones(db))
	protected.POST("/phones", handlers.AddUserPhones(db))

	return r
}

// initializeTestDB creates an in-memory SQLite DB and migrates the Phone model.
func initializeTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to in-memory db: %v", err)
	}
	// Migrate Phone model
	err = db.AutoMigrate(&models.Phone{})
	if err != nil {
		t.Fatalf("failed to migrate Phone model: %v", err)
	}
	return db
}

func TestGetUserPhones_Empty(t *testing.T) {
	db := initializeTestDB(t)
	r := setupRouter(db)

	req, _ := http.NewRequest(http.MethodGet, "/user/phones", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var body map[string][]string
	err := json.Unmarshal(rec.Body.Bytes(), &body)
	assert.NoError(t, err)
	assert.Empty(t, body["phones"], "Expected no phones for new user")
}

func TestAddUserPhones_ThenGet(t *testing.T) {
	db := initializeTestDB(t)
	r := setupRouter(db)

	// Prepare payload
	numbers := []string{"+1111111111", "+2222222222"}
	payload := map[string][]string{"numbers": numbers}
	b, _ := json.Marshal(payload)

	// POST /user/phones
	reqPost, _ := http.NewRequest(http.MethodPost, "/user/phones", bytes.NewBuffer(b))
	reqPost.Header.Set("Content-Type", "application/json")
	recPost := httptest.NewRecorder()
	r.ServeHTTP(recPost, reqPost)

	assert.Equal(t, http.StatusCreated, recPost.Code)

	// GET /user/phones
	reqGet, _ := http.NewRequest(http.MethodGet, "/user/phones", nil)
	recGet := httptest.NewRecorder()
	r.ServeHTTP(recGet, reqGet)

	assert.Equal(t, http.StatusOK, recGet.Code)

	var body map[string][]string
	err := json.Unmarshal(recGet.Body.Bytes(), &body)
	assert.NoError(t, err)
	assert.ElementsMatch(t, numbers, body["phones"], "Expected returned phones to match posted ones")
}

func TestAddUserPhones_InvalidPayload(t *testing.T) {
	db := initializeTestDB(t)
	r := setupRouter(db)

	// Missing 'numbers'
	req, _ := http.NewRequest(http.MethodPost, "/user/phones", bytes.NewBuffer([]byte(`{"nums": []}`)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
