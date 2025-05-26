package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CustomRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				// Log the panic for internal tracking
				log.Printf("[PANIC RECOVERED] %v\n", r)

				// Respond with a generic message
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "Internal server error",
				})
			}
		}()

		// Continue down the chain
		c.Next()
	}
}
