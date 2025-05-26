package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"log"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"eventure/internal/database"
	"eventure/internal/handlers"
	"eventure/internal/middleware"

	"github.com/joho/godotenv"
)

func setupTestServer() *gin.Engine {
	gin.SetMode(gin.TestMode)

	// Load environment variables from .env (for DB_DSN)
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Failed to load .env file for tests:", err)
	}

	err = database.Connect()
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	r := gin.Default()
	r.Use(middleware.SecureHeaders())
	r.Use(middleware.CustomRecovery())
	handlers.RegisterRoutes(r)
	return r
}

func TestSecurityFlow(t *testing.T) {
	r := setupTestServer()

	// 1. Register user
	registerBody := `{"username":"testuser","password":"password123"}`
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer([]byte(registerBody)))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusCreated, resp.Code)

	// 2. Login user
	loginBody := `{"username":"testuser","password":"password123"}`
	req = httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer([]byte(loginBody)))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	var loginResp map[string]string
	json.Unmarshal(resp.Body.Bytes(), &loginResp)
	accessToken := loginResp["access_token"]
	refreshToken := loginResp["refresh_token"]
	assert.NotEmpty(t, accessToken)
	assert.NotEmpty(t, refreshToken)

	// 3. Access profile (with token)
	req = httptest.NewRequest(http.MethodGet, "/auth/profile", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	// 4. Access profile (without token)
	req = httptest.NewRequest(http.MethodGet, "/auth/profile", nil)
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusUnauthorized, resp.Code)

	// 5. Refresh token
	refreshPayload := map[string]string{
		"username":      "testuser",
		"refresh_token": refreshToken,
	}
	bodyBytes, _ := json.Marshal(refreshPayload)
	req = httptest.NewRequest(http.MethodPost, "/auth/refresh", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	var refreshResp map[string]string
	json.Unmarshal(resp.Body.Bytes(), &refreshResp)
	assert.NotEmpty(t, refreshResp["access_token"])
}
