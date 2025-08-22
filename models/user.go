package models

import (
	"gorm.io/gorm"
)

// Define a struct matching the columns (use pointers for nullable FKs)
type User struct {
	gorm.Model
	StudentID      string `gorm:"size:10;uniqueIndex;not null"`
	ProfilePicture []byte
	FirstName      string  `gorm:"size:30;not null"`
	LastName       string  `gorm:"size:30;not null"`
	Gender         string  `gorm:"size:6"`
	PhoneNumber    string  `gorm:"size:20"`
	Balance        float64 `gorm:"type:numeric(12,2);default:0;check:balance >= 0"`

	Learner *Learner
	Teacher *Teacher
	Admin   *Admin
}

// ---- DOC-ONLY STRUCT FOR SWAGGER BELOW ----

type UserDoc struct {
	ID             uint    `json:"id" example:"101"`
	StudentID      string  `json:"student_id" example:"6610505511"`
	ProfilePicture string  `json:"profile_picture,omitempty" example:"<base64-encoded-image>"`
	FirstName      string  `json:"first_name" example:"Alice"`
	LastName       string  `json:"last_name" example:"Smith"`
	Gender         string  `json:"gender" example:"Female"`
	PhoneNumber    string  `json:"phone_number" example:"+66912345678"`
	Balance        float64 `json:"balance" example:"250.75"`
}
