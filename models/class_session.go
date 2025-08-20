package models

import (
	"time"

	"gorm.io/gorm"
)

type ClassSession struct {
	gorm.Model
	ClassID            uint
	Description        string    `gorm:"size:1000"`
	EnrollmentDeadline time.Time `gorm:"not null"`
	ClassStart         time.Time `gorm:"not null"`
	ClassFinish        time.Time `gorm:"not null"`
	ClassStatus        string    `gorm:"size:20"`

	Class Class `gorm:"foreignKey:ClassID;references:ID;constraint:OnDelete:CASCADE"`
}

// ---- DOC-ONLY STRUCT FOR SWAGGER BELOW ----

type ClassSessionDoc struct {
	ID                 uint      `json:"id" example:"15"`
	ClassID            uint      `json:"class_id" example:"12"`
	Description        string    `json:"description" example:"Weekly tutoring session for calculus"`
	EnrollmentDeadline time.Time `json:"enrollment_deadline" example:"2025-09-01T23:59:59Z"`
	ClassStart         time.Time `json:"class_start" example:"2025-09-05T14:00:00Z"`
	ClassFinish        time.Time `json:"class_finish" example:"2025-09-05T16:00:00Z"`
	ClassStatus        string    `json:"class_status" example:"Scheduled"`
}
