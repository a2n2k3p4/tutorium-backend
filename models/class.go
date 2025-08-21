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

	Teacher    Teacher         `gorm:"foreignKey:TeacherID;references:ID;constraint:OnDelete:CASCADE"`
	Categories []ClassCategory `gorm:"many2many:class_class_categories;constraint:OnDelete:CASCADE"`
}

// ---- DOC-ONLY STRUCT FOR SWAGGER BELOW ----

type ClassDoc struct {
	ID                 uint      `json:"id" example:"21"`
	TeacherID          uint      `json:"teacher_id" example:"7"`
	LearnerLimit       int       `json:"learner_limit" example:"50"`
	ClassDescription   string    `json:"class_description" example:"Advanced Python programming course"`
	EnrollmentDeadline time.Time `json:"enrollment_deadline" example:"2025-09-10T23:59:59Z"`
}
