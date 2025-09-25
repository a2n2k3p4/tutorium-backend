package models

import (
	"time"

	"gorm.io/gorm"
)

type Report struct {
	gorm.Model
	ReportUserID      uint      `json:"report_user_id"`
	ReportedUserID    uint      `json:"reported_user_id" gorm:"not null"`
	ClassSessionID    uint      `json:"class_session_id" gorm:"not null"`
	ReportType        string    `json:"report_type" gorm:"size:20;not null"`
	ReportReason      string    `json:"report_reason" gorm:"size:50;not null"`
	ReportDescription string    `json:"report_description" gorm:"size:255"`
	ReportPictureURL  string    `json:"report_picture,omitempty"`
	ReportDate        time.Time `json:"report_date" gorm:"default:CURRENT_TIMESTAMP"`
	ReportStatus      string    `json:"report_status" gorm:"size:10;default:'pending'"`
	ReportResult      string    `json:"report_result" gorm:"size:255"`

	Reporter     User         `gorm:"foreignKey:ReportUserID;constraint:OnDelete:SET NULL"`
	Reported     User         `gorm:"foreignKey:ReportedUserID;constraint:OnDelete:CASCADE"`
	ClassSession ClassSession `gorm:"foreignKey:ClassSessionID;constraint:OnDelete:CASCADE"`
}

// ---- DOC-ONLY STRUCT FOR SWAGGER BELOW ----

type ReportDoc struct {
	ID                uint      `json:"id" example:"12"`
	ReportUserID      uint      `json:"report_user_id" example:"5"`
	ReportedUserID    uint      `json:"reported_user_id" example:"8"`
	ClassSessionID    uint      `json:"class_session_id" example:"20"`
	ReportType        string    `json:"report_type" example:"Abuse"`
	ReportReason      string    `json:"report_reason" example:"teacher_absent"`
	ReportDescription string    `json:"report_description" example:"User sent inappropriate messages"`
	ReportPicture     string    `json:"report_picture,omitempty" example:"<base64-encoded-image>"`
	ReportDate        time.Time `json:"report_date" example:"2025-08-20T14:30:00Z"`
	ReportStatus      string    `json:"report_status" example:"pending"`
}
