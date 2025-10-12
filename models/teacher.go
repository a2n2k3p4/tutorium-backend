package models

import "gorm.io/gorm"

type Teacher struct {
	gorm.Model
	UserID      uint   `json:"user_id" gorm:"unique;not null"`
	Description string `json:"description" gorm:"size:255"`
	FlagCount   int    `json:"flag_count" gorm:"default:0;not null"`
	Email       string `json:"email" gorm:"size:100;unique;not null"`
}

// ---- DOC-ONLY STRUCT FOR SWAGGER BELOW ----

type TeacherDoc struct {
	UserID      uint   `json:"user_id" example:"5"`
	Description string `json:"description" example:"Experienced Mathematics teacher specializing in calculus and linear algebra."`
	FlagCount   int    `json:"flag_count" example:"3"`
	Email       string `json:"email" example:"teacher@example.com"`
}

type TeacherAverageRating struct {
	TeacherID     uint    `json:"teacher_id,omitempty" example:"123"`
	AverageRating float64 `json:"average_rating" example:"4.5"`
}
