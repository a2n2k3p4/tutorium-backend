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

// ---- DOC-ONLY STRUCT FOR SWAGGER BELOW ----

type EnrollmentDoc struct {
	ID               uint   `json:"id" example:"101"`
	LearnerID        uint   `json:"learner_id" example:"42"`
	ClassSessionID   uint   `json:"class_session_id" example:"21"`
	EnrollmentStatus string `json:"enrollment_status" example:"active"`
}
