package controller

import (
	"bytes"
	"encoding/json"
	"mess-registration/internal/schema"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	database, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	database.AutoMigrate(&schema.User{}, &schema.SwapRequest{})
	return database
}

func setupTestRouter(database *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	userController := NewUserController(database)

	r.POST("/api/register", userController.RegisterUser)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return r
}

func TestHealthCheck(t *testing.T) {
	database := setupTestDB()
	router := setupTestRouter(database)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "ok", response["status"])
}

func TestRegisterUser(t *testing.T) {
	database := setupTestDB()
	router := setupTestRouter(database)

	// Test successful registration
	t.Run("successful registration", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"userID": "test_user_123",
		}
		jsonBody, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/register?mess=0", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Registration successful", response["message"])
	})

	// Test missing mess parameter
	t.Run("missing mess parameter", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"userID": "test_user_456",
		}
		jsonBody, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/register", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
	})

	// Test invalid mess number
	t.Run("invalid mess number", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"userID": "test_user_789",
		}
		jsonBody, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/register?mess=5", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
	})
}

func TestMain(m *testing.M) {
	// Setup
	os.Setenv("JWT_SECRET", "test-secret")
	os.Setenv("GIN_MODE", "test")

	// Run tests
	code := m.Run()

	// Cleanup
	os.Exit(code)
}
