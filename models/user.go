package models

import (
	"gorm.io/gorm"
)

// Define a struct matching the columns (use pointers for nullable FKs)
type User struct {
	gorm.Model
	SessionID      string `gorm:"size:255"`
	ProfilePicture []byte
	FirstName      string  `gorm:"size:30;not null"`
	LastName       string  `gorm:"size:30;not null"`
	Gender         string  `gorm:"size:6"`
	PhoneNumber    string  `gorm:"size:20"`
	Balance        float64 `gorm:"type:numeric(12,2);default:0;check:balance >= 0"`

	LearnerID *uint `gorm:"unique"`
	TeacherID *uint `gorm:"unique"`
	AdminID   *uint `gorm:"unique"`

	Learner *Learner `gorm:"foreignKey:LearnerID;references:ID;constraint:OnDelete:SET NULL;"`
	Teacher *Teacher `gorm:"foreignKey:TeacherID;references:ID;constraint:OnDelete:SET NULL;"`
	Admin   *Admin   `gorm:"foreignKey:AdminID;references:ID;constraint:OnDelete:SET NULL;"`
}

// ---- DOC-ONLY STRUCT FOR SWAGGER BELOW ----

type UserDoc struct {
	ID             uint    `json:"id" example:"101"`
	SessionID      string  `json:"session_id" example:"sess-abc123xyz"`
	ProfilePicture string  `json:"profile_picture,omitempty" example:"<base64-encoded-image>"`
	FirstName      string  `json:"first_name" example:"Alice"`
	LastName       string  `json:"last_name" example:"Smith"`
	Gender         string  `json:"gender" example:"Female"`
	PhoneNumber    string  `json:"phone_number" example:"+66912345678"`
	Balance        float64 `json:"balance" example:"250.75"`
	LearnerID      *uint   `json:"learner_id,omitempty" example:"42"`
	TeacherID      *uint   `json:"teacher_id,omitempty" example:"7"`
	AdminID        *uint   `json:"admin_id,omitempty" example:"3"`
}
