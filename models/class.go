package models

import (
	"time"

	"gorm.io/gorm"
)

type Class struct {
	gorm.Model
	TeacherID          uint      `json:"teacher_id" gorm:"not null"`
	ClassName          string    `json:"class_name" gorm:"size:255;not null"`
	LearnerLimit       int       `json:"learner_limit" gorm:"not null;default:50"`
	ClassDescription   string    `json:"class_description" gorm:"size:1000"`
	BannerPictureURL   string    `json:"banner_picture,omitempty"`
	Price              float64   `json:"price" gorm:"type:numeric(12,2);default:0;check:price >= 0"`
	Rating             float64   `json:"rating" gorm:"type:numeric(3,2);default:0;check:rating >= 0 AND rating <= 5"`
	EnrollmentDeadline time.Time `json:"enrollment_deadline" gorm:"not null"`

	Teacher    Teacher         `gorm:"foreignKey:TeacherID;references:ID;constraint:OnDelete:CASCADE"`
	Categories []ClassCategory `gorm:"many2many:class_class_categories;constraint:OnDelete:CASCADE"`
}

// ---- DOC-ONLY STRUCT FOR SWAGGER BELOW ----

type ClassDoc struct {
	ID                 uint      `json:"id" example:"21"`
	TeacherID          uint      `json:"teacher_id" example:"7"`
	ClassName          string    `json:"class_name" example:"Advanced Python Programming"`
	LearnerLimit       int       `json:"learner_limit" example:"50"`
	ClassDescription   string    `json:"class_description" example:"Advanced Python programming course"`
	BannerPicture	   string    `json:"banner_picture,omitempty" example:"<base64-encoded-image>"`
	Price              float64   `json:"price" example:"1999.99"`
	Rating             float64   `json:"rating" example:"4.7"`
	EnrollmentDeadline time.Time `json:"enrollment_deadline" example:"2025-09-10T23:59:59Z"`
}
