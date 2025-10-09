package handlers

import (
	"errors"

	"github.com/a2n2k3p4/tutorium-backend/middlewares"
	"github.com/a2n2k3p4/tutorium-backend/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func TeacherRoutes(app *fiber.App) {
	teacher := app.Group("/teachers", middlewares.ProtectedMiddleware())
	teacher.Get("/", GetTeachers)
	teacher.Get("/:id", GetTeacher)
	teacher.Post("/", CreateTeacher)
	teacher.Put("/:id", UpdateTeacher)
	teacher.Delete("/:id", DeleteTeacher)
}

// CreateTeacher godoc
//
//	@Summary		Create a new teacher
//	@Description	CreateTeacher creates a new Teacher record
//	@Tags			Teachers
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			teacher	body		models.TeacherDoc	true	"Teacher payload"
//	@Success		201		{object}	models.TeacherDoc
//	@Failure		400		{string}	string	"Invalid input"
//	@Failure		500		{string}	string	"Server error"
//	@Router			/teachers [post]
func CreateTeacher(c *fiber.Ctx) error {
	var teacher models.Teacher

	if err := c.BodyParser(&teacher); err != nil {
		return c.Status(400).JSON(err.Error())
	}
	db, err := middlewares.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	if err := db.Create(&teacher).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(201).JSON(teacher)
}

// GetTeachers godoc
//
//	@Summary		List all teachers
//	@Description	GetTeachers retrieves all Teacher records
//	@Tags			Teachers
//	@Security		BearerAuth
//	@Produce		json
//	@Success		200	{array}		models.TeacherDoc
//	@Failure		500	{string}	string	"Server error"
//	@Router			/teachers [get]
func GetTeachers(c *fiber.Ctx) error {
	teachers := []models.Teacher{}
	db, err := middlewares.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	if err := db.Find(&teachers).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(200).JSON(teachers)
}

func findTeacher(db *gorm.DB, id int, teacher *models.Teacher) error {
	return db.First(teacher, "id = ?", id).Error
}

// GetTeacher godoc
//
//	@Summary		Get teacher by ID
//	@Description	GetTeacher retrieves a single Teacher by its ID
//	@Tags			Teachers
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id	path		int	true	"Teacher ID"
//	@Success		200	{object}	models.TeacherDoc
//	@Failure		400	{string}	string	"Invalid ID"
//	@Failure		404	{string}	string	"Teacher not found"
//	@Failure		500	{string}	string	"Server error"
//	@Router			/teachers/{id} [get]
func GetTeacher(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var teacher models.Teacher

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}
	db, err := middlewares.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	err = findTeacher(db, id, &teacher)
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return c.Status(404).JSON("teacher not found")
	case err != nil:
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(200).JSON(teacher)
}

// UpdateTeacher godoc
//
//	@Summary		Update an existing teacher
//	@Description	UpdateTeacher updates a Teacher record by its ID
//	@Tags			Teachers
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int					true	"Teacher ID"
//	@Param			teacher	body		models.TeacherDoc	true	"Updated teacher payload"
//	@Success		200		{object}	models.TeacherDoc
//	@Failure		400		{string}	string	"Invalid input"
//	@Failure		404		{string}	string	"Teacher not found"
//	@Failure		500		{string}	string	"Server error"
//	@Router			/teachers/{id} [put]
func UpdateTeacher(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var teacher models.Teacher

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}
	db, err := middlewares.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	err = findTeacher(db, id, &teacher)
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return c.Status(404).JSON("teacher not found")
	case err != nil:
		return c.Status(500).JSON(err.Error())
	}

	var teacher_update models.Teacher
	if err := c.BodyParser(&teacher_update); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if err := db.Model(&teacher).Updates(teacher_update).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(200).JSON(teacher)
}

// DeleteTeacher godoc
//
//	@Summary		Delete a teacher by ID
//	@Description	DeleteTeacher removes a Teacher record by its ID
//	@Tags			Teachers
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id	path		int					true	"Teacher ID"
//	@Success		200	{string}	string				"Successfully deleted Teacher"
//	@Failure		400	{string}	string	"Invalid ID"
//	@Failure		404	{string}	string	"Teacher not found"
//	@Failure		500	{string}	string	"Server error"
//	@Router			/teachers/{id} [delete]
func DeleteTeacher(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var teacher models.Teacher

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}
	db, err := middlewares.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	err = findTeacher(db, id, &teacher)
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return c.Status(404).JSON("teacher not found")
	case err != nil:
		return c.Status(500).JSON(err.Error())
	}

	if err = db.Delete(&teacher).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}
	return c.Status(200).JSON("Successfully deleted Teacher")
}
