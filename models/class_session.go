package models

import (
	"time"

	"gorm.io/gorm"
)

// ---- CreateClassSessionRequest temporary holds the POST API json data ----
type CreateClassSessionRequest struct {
	ClassID            uint      `json:"class_id"`
	Description        string    `json:"description"`
	Price              float64   `json:"price" `
	LearnerLimit       int       `json:"learner_limit"`
	EnrollmentDeadline time.Time `json:"enrollment_deadline"`
	ClassStart         time.Time `json:"class_start"`
	ClassFinish        time.Time `json:"class_finish"`
	ClassStatus        string    `json:"class_status"`
}

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
	MeetingUrl         string    `json:"class_url" gorm:"size:128"`

	Class Class `gorm:"foreignKey:ClassID;references:ID;constraint:OnDelete:CASCADE"`
}

// ---- DOC-ONLY STRUCT FOR SWAGGER BELOW ----

type CreateClassSessionRequestDoc struct {
	ClassID            uint      `json:"class_id" example:"12"`
	Description        string    `json:"description" example:"Weekly tutoring session for calculus"`
	Price              float64   `json:"price" example:"1999.99"`
	LearnerLimit       int       `json:"learner_limit" example:"5"`
	EnrollmentDeadline time.Time `json:"enrollment_deadline" example:"2025-09-01T23:59:59Z"`
	ClassStart         time.Time `json:"class_start" example:"2025-09-05T14:00:00Z"`
	ClassFinish        time.Time `json:"class_finish" example:"2025-09-05T16:00:00Z"`
	ClassStatus        string    `json:"class_status" example:"Scheduled"`
}

type ClassSessionDoc struct {
	ClassID            uint      `json:"class_id" example:"12"`
	Description        string    `json:"description" example:"Weekly tutoring session for calculus"`
	Price              float64   `json:"price" example:"199.99"`
	LearnerLimit       int       `json:"learner_limit" example:"50"`
	EnrollmentDeadline time.Time `json:"enrollment_deadline" example:"2025-09-01T23:59:59Z"`
	ClassStart         time.Time `json:"class_start" example:"2025-09-05T14:00:00Z"`
	ClassFinish        time.Time `json:"class_finish" example:"2025-09-05T16:00:00Z"`
	ClassStatus        string    `json:"class_status" example:"Scheduled"`
	MeetingUrl         string    `json:"class_url" example:"https://meet.jit.si/KUtutorium_12_1758630058"`
}
