package handlers

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var db *gorm.DB

func AllRoutes(database *gorm.DB, app *fiber.App) {
	db = database

	AdminRoutes(app)
	BanLearnerRoutes(app)
	BanTeacherRoutes(app)
	ClassCategoryRoutes(app)
	ClassRoutes(app)
	ClassSessionRoutes(app)
	EnrollmentRoutes(app)
}
