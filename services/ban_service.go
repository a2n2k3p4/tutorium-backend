package services

import (
	"time"

	"github.com/a2n2k3p4/tutorium-backend/models"
	"gorm.io/gorm"
)

// ApplyLearnerFlags adds flags to a learner. If the flag count reaches 3,
// it issues a 7-day ban and updates the main user's total ban count.
func ApplyLearnerFlags(db *gorm.DB, learnerID uint, flagsToAdd int, reason string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var learner models.Learner
		if err := tx.First(&learner, learnerID).Error; err != nil {
			return err
		}

		learner.FlagCount += flagsToAdd

		if learner.FlagCount >= 3 {
			var user models.User
			if err := tx.First(&user, learner.UserID).Error; err != nil {
				return err
			}

			learner.FlagCount -= 3
			user.BanCount += 1

			banDetails := models.BanDetailsLearner{
				LearnerID:      learnerID,
				BanStart:       time.Now(),
				BanEnd:         time.Now().Add(7 * 24 * time.Hour),
				BanDescription: reason,
			}
			if err := tx.Create(&banDetails).Error; err != nil {
				return err
			}
			if err := tx.Save(&user).Error; err != nil {
				return err
			}
		}
		return tx.Save(&learner).Error
	})
}

// ApplyTeacherFlags adds flags to a teacher. If the flag count reaches 3,
// it issues a 7-day ban and updates the main user's total ban count.
func ApplyTeacherFlags(db *gorm.DB, teacherID uint, flagsToAdd int, reason string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var teacher models.Teacher
		if err := tx.First(&teacher, teacherID).Error; err != nil {
			return err
		}

		teacher.FlagCount += flagsToAdd

		if teacher.FlagCount >= 3 {
			var user models.User
			if err := tx.First(&user, teacher.UserID).Error; err != nil {
				return err
			}

			teacher.FlagCount -= 3
			user.BanCount += 1

			banDetails := models.BanDetailsTeacher{
				TeacherID:      teacherID,
				BanStart:       time.Now(),
				BanEnd:         time.Now().Add(7 * 24 * time.Hour),
				BanDescription: reason,
			}
			if err := tx.Create(&banDetails).Error; err != nil {
				return err
			}
			if err := tx.Save(&user).Error; err != nil {
				return err
			}
		}
		return tx.Save(&teacher).Error
	})
}

// AddTeacherFlag is a simple function to add a flag without triggering a ban.
// Used by the cron job when a ban is not imminent.
func AddTeacherFlag(db *gorm.DB, teacherID uint, flagsToAdd int) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var teacher models.Teacher
		if err := tx.First(&teacher, teacherID).Error; err != nil {
			return err
		}
		teacher.FlagCount += flagsToAdd
		return tx.Save(&teacher).Error
	})
}
