package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func AllRoutes(app *fiber.App) {
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
	PaymentRoutes(app)
}
