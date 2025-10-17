package middlewares

import (
	"errors"
	"time"

	"github.com/a2n2k3p4/tutorium-backend/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// BanMiddleware checks if the authenticated user has an active ban in any role.
func BanMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if Status() == "development" {
			return c.Next()
		}
		user, ok := c.Locals("currentUser").(*models.User)
		if !ok {
			return c.Next()
		}

		db, err := GetDB(c)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "database not available"})
		}

		banned, banEnd, err := isUserBanned(db, user)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to check ban status"})
		}
		if banned {
			return c.Status(403).JSON(fiber.Map{
				"error": "Your account is suspended until " + banEnd.Format(time.RFC3339),
			})
		}
		return c.Next()
	}
}

// isUserBanned is a helper function that checks all roles for an active ban.
func isUserBanned(db *gorm.DB, user *models.User) (bool, time.Time, error) {
	if user.Teacher != nil {
		var teacherBan models.BanDetailsTeacher
		err := db.Where("teacher_id = ? AND ban_end > ?", user.Teacher.ID, time.Now()).First(&teacherBan).Error
		if err == nil {
			return true, teacherBan.BanEnd, nil
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return false, time.Time{}, err
		}
	}

	if user.Learner != nil {
		var learnerBan models.BanDetailsLearner
		err := db.Where("learner_id = ? AND ban_end > ?", user.Learner.ID, time.Now()).First(&learnerBan).Error
		if err == nil {
			return true, learnerBan.BanEnd, nil
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return false, time.Time{}, err
		}
	}

	return false, time.Time{}, nil
}
