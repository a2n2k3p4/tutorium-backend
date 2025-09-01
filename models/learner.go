package models

import "gorm.io/gorm"

type Learner struct {
	gorm.Model
	UserID    uint `json:"user_id" gorm:"unique;not null"`
	FlagCount int  `json:"flag_count" gorm:"default:0;not null"`
}

// ---- DOC-ONLY STRUCT FOR SWAGGER BELOW ----

type LearnerDoc struct {
	ID        uint `json:"id" example:"42"`
	UserID    uint `json:"user_id" example:"5"`
	FlagCount int  `json:"flag_count" example:"3"`
}
