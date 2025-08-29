package models

import "gorm.io/gorm"

type Teacher struct {
	gorm.Model
	UserID      uint   `json:"user_id" gorm:"unique;not null"`
	Description string `json:"description" gorm:"size:255"`
}

// ---- DOC-ONLY STRUCT FOR SWAGGER BELOW ----

type TeacherDoc struct {
	ID          uint   `json:"id" example:"12"`
	UserID      uint   `json:"user_id" example:"5"`
	Description string `json:"description" example:"Experienced Mathematics teacher specializing in calculus and linear algebra."`
}
