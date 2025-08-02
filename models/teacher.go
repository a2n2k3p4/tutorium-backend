package models

import "gorm.io/gorm"

type Teacher struct {
	gorm.Model
	Description string `gorm:"size:255"`
}
