package services

import (
	"log"

	"fmt"
	"time"

	"github.com/a2n2k3p4/tutorium-backend/models"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

// StartScheduler initializes all cron jobs for the application.
func StartScheduler(db *gorm.DB) {
	c := cron.New()
	c.AddFunc("@every 5m", func() {
		log.Println("Running teacher absence checker job...")
		CheckForAbsentTeachers(db)
	})
	c.Start()
	log.Println("Cron job scheduler started.")
}

// CheckForAbsentTeachers finds absent teachers and either flags them or creates a report.
func CheckForAbsentTeachers(db *gorm.DB) {
	var sessions []models.ClassSession
	fifteenMinutesAgo := time.Now().Add(-15 * time.Minute)

	err := db.Where("class_start < ? AND class_status = ?", fifteenMinutesAgo, "not_started").Find(&sessions).Error
	if err != nil {
		log.Printf("Error finding absent sessions: %v", err)
		return
	}

	for _, session := range sessions {
		var class models.Class
		if err := db.First(&class, session.ClassID).Error; err != nil {
			log.Printf("Failed to find class %d: %v", session.ClassID, err)
			continue
		}

		var teacher models.Teacher
		if err := db.First(&teacher, class.TeacherID).Error; err != nil {
			log.Printf("Failed to find teacher %d: %v", class.TeacherID, err)
			continue
		}

		// Auto-flag absent teacher; create system report for admin review if flag threshold exceeded
		if teacher.FlagCount < 2 {
			log.Printf("Automatically flagging teacher %d (current flags: %d)", teacher.ID, teacher.FlagCount)
			if err := AddTeacherFlag(db, teacher.ID, 1); err != nil {
				log.Printf("Error applying immediate flag to teacher %d: %v", teacher.ID, err)
			} else {
				desc := fmt.Sprintf("You have been automatically flagged for absence from your class session on %s.", session.ClassStart.Format(time.RFC822))
				CreateNotification(db, teacher.UserID, "system", desc)
			}
		} else {
			// A report must have a reporter; use the admin account for system-generated reports
			var firstAdmin models.Admin
			if err := db.Order("id asc").First(&firstAdmin).Error; err != nil {
				log.Printf("CRITICAL: Could not find an admin to assign system report. Skipping for session %d.", session.ID)
				continue
			}

			log.Printf("Absence would trigger ban for teacher %d. Creating system report on behalf of admin %d.", teacher.ID, firstAdmin.UserID)
			systemReport := models.Report{
				ReportUserID:      firstAdmin.UserID,
				ReportedUserID:    teacher.UserID,
				ClassSessionID:    session.ID,
				ReportType:        "learner",
				ReportReason:      "teacher_absent",
				ReportDescription: fmt.Sprintf("System detected teacher absence that requires admin review before banning. Teacher already has %d flags.", teacher.FlagCount),
				ReportStatus:      "pending",
			}
			if err := db.Create(&systemReport).Error; err != nil {
				log.Printf("Failed to create system report for session %d: %v", session.ID, err)
			} else {
				desc := fmt.Sprintf("Your absence from the class on %s has been flagged for admin review due to your current flag count.", session.ClassStart.Format(time.RFC822))
				CreateNotification(db, teacher.UserID, "system", desc)
			}
		}
		db.Model(&session).Update("class_status", "absent")
	}
}
