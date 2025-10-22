package models

import (
	"time"

	"gorm.io/gorm"
)

type Notification struct {
	gorm.Model
	UserID                  uint      `json:"user_id" gorm:"not null"`
	NotificationTitle       string    `json:"notification_title" gorm:"size:100;not null"`
	NotificationType        string    `json:"notification_type" gorm:"size:30;not null"`
	NotificationDescription string    `json:"notification_description" gorm:"size:255"`
	NotificationDate        time.Time `json:"notification_date" gorm:"default:CURRENT_TIMESTAMP"`
	ReadFlag                bool      `json:"read_flag" gorm:"default:false"`

	User User `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
}

// ---- DOC-ONLY STRUCT FOR SWAGGER BELOW ----

type NotificationDoc struct {
	UserID                  uint      `json:"user_id" example:"42"`
	NotificationTitle       string    `json:"notification_title" example:"Review received"`
	NotificationType        string    `json:"notification_type" example:"System Alert"`
	NotificationDescription string    `json:"notification_description" example:"Your class has been rescheduled"`
	NotificationDate        time.Time `json:"notification_date" example:"2025-08-20T15:04:05Z"`
	ReadFlag                bool      `json:"read_flag" example:"false"`
}
