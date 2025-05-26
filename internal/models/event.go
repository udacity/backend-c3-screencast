package models

import "gorm.io/gorm"

type Event struct {
	gorm.Model
	Title   string
	Details string
	OwnerID uint
	Owner   User `gorm:"foreignKey:OwnerID"`
}
