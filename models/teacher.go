package models

import "gorm.io/gorm"

type Teacher struct {
	gorm.Model
	Description string `gorm:"size:255"`
}

// ---- DOC-ONLY STRUCT FOR SWAGGER BELOW ----

type TeacherDoc struct {
	ID          uint   `json:"id" example:"12"`
	Description string `json:"description" example:"Experienced Mathematics teacher specializing in calculus and linear algebra."`
}
