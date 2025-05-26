package middleware

import (
	"eventure/internal/models"
	"log"
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

		log.Printf("[ACCESS DENIED] User '%s' (ID %d) with role '%s' lacks permission '%s'",
			user.Username, user.ID, role.Name, required)

		c.JSON(http.StatusForbidden, gin.H{
			"error":      "Permission denied",
			"required":   required,
			"user_role":  role.Name,
			"user_perms": role.Permissions, // Helpful for debugging
		})
		c.Abort()
	}
}
