package database

import (
	"fmt"
	"os"

	"eventure/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() error {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		return fmt.Errorf("DB_DSN is not set in environment variables")
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	fmt.Println("Connected to Postgres!")

	if err := DB.AutoMigrate(&models.User{}); err != nil {
		return fmt.Errorf("AutoMigrate failed: %w", err)
	}

	return nil
}
