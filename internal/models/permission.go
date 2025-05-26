package models

import "gorm.io/gorm"

type Permission struct {
	gorm.Model
	Name  string `gorm:"uniqueIndex;not null"`
	Roles []Role `gorm:"many2many:role_permissions"`
}
