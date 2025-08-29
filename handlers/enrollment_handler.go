package handlers

import (
	"errors"

	"github.com/a2n2k3p4/tutorium-backend/middlewares"
	"github.com/a2n2k3p4/tutorium-backend/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func EnrollmentRoutes(app *fiber.App) {
	enrollment := app.Group("/enrollments", middlewares.ProtectedMiddleware(), middlewares.LearnerRequired())

	enrollment.Post("/", CreateEnrollment)
	enrollment.Get("/", GetEnrollments)
	enrollment.Get("/:id", GetEnrollment)
	enrollment.Put("/:id", UpdateEnrollment)
	enrollment.Delete("/:id", DeleteEnrollment)
}

// CreateEnrollment godoc
//
//	@Summary		Create a new enrollment
//	@Description	CreateEnrollment creates a new Enrollment record
//	@Tags			Enrollments
//  @Security 		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			enrollment	body		models.EnrollmentDoc	true	"Enrollment payload"
//	@Success		201			{object}	models.EnrollmentDoc
//	@Failure		400			{object}	map[string]string	"Invalid input"
//	@Failure		500			{object}	map[string]string	"Server error"
//	@Router			/enrollments [post]
func CreateEnrollment(c *fiber.Ctx) error {
	var enrollment models.Enrollment

	if err := c.BodyParser(&enrollment); err != nil {
		return c.Status(400).JSON(err.Error())
	}
	db, err := middlewares.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	if err := db.Create(&enrollment).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(201).JSON(enrollment)
}

// GetEnrollments godoc
//
//	@Summary		List all enrollments
//	@Description	GetEnrollments retrieves all Enrollment records with associated Learner and Class
//	@Tags			Enrollments
//  @Security 		BearerAuth
//	@Produce		json
//	@Success		200	{array}		models.EnrollmentDoc
//	@Failure		500	{object}	map[string]string	"Server error"
//	@Router			/enrollments [get]
func GetEnrollments(c *fiber.Ctx) error {
	enrollments := []models.Enrollment{}
	db, err := middlewares.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	if err := db.Preload("Learner").Preload("Class").Find(&enrollments).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}
	return c.Status(200).JSON(enrollments)
}

func findEnrollment(db *gorm.DB, id int, enrollment *models.Enrollment) error {
	return db.Preload("Learner").Preload("Class").First(enrollment, "id = ?", id).Error
}

// GetEnrollment godoc
//
//	@Summary		Get enrollment by ID
//	@Description	GetEnrollment retrieves a single Enrollment by its ID, including related Learner and Class
//	@Tags			Enrollments
//  @Security 		BearerAuth
//	@Produce		json
//	@Param			id	path		int	true	"Enrollment ID"
//	@Success		200	{object}	models.EnrollmentDoc
//	@Failure		400	{object}	map[string]string	"Invalid ID"
//	@Failure		404	{object}	map[string]string	"Enrollment not found"
//	@Failure		500	{object}	map[string]string	"Server error"
//	@Router			/enrollments/{id} [get]
func GetEnrollment(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var enrollment models.Enrollment

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}
	db, err := middlewares.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	err = findEnrollment(db, id, &enrollment)
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return c.Status(404).JSON("enrollment not found")
	case err != nil:
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(200).JSON(enrollment)
}

// UpdateEnrollment godoc
//
//	@Summary		Update an existing enrollment
//	@Description	UpdateEnrollment updates an Enrollment record by its ID
//	@Tags			Enrollments
//  @Security 		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id			path		int						true	"Enrollment ID"
//	@Param			enrollment	body		models.EnrollmentDoc	true	"Updated enrollment payload"
//	@Success		200			{object}	models.EnrollmentDoc
//	@Failure		400			{object}	map[string]string	"Invalid input"
//	@Failure		404			{object}	map[string]string	"Enrollment not found"
//	@Failure		500			{object}	map[string]string	"Server error"
//	@Router			/enrollments/{id} [put]
func UpdateEnrollment(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var enrollment models.Enrollment

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}
	db, err := middlewares.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	err = findEnrollment(db, id, &enrollment)
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

// DeleteEnrollment godoc
//
//	@Summary		Delete an enrollment by ID
//	@Description	DeleteEnrollment removes an Enrollment record by its ID
//	@Tags			Enrollments
//  @Security 		BearerAuth
//	@Produce		json
//	@Param			id	path		int					true	"Enrollment ID"
//	@Success		200	{string}	string				"Successfully deleted enrollment"
//	@Failure		400	{object}	map[string]string	"Invalid ID"
//	@Failure		404	{object}	map[string]string	"Enrollment not found"
//	@Failure		500	{object}	map[string]string	"Server error"
//	@Router			/enrollments/{id} [delete]
func DeleteEnrollment(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var enrollment models.Enrollment

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}
	db, err := middlewares.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	err = findEnrollment(db, id, &enrollment)
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
