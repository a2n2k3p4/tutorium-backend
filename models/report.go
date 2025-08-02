package models

import (
	"time"

	"gorm.io/gorm"
)

type Report struct {
	gorm.Model
	ReportUserID      uint   `gorm:"not null"`
	ReportedUserID    uint   `gorm:"not null"`
	ReportType        string `gorm:"size:20;not null"`
	ReportDescription string `gorm:"size:255"`
	ReportPicture     []byte
	ReportDate        time.Time `gorm:"default:CURRENT_TIMESTAMP"`

	Reporter User `gorm:"foreignKey:ReportUserID;constraint:OnDelete:CASCADE"`
	Reported User `gorm:"foreignKey:ReportedUserID;constraint:OnDelete:CASCADE"`
}
