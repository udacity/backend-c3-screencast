package middleware

import (
	"eventure/internal/database"
	"eventure/internal/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if !strings.HasPrefix(auth, "Bearer access-token-") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid token"})
			c.Abort()
			return
		}

		username := strings.TrimPrefix(auth, "Bearer access-token-")
		var user models.User
		if err := database.DB.
			Preload("Role").
			Preload("Role.Permissions").
			Where("username = ?", username).
			First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token or user not found"})
			c.Abort()
			return
		}

		// Set user in context
		c.Set("currentUser", user)
		c.Next()
	}
}
