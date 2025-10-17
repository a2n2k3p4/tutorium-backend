package services

import (
	"github.com/a2n2k3p4/tutorium-backend/models"
	"gorm.io/gorm"
)

// Creates and saves a new notification for a specific user.
func CreateNotification(db *gorm.DB, userID uint, notifType string, description string) error {
	notification := models.Notification{
		UserID:                  userID,
		NotificationType:        notifType,
		NotificationDescription: description,
	}
	return db.Create(&notification).Error
}
