package models

import "gorm.io/gorm"

type ClassCategory struct {
	gorm.Model
	ClassCategory string `gorm:"size:30;primaryKey"`

	Class Class `gorm:"foreignKey:ID;references:ID;constraint:OnDelete:CASCADE"`
}
