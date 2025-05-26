package handlers

import (
	"eventure/internal/database"
	"eventure/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type EventInput struct {
	Title   string `json:"title"`
	Details string `json:"details"`
}

func CreateEventHandler(c *gin.Context) {
	userRaw, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	user := userRaw.(models.User)

	var input EventInput
	if err := c.ShouldBindJSON(&input); err != nil || input.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	event := models.Event{
		Title:   input.Title,
		Details: input.Details,
		OwnerID: user.ID,
	}

	if err := database.DB.Create(&event).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create event"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":      event.ID,
		"title":   event.Title,
		"details": event.Details,
		"owner":   user.Username,
	})
}

func UpdateEventHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	var event models.Event
	if err := database.DB.First(&event, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	var input EventInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	event.Title = input.Title
	event.Details = input.Details

	if err := database.DB.Save(&event).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update event"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Event updated",
		"id":      event.ID,
	})
}

func DeleteEventHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	if err := database.DB.Delete(&models.Event{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete event"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Event deleted",
		"id":      id,
	})
}

func SearchEventsHandler(c *gin.Context) {
	query := c.Query("title")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter 'title' is required"})
		return
	}

	var unsafeResults []models.Event
	var safeResults []models.Event

	// Unsafe: string interpolation can be exploited
	// rawQuery := fmt.Sprintf("SELECT * FROM events WHERE title LIKE '%%%s%%'", query)
	// if err := database.DB.Raw(rawQuery).Scan(&unsafeResults).Error; err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Unsafe query failed"})
	// 	return
	// }

	// Safe: uses parameterized queries
	if err := database.DB.Where("title ILIKE ?", "%"+query+"%").Find(&safeResults).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Safe query failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"unsafe_results": unsafeResults,
		"safe_results":   safeResults,
	})
}
