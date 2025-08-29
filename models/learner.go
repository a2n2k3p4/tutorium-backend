package models

import "gorm.io/gorm"

type Learner struct {
	gorm.Model
	UserID uint `json:"user_id" gorm:"unique;not null"`
}

// ---- DOC-ONLY STRUCT FOR SWAGGER BELOW ----

type LearnerDoc struct {
	ID     uint `json:"id" example:"42"`
	UserID uint `json:"user_id" example:"5"`
}
