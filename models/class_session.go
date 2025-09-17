package models

import (
	"time"

	"gorm.io/gorm"
)

type ClassSession struct {
	gorm.Model
	ClassID            uint      `json:"class_id"`
	Description        string    `json:"description" gorm:"size:1000"`
	Price              float64   `json:"price" gorm:"type:numeric(12,2);default:0;check:price >= 0"`
	LearnerLimit       int       `json:"learner_limit" gorm:"not null;default:50"`
	EnrollmentDeadline time.Time `json:"enrollment_deadline" gorm:"not null"`
	ClassStart         time.Time `json:"class_start" gorm:"not null"`
	ClassFinish        time.Time `json:"class_finish" gorm:"not null"`
	ClassStatus        string    `json:"class_status" gorm:"size:20"`

	Class Class `gorm:"foreignKey:ClassID;references:ID;constraint:OnDelete:CASCADE"`
}

// ---- DOC-ONLY STRUCT FOR SWAGGER BELOW ----

type ClassSessionDoc struct {
	ID                 uint      `json:"id" example:"15"`
	ClassID            uint      `json:"class_id" example:"12"`
	Description        string    `json:"description" example:"Weekly tutoring session for calculus"`
	Price              float64   `json:"price" example:"1999.99"`
	LearnerLimit       int       `json:"learner_limit" example:"50"`
	EnrollmentDeadline time.Time `json:"enrollment_deadline" example:"2025-09-01T23:59:59Z"`
	ClassStart         time.Time `json:"class_start" example:"2025-09-05T14:00:00Z"`
	ClassFinish        time.Time `json:"class_finish" example:"2025-09-05T16:00:00Z"`
	ClassStatus        string    `json:"class_status" example:"Scheduled"`
}
