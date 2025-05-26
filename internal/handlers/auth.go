package handlers

import (
	"eventure/internal/authz"
	"eventure/internal/database"
	"eventure/internal/models"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type RegisterInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RefreshRequest struct {
	Username     string `json:"username" binding:"required"`
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// Simulated refresh token store (in-memory)
var (
	refreshTokenStore = make(map[string]string)
	storeMutex        = sync.Mutex{}
)

func RegisterHandler(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	input.Username = strings.TrimSpace(input.Username)
	if input.Username == "" || len(input.Password) < 6 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username and password (min 6 characters) are required"})
		return
	}

	var existing models.User
	if err := database.DB.Where("username = ?", input.Username).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not hash password"})
		return
	}

	var role models.Role
	if err := database.DB.Where("name = ?", "user").First(&role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Default role not found"})
		return
	}

	user := models.User{
		Username: input.Username,
		Password: string(hashedPassword),
		RoleID:   role.ID,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"role":    role.Name,
	})
}

func LoginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	var user models.User
	if err := database.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Simulate token strings (NOT real JWTs)
	accessToken := "access-token-" + strings.ToLower(user.Username)
	refreshToken := "refresh-token-" + strings.ToLower(user.Username)

	storeMutex.Lock()
	refreshTokenStore[user.Username] = refreshToken
	storeMutex.Unlock()

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func RefreshHandler(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	storeMutex.Lock()
	storedToken, ok := refreshTokenStore[req.Username]
	storeMutex.Unlock()

	if !ok || storedToken != req.RefreshToken {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	newAccessToken := "access-token-" + strings.ToLower(req.Username)
	c.JSON(http.StatusOK, gin.H{"access_token": newAccessToken})
}

func ProfileHandler(c *gin.Context) {
	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found in context"})
		return
	}

	u := user.(models.User)
	c.JSON(http.StatusOK, gin.H{
		"username": u.Username,
		"role":     u.Role,
	})
}

func CanViewEventsHandler(c *gin.Context) {
	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not in context"})
		return
	}

	u := user.(models.User)

	if authz.CanPerform(u, "read", "event") {
		c.JSON(http.StatusOK, gin.H{"message": "✅ Access granted via CanPerform()"})
	} else {
		c.JSON(http.StatusForbidden, gin.H{"message": "❌ Access denied via CanPerform()"})
	}
}
