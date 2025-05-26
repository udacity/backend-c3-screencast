package handlers

import (
	"eventure/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	auth := r.Group("/auth")
	{
		auth.POST("/register", RegisterHandler)
		auth.POST("/login", LoginHandler)
		auth.POST("/refresh", RefreshHandler)
		auth.GET("/profile", middleware.RequireAuth(), ProfileHandler)
	}

	admin := r.Group("/admin")
	admin.Use(middleware.RequireAuth(), middleware.RequireRole("admin"))
	{
		admin.GET("/users", ListUsersHandler)
	}

	event := r.Group("/events")
	event.Use(middleware.RequireAuth())
	{
		event.POST("", CreateEventHandler)
		event.PUT("/:id", middleware.RequireEventOwnershipOrPermission("update"), UpdateEventHandler)
		event.DELETE("/:id", middleware.RequireEventOwnershipOrPermission("delete"), DeleteEventHandler)
	}
}
