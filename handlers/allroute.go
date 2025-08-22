package handlers

import (
	"github.com/a2n2k3p4/tutorium-backend/config/dbserver"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var db *gorm.DB

func AllRoutes(app *fiber.App) {
	db = dbserver.DB

	AdminRoutes(app)
	BanLearnerRoutes(app)
	BanTeacherRoutes(app)
	ClassCategoryRoutes(app)
	ClassRoutes(app)
	ClassSessionRoutes(app)
	EnrollmentRoutes(app)
	LearnerRoutes(app)
	NotificationRoutes(app)
	ReportRoutes(app)
	ReviewRoutes(app)
	TeacherRoutes(app)
	UserRoutes(app)
	LoginRoutes(app)
}
