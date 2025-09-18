package models

import (
	"gorm.io/gorm"
)

// Define a struct matching the columns (use pointers for nullable FKs)
type User struct {
	gorm.Model
	StudentID         string  `json:"student_id" gorm:"size:20;uniqueIndex;not null"`
	ProfilePictureURL string  `json:"profile_picture,omitempty"`
	FirstName         string  `json:"first_name" gorm:"size:30;not null"`
	LastName          string  `json:"last_name" gorm:"size:30;not null"`
	Gender            string  `json:"gender" gorm:"size:6"`
	PhoneNumber       string  `json:"phone_number" gorm:"size:20"`
	Balance           float64 `json:"balance" gorm:"type:numeric(12,2);default:0;check:balance >= 0"`
	BanCount          int     `json:"ban_count" gorm:"default:0;not null"`

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
	BanCount       int     `json:"ban_count" example:"1"`
}
