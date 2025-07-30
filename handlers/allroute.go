package handlers

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var db *gorm.DB

func AllRoutes(database *gorm.DB, app *fiber.App) {
	db = database

	AdminRoutes(database, app)
	BanLearnerRoutes(database, app)
	BanTeacherRoutes(database, app)
	ClassCategoryRoutes(database, app)
	ClassRoutes(database, app)
	ClassSessionRoutes(database, app)
	EnrollmentRoutes(database, app)
}
