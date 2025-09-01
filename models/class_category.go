package models

import "gorm.io/gorm"

type ClassCategory struct {
	gorm.Model
	ClassCategory string `json:"class_category" gorm:"size:30"`

	Classes []Class `gorm:"many2many:class_class_categories;constraint:OnDelete:CASCADE"`
}

// ---- DOC-ONLY STRUCT FOR SWAGGER BELOW ----

type ClassCategoryDoc struct {
	ID            uint       `json:"id" example:"3"`
	ClassCategory string     `json:"class_category" example:"Mathematics"`
	Classes       []ClassDoc `json:"classes,omitempty"`
}
