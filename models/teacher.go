package models

import "gorm.io/gorm"

type Teacher struct {
	gorm.Model
	UserID      uint   `json:"user_id" gorm:"unique;not null"`
	Description string `json:"description" gorm:"size:255"`
	FlagCount   int    `json:"flag_count" gorm:"default:0;not null"`
	BanCount    int    `json:"ban_count" gorm:"default:0;not null"`
	Email       string `json:"email" gorm:"size:100;unique;not null"`
}

// ---- DOC-ONLY STRUCT FOR SWAGGER BELOW ----

type TeacherDoc struct {
	ID          uint   `json:"id" example:"12"`
	UserID      uint   `json:"user_id" example:"5"`
	Description string `json:"description" example:"Experienced Mathematics teacher specializing in calculus and linear algebra."`
	FlagCount   int    `json:"flag_count" example:"3"`
	BanCount    int    `json:"ban_count" example:"1"`
	Email       string `json:"email" example:"teacher@example.com"`
}
