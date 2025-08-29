package models

import (
	"gorm.io/gorm"
)

type Enrollment struct {
	gorm.Model
	LearnerID        uint   `json:"learner_id" gorm:"not null;uniqueIndex:idx_learner_class_enroll"`
	ClassID          uint   `json:"class_id" gorm:"not null;uniqueIndex:idx_learner_class_enroll"`
	EnrollmentStatus string `json:"enrollment_status" gorm:"size:20"`

	Learner Learner `gorm:"foreignKey:LearnerID;references:ID;constraint:OnDelete:CASCADE"`
	Class   Class   `gorm:"foreignKey:ClassID;references:ID;constraint:OnDelete:CASCADE"`
}

// ---- DOC-ONLY STRUCT FOR SWAGGER BELOW ----

type EnrollmentDoc struct {
	ID               uint   `json:"id" example:"101"`
	LearnerID        uint   `json:"learner_id" example:"42"`
	ClassID          uint   `json:"class_id" example:"21"`
	EnrollmentStatus string `json:"enrollment_status" example:"active"`
}
