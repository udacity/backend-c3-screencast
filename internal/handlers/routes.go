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

	// Temporary (for testing)
	admin := r.Group("/admin")
	admin.Use(middleware.RequireAuth(), middleware.RequireRole("admin"))
	{
		admin.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "pong"})
		})
	}
}
