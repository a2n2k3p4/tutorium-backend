package models

import (
	"time"

	"gorm.io/gorm"
)

type BanDetailsLearner struct {
	gorm.Model
	LearnerID      uint      `json:"learner_id" gorm:"not null"`
	BanStart       time.Time `json:"ban_start" gorm:"default:CURRENT_TIMESTAMP"`
	BanEnd         time.Time `json:"ban_end" gorm:"not null"`
	BanDescription string    `json:"ban_description" gorm:"size:255"`

	Learner Learner `gorm:"foreignKey:LearnerID;references:ID;constraint:OnDelete:CASCADE"`
}

// ---- DOC-ONLY STRUCT FOR SWAGGER BELOW ----

type BanDetailsLearnerDoc struct {
	LearnerID      uint      `json:"learner_id" example:"42"`
	BanStart       time.Time `json:"ban_start" example:"2025-08-20T12:00:00Z"`
	BanEnd         time.Time `json:"ban_end" example:"2025-08-30T12:00:00Z"`
	BanDescription string    `json:"ban_description" example:"Spamming inappropriate content"`
}
