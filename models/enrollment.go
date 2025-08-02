package models

import (
	"time"

	"gorm.io/gorm"
)

type Enrollment struct {
	LearnerID        uint   `gorm:"primaryKey"`
	ClassID          uint   `gorm:"primaryKey"`
	EnrollmentStatus string `gorm:"size:20"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Learner Learner `gorm:"foreignKey:LearnerID;references:ID;constraint:OnDelete:CASCADE"`
	Class   Class   `gorm:"foreignKey:ClassID;references:ID;constraint:OnDelete:CASCADE"`
}
