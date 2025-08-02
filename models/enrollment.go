package models

import (
	"gorm.io/gorm"
)

type Enrollment struct {
	gorm.Model
	LearnerID        uint   `gorm:"not null;uniqueIndex:idx_learner_class_enroll"`
	ClassID          uint   `gorm:"not null;uniqueIndex:idx_learner_class_enroll"`
	EnrollmentStatus string `gorm:"size:20"`

	Learner Learner `gorm:"foreignKey:LearnerID;references:ID;constraint:OnDelete:CASCADE"`
	Class   Class   `gorm:"foreignKey:ClassID;references:ID;constraint:OnDelete:CASCADE"`
}
