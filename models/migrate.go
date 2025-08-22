package models

import (
	"log"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&User{},
		&Admin{},
		&Learner{},
		&Teacher{},
		&BanDetailsLearner{},
		&BanDetailsTeacher{},
		&Class{},
		&ClassCategory{},
		&ClassSession{},
		&Enrollment{},
		&Notification{},
		&Report{},
		&Review{},
	)
	if err != nil {
		log.Fatalf("migration failed: %v", err)
	}
	log.Println("Database migrated successfully")
}
