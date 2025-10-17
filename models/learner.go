package models

import "gorm.io/gorm"

type Learner struct {
	gorm.Model
	UserID     uint            `json:"user_id" gorm:"unique;not null"`
	FlagCount  int             `json:"flag_count" gorm:"default:0;not null"`
	Interested []ClassCategory `gorm:"many2many:interested_class_categories;constraint:OnDelete:CASCADE"`
}

// ---- DOC-ONLY STRUCT FOR SWAGGER BELOW ----

type LearnerDoc struct {
	UserID    uint `json:"user_id" example:"1"`
	FlagCount int  `json:"flag_count" example:"1"`
}
