package models

import "gorm.io/gorm"

type ClassCategory struct {
	gorm.Model
	ClassCategory string `gorm:"size:30"`

	Classes []Class `gorm:"many2many:class_class_categories;constraint:OnDelete:CASCADE"`
}
