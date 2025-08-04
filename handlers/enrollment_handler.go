package handlers

import (
	"errors"

	"github.com/a2n2k3p4/tutorium-backend/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func EnrollmentRoutes(app *fiber.App) {
	app.Post("/enrollment", CreateEnrollment)
	app.Get("/enrollments", GetEnrollments)
	app.Get("/enrollment/:id", GetEnrollment)
	app.Put("/enrollment/:id", UpdateEnrollment)
	app.Delete("/enrollment/:id", DeleteEnrollment)
}

func CreateEnrollment(c *fiber.Ctx) error {
	var enrollment models.Enrollment

	if err := c.BodyParser(&enrollment); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if err := db.Create(&enrollment).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(201).JSON(enrollment)
}

func GetEnrollments(c *fiber.Ctx) error {
	enrollments := []models.Enrollment{}
	if err := db.Preload("Learner").Preload("Class").Find(&enrollments).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}
	return c.Status(200).JSON(enrollments)
}

func findEnrollment(id int, enrollment *models.Enrollment) error {
	return db.Preload("Learner").Preload("Class").First(enrollment, "id = ?", id).Error
}

func GetEnrollment(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var enrollment models.Enrollment

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	err = findEnrollment(id, &enrollment)
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return c.Status(404).JSON("enrollment not found")
	case err != nil:
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(200).JSON(enrollment)
}

func UpdateEnrollment(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var enrollment models.Enrollment

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	err = findEnrollment(id, &enrollment)
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return c.Status(404).JSON("enrollment not found")
	case err != nil:
		return c.Status(500).JSON(err.Error())
	}

	var enrollment_update models.Enrollment
	if err := c.BodyParser(&enrollment_update); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if err := db.Model(&enrollment).Updates(enrollment_update).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(200).JSON(enrollment)

}

func DeleteEnrollment(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var enrollment models.Enrollment

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	err = findEnrollment(id, &enrollment)
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return c.Status(404).JSON("enrollment not found")
	case err != nil:
		return c.Status(500).JSON(err.Error())
	}

	if err = db.Delete(&enrollment).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}
	return c.Status(200).JSON("Successfully deleted enrollment")
}
