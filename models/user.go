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
