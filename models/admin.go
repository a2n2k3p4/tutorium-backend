package models

import "gorm.io/gorm"

type Admin struct {
	gorm.Model
	UserID uint `gorm:"unique;not null"`
}

// ---- DOC-ONLY STRUCT FOR SWAGGER BELOW ----

type AdminDoc struct {
	ID     uint `json:"id" example:"43"`
	UserID uint `json:"user_id" example:"5"`
}
