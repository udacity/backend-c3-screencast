package handlers

import (
	"eventure/internal/database"
	"eventure/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ListUsersHandler(c *gin.Context) {
	var users []models.User
	if err := database.DB.Preload("Role").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	var output []gin.H
	for _, u := range users {
		output = append(output, gin.H{
			"id":       u.ID,
			"username": u.Username,
			"role":     u.Role.Name,
		})
	}

	c.JSON(http.StatusOK, output)
}
