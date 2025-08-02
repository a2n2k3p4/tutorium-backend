package models

import (
	"time"

	"gorm.io/gorm"
)

type BanDetailsTeacher struct {
	gorm.Model
	TeacherID      uint      `gorm:"not null"`
	BanStart       time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	BanEnd         time.Time `gorm:"not null"`
	BanDescription string    `gorm:"size:255"`

	Teacher Teacher `gorm:"foreignKey:TeacherID;references:ID;constraint:OnDelete:CASCADE"`
}
