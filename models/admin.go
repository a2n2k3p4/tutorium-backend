package models

import "gorm.io/gorm"

type Admin struct {
	gorm.Model
	UserID uint `json:"user_id" gorm:"unique;not null"`
}

// ---- DOC-ONLY STRUCT FOR SWAGGER BELOW ----

type AdminDoc struct {
	UserID uint `json:"user_id" example:"5"`
}
