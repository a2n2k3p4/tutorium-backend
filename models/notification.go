package models

import (
	"time"

	"gorm.io/gorm"
)

type Notification struct {
	gorm.Model
	UserID                  uint      `gorm:"not null"`
	NotificationType        string    `gorm:"size:30;not null"`
	NotificationDescription string    `gorm:"size:255"`
	NotificationDate        time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	ReadFlag                bool      `gorm:"default:false"`

	User User `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
}

// ---- DOC-ONLY STRUCT FOR SWAGGER BELOW ----

type NotificationDoc struct {
	ID                      uint      `json:"id" example:"100"`
	UserID                  uint      `json:"user_id" example:"42"`
	NotificationType        string    `json:"notification_type" example:"System Alert"`
	NotificationDescription string    `json:"notification_description" example:"Your class has been rescheduled"`
	NotificationDate        time.Time `json:"notification_date" example:"2025-08-20T15:04:05Z"`
	ReadFlag                bool      `json:"read_flag" example:"false"`
}
