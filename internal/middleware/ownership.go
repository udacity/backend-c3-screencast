package middleware

import (
	"eventure/internal/database"
	"eventure/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func RequireEventOwnershipOrPermission(action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRaw, exists := c.Get("currentUser")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
			c.Abort()
			return
		}
		user := userRaw.(models.User)

		// Get event ID from URL
		idParam := c.Param("id")
		eventID, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
			c.Abort()
			return
		}

		// Load event
		var event models.Event
		if err := database.DB.First(&event, eventID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
			c.Abort()
			return
		}

		// Allow if user is owner
		if event.OwnerID == user.ID {
			c.Next()
			return
		}

		// Check if user has override permission
		required := action + ":any"
		for _, perm := range user.Role.Permissions {
			if perm.Name == "event:"+required {
				c.Next()
				return
			}
		}

		logAccessDenied("Ownership and override check failed", map[string]interface{}{
			"event_id": eventID,
			"owner_id": event.OwnerID,
			"user_id":  user.ID,
			"user":     user.Username,
			"role":     user.Role.Name,
			"required": "event:" + required,
			"path":     c.Request.URL.Path,
		})

		c.JSON(http.StatusForbidden, gin.H{
			"error":       "Not authorized for this event",
			"owner_id":    event.OwnerID,
			"user_id":     user.ID,
			"required":    "event:" + required,
			"user_role":   user.Role.Name,
			"user_perms":  user.Role.Permissions,
			"event_title": event.Title,
		})
		c.Abort()
	}
}
