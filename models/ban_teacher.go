package models

import (
	"time"

	"gorm.io/gorm"
)

type BanDetailsTeacher struct {
	gorm.Model
	TeacherID      uint      `json:"teacher_id" gorm:"not null"`
	BanStart       time.Time `json:"ban_start" gorm:"default:CURRENT_TIMESTAMP"`
	BanEnd         time.Time `json:"ban_end" gorm:"not null"`
	BanDescription string    `json:"ban_description" gorm:"size:255"`

	Teacher Teacher `gorm:"foreignKey:TeacherID;references:ID;constraint:OnDelete:CASCADE"`
}

// ---- DOC-ONLY STRUCT FOR SWAGGER BELOW ----

type BanDetailsTeacherDoc struct {
	ID             uint      `json:"id" example:"1"`
	TeacherID      uint      `json:"teacher_id" example:"7"`
	BanStart       time.Time `json:"ban_start" example:"2025-08-20T12:00:00Z"`
	BanEnd         time.Time `json:"ban_end" example:"2025-08-30T12:00:00Z"`
	BanDescription string    `json:"ban_description" example:"Repeated policy violations"`
}
