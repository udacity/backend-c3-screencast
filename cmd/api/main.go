package main

import (
	"log"
	"os"

	"eventure/internal/database"
	"eventure/internal/handlers"
	"eventure/internal/middleware"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found or couldn't be loaded")
	}

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(middleware.CustomRecovery())
	r.Use(middleware.SecureHeaders())

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Front end origin
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	err = database.Connect()
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	handlers.RegisterRoutes(r)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s...", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
