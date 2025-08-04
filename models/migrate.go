package models

import (
	"log"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&Admin{},
		&BanDetailsLearner{},
		&BanDetailsTeacher{},
		&Class{},
		&ClassCategory{},
		&ClassSession{},
		&Enrollment{},
		&Learner{},
		&Notification{},
		&Report{},
		&Review{},
		&Teacher{},
		&User{},
	)
	if err != nil {
		log.Fatalf("migration failed: %v", err)
	}
	log.Println("Database migrated successfully")
}
