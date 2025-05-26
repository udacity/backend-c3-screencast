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

	// Test route for permission-based middleware
	r.GET("/events", middleware.RequireAuth(), middleware.RequirePermission("event", "read"), func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Access granted to view events!"})
	})
}
