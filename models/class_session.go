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
