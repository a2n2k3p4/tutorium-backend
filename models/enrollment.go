package models

import (
	"gorm.io/gorm"
)

type Enrollment struct {
	gorm.Model
	LearnerID        uint   `json:"learner_id" gorm:"not null;uniqueIndex:idx_learner_session"`
	ClassSessionID   uint   `json:"class_session_id" gorm:"not null;uniqueIndex:idx_learner_session"`
	EnrollmentStatus string `json:"enrollment_status" gorm:"size:20"`

	Learner      Learner      `gorm:"foreignKey:LearnerID;references:ID;constraint:OnDelete:CASCADE"`
	ClassSession ClassSession `gorm:"foreignKey:ClassSessionID;references:ID;constraint:OnDelete:CASCADE"`
}

type EnrollmentResponse struct {
	Enrollment
	User *User `json:"user,omitempty"`
}

// ---- DOC-ONLY STRUCT FOR SWAGGER BELOW ----

type EnrollmentDoc struct {
	LearnerID        uint   `json:"learner_id" example:"1"`
	ClassSessionID   uint   `json:"class_session_id" example:"3"`
	EnrollmentStatus string `json:"enrollment_status" example:"active"`
}
