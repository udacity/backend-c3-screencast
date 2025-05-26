package middleware

import (
	"eventure/internal/database"
	"eventure/internal/models"
	"log"
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

		log.Printf("[ACCESS DENIED] User '%s' (ID %d, role '%s') tried to %s event %d but is not the owner and lacks permission '%s'",
			user.Username, user.ID, user.Role.Name, action, event.ID, "event:"+required)

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
