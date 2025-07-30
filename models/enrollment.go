package models

import "gorm.io/gorm"

type Enrollment struct {
	gorm.Model
	LearnerID        uint   `gorm:"primaryKey"`
	EnrollmentStatus string `gorm:"size:20"`

	Learner Learner `gorm:"foreignKey:LearnerID;references:ID;constraint:OnDelete:CASCADE"`
	Class   Class   `gorm:"foreignKey:ClassID;references:ID;constraint:OnDelete:CASCADE"`
}
