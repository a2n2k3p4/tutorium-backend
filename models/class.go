package models

import (
	"time"

	"gorm.io/gorm"
)

type Class struct {
	gorm.Model
	TeacherID          uint      `gorm:"not null"`
	LearnerLimit       int       `gorm:"not null;default:50"`
	ClassDescription   string    `gorm:"size:1000"`
	EnrollmentDeadline time.Time `gorm:"not null"`

	Teacher Teacher `gorm:"foreignKey:TeacherID;references:ID;constraint:OnDelete:CASCADE"`
}
