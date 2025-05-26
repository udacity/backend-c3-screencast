package middleware

import (
	"net/http"

	"eventure/internal/models"

	"github.com/gin-gonic/gin"
)

func RequireRole(roleName string) gin.HandlerFunc {

	return func(c *gin.Context) {

		user, exists := c.Get("currentUser")

		if !exists {

			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})

			c.Abort()

			return

		}

		u := user.(models.User)

		if u.Role.Name != roleName {

			c.JSON(http.StatusForbidden, gin.H{

				"error": "Insufficient permissions",

				"your_role": u.Role.Name,

				"required": roleName,
			})

			c.Abort()

			return

		}

		c.Next()

	}

}
