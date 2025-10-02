package handlers

import (
	"errors"

	"github.com/a2n2k3p4/tutorium-backend/middlewares"
	"github.com/a2n2k3p4/tutorium-backend/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func LearnerRoutes(app *fiber.App) {
	learner := app.Group("/learners", middlewares.ProtectedMiddleware())

	learner.Post("/", CreateLearner)
	learner.Get("/", GetLearners)
	learner.Get("/:id", GetLearner)
	// learner.Put("/:id", UpdateLearner) No application logic for updating learner
	learner.Delete("/:id", DeleteLearner)
}

// CreateLearner godoc
//
//	@Summary		Create a new learner
//	@Description	CreateLearner creates a new Learner record
//	@Tags			Learners
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			learner	body		models.LearnerDoc	true	"Learner payload"
//	@Success		201		{object}	models.LearnerDoc
//	@Failure		400		{string}	string	"Invalid input"
//	@Failure		500		{string}	string	"Server error"
//	@Router			/learners [post]
func CreateLearner(c *fiber.Ctx) error {
	var learner models.Learner

	if err := c.BodyParser(&learner); err != nil {
		return c.Status(400).JSON(err.Error())
	}
	db, err := middlewares.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	if err := db.Create(&learner).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(201).JSON(learner)
}

// GetLearners godoc
//
//	@Summary		List all learners
//	@Description	GetLearners retrieves all Learner records
//	@Tags			Learners
//	@Security		BearerAuth
//	@Produce		json
//	@Success		200	{array}		models.LearnerDoc
//	@Failure		500	{string}	string	"Server error"
//	@Router			/learners [get]
func GetLearners(c *fiber.Ctx) error {
	learners := []models.Learner{}
	db, err := middlewares.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	if err := db.Find(&learners).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(200).JSON(learners)
}

func findLearner(db *gorm.DB, id int, learner *models.Learner) error {
	return db.First(learner, "id = ?", id).Error
}

// GetLearner godoc
//
//	@Summary		Get learner by ID
//	@Description	GetLearner retrieves a single Learner by its ID
//	@Tags			Learners
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id	path		int	true	"Learner ID"
//	@Success		200	{object}	models.LearnerDoc
//	@Failure		400	{string}	string	"Invalid ID"
//	@Failure		404	{string}	string	"Learner not found"
//	@Failure		500	{string}	string	"Server error"
//	@Router			/learners/{id} [get]
func GetLearner(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var learner models.Learner

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}
	db, err := middlewares.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	err = findLearner(db, id, &learner)
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return c.Status(404).JSON("Learner not found")
	case err != nil:
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(200).JSON(learner)
}

// DeleteLearner godoc
//
//	@Summary		Delete a learner by ID
//	@Description	DeleteLearner removes a Learner record by its ID
//	@Tags			Learners
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id	path		int		true	"Learner ID"
//	@Success		200	{string}	string	"Successfully deleted Learner"
//	@Failure		400	{string}	string	"Invalid ID"
//	@Failure		404	{string}	string	"Learner not found"
//	@Failure		500	{string}	string	"Server error"
//	@Router			/learners/{id} [delete]
func DeleteLearner(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var learner models.Learner

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}
	db, err := middlewares.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	err = findLearner(db, id, &learner)
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return c.Status(404).JSON("Learner not found")
	case err != nil:
		return c.Status(500).JSON(err.Error())
	}

	if err = db.Delete(&learner).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}
	return c.Status(200).JSON("Successfully deleted Learner")
}
