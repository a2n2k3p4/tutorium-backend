package models

import (
	"time"

	"gorm.io/gorm"
)

type Review struct {
	LearnerID uint   `gorm:"primaryKey"`
	ClassID   uint   `gorm:"primaryKey"`
	Rating    int    `gorm:"not null;check:rating >= 1 AND rating <= 5"`
	Comment   string `gorm:"size:255"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Learner Learner `gorm:"foreignKey:LearnerID;references:ID;constraint:OnDelete:CASCADE"`
	Class   Class   `gorm:"foreignKey:ClassID;references:ID;constraint:OnDelete:CASCADE"`
}
