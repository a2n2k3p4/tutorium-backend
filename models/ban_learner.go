package models

import (
	"time"

	"gorm.io/gorm"
)

type BanDetailsLearner struct {
	gorm.Model
	LearnerID      uint      `gorm:"not null"`
	BanStart       time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	BanEnd         time.Time `gorm:"not null"`
	BanDescription string    `gorm:"size:255"`

	Learner Learner `gorm:"foreignKey:LearnerID;references:ID;constraint:OnDelete:CASCADE"`
}
