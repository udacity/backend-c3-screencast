package handlers

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	auth := r.Group("/auth")
	{
		auth.POST("/register", RegisterHandler)
		auth.POST("/login", LoginHandler)
		auth.POST("/refresh", RefreshHandler)
	}
}
