package models

import "gorm.io/gorm"

type Admin struct {
	gorm.Model
}

// ---- DOC-ONLY STRUCT FOR SWAGGER BELOW ----

type AdminDoc struct {
	ID uint `json:"id" example:"43"`
}
