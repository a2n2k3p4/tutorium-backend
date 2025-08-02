package handlers

import (
	"errors"

	"github.com/a2n2k3p4/tutorium-backend/models"
	"github.com/gofiber/fiber/v2"
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

	db.Create(&enrollment)
	return c.Status(200).JSON(enrollment)
}

func GetEnrollments(c *fiber.Ctx) error {
	enrollments := []models.Enrollment{}
	dbErr := db.Find(&enrollments).Error
	if dbErr != nil {
		return c.Status(404).JSON(dbErr)
	}

	return c.Status(200).JSON(enrollments)
}

func findenrollment(id int, enrollment *models.Enrollment) error {
	dbErr := db.Find(&enrollment, "id = ?", id).Error
	if dbErr != nil {
		return errors.New(dbErr.Error())
	}

	if enrollment.ID == 0 {
		return errors.New("enrollment does not exist")
	}
	return nil
}

func GetEnrollment(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var enrollment models.Enrollment

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	if err := findenrollment(id, &enrollment); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	return c.Status(200).JSON(enrollment)
}

func UpdateEnrollment(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var enrollment models.Enrollment

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	err = findenrollment(id, &enrollment)

	if err != nil {
		return c.Status(400).JSON(err.Error())
	}

	var enrollment_update models.Enrollment
	if err := c.BodyParser(&enrollment_update); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	db.Model(&enrollment).Updates(enrollment_update)

	return c.Status(200).JSON(enrollment)

}

func DeleteEnrollment(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var enrollment models.Enrollment

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	err = findenrollment(id, &enrollment)

	if err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if err = db.Delete(&enrollment).Error; err != nil {
		return c.Status(404).JSON(err.Error())
	}
	return c.Status(200).JSON("Successfully deleted enrollment")
}
