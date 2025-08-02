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
