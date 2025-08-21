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

// ---- DOC-ONLY STRUCT FOR SWAGGER BELOW ----

type ReportDoc struct {
	ID                uint      `json:"id" example:"12"`
	ReportUserID      uint      `json:"report_user_id" example:"5"`
	ReportedUserID    uint      `json:"reported_user_id" example:"8"`
	ReportType        string    `json:"report_type" example:"Abuse"`
	ReportDescription string    `json:"report_description" example:"User sent inappropriate messages"`
	ReportPicture     string    `json:"report_picture" example:"base64-encoded-image-string"`
	ReportDate        time.Time `json:"report_date" example:"2025-08-20T14:30:00Z"`
}
