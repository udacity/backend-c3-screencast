package middleware

import (
	"eventure/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequirePermission(resource string, action string) gin.HandlerFunc {
	required := resource + ":" + action

	return func(c *gin.Context) {
		userVal, exists := c.Get("currentUser")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
			c.Abort()
			return
		}

		user := userVal.(models.User)
		role := user.Role

		for _, p := range role.Permissions {
			if p.Name == required {
				c.Next()
				return
			}
		}

		logAccessDenied("Permission check failed", map[string]interface{}{
			"username":   user.Username,
			"user_id":    user.ID,
			"user_role":  role.Name,
			"required":   required,
			"path":       c.Request.URL.Path,
			"method":     c.Request.Method,
			"user_perms": role.Permissions,
		})

		c.JSON(http.StatusForbidden, gin.H{
			"error":      "Permission denied",
			"required":   required,
			"user_role":  role.Name,
			"user_perms": role.Permissions, // Helpful for debugging
		})
		c.Abort()
	}
}
