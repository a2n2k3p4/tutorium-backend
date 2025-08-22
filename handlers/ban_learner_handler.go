package handlers

import (
	"errors"

	"github.com/a2n2k3p4/tutorium-backend/middleware"
	"github.com/a2n2k3p4/tutorium-backend/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func BanLearnerRoutes(app *fiber.App) {
	banLearner := app.Group("/banlearners", middleware.ProtectedMiddleware(), middleware.AdminRequired())

	banLearner.Post("/", CreateBanLearner)
	banLearner.Get("/", GetBanLearners)
	banLearner.Get("/:id", GetBanLearner)
	banLearner.Put("/:id", UpdateBanLearner)
	banLearner.Delete("/:id", DeleteBanLearner)
}

// CreateBanLearner godoc
//
//	@Summary		Create a new banned learner record
//	@Description	CreateBanLearner creates a new ban record for a learner
//	@Tags			BanLearners
//	@Accept			json
//	@Produce		json
//	@Param			banlearner	body		models.BanDetailsLearnerDoc	true	"Ban Learner Payload"
//	@Success		201			{object}	models.BanDetailsLearnerDoc
//	@Failure		400			{object}	map[string]string
//	@Failure		500			{object}	map[string]string
//	@Router			/banlearners [post]
func CreateBanLearner(c *fiber.Ctx) error {
	var banlearner models.BanDetailsLearner

	if err := c.BodyParser(&banlearner); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if err := db.Create(&banlearner).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}
	return c.Status(201).JSON(banlearner)
}

// GetBanLearners godoc
//
//	@Summary		Get all banned learners
//	@Description	GetBanLearners returns a list of all ban records
//	@Tags			BanLearners
//	@Produce		json
//	@Success		200	{array}		models.BanDetailsLearnerDoc
//	@Failure		500	{object}	map[string]string
//	@Router			/banlearners [get]
func GetBanLearners(c *fiber.Ctx) error {
	banlearners := []models.BanDetailsLearner{}
	if err := db.Preload("Learner").Find(&banlearners).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(200).JSON(banlearners)
}

func findBanLearner(id int, banlearner *models.BanDetailsLearner) error {
	return db.Preload("Learner").First(banlearner, "id = ?", id).Error
}

// GetBanLearner godoc
//
//	@Summary		Get a banned learner by ID
//	@Description	GetBanLearner returns a single ban record by ID
//	@Tags			BanLearners
//	@Produce		json
//	@Param			id	path		int	true	"Ban Learner ID"
//	@Success		200	{object}	models.BanDetailsLearnerDoc
//	@Failure		400	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/banlearners/{id} [get]
func GetBanLearner(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var banlearner models.BanDetailsLearner

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	err = findBanLearner(id, &banlearner)
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return c.Status(404).JSON("banlearner not found")
	case err != nil:
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(200).JSON(banlearner)
}

// UpdateBanLearner godoc
//
//	@Summary		Update a banned learner record
//	@Description	UpdateBanLearner modifies an existing ban record
//	@Tags			BanLearners
//	@Accept			json
//	@Produce		json
//	@Param			id			path		int							true	"Ban Learner ID"
//	@Param			banlearner	body		models.BanDetailsLearnerDoc	true	"Updated ban record"
//	@Success		200			{object}	models.BanDetailsLearnerDoc
//	@Failure		400			{object}	map[string]string
//	@Failure		404			{object}	map[string]string
//	@Failure		500			{object}	map[string]string
//	@Router			/banlearners/{id} [put]
func UpdateBanLearner(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var banlearner models.BanDetailsLearner

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	err = findBanLearner(id, &banlearner)
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return c.Status(404).JSON("banlearner not found")
	case err != nil:
		return c.Status(500).JSON(err.Error())
	}

	var banlearner_update models.BanDetailsLearner
	if err := c.BodyParser(&banlearner_update); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if err := db.Model(&banlearner).Updates(banlearner_update).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(200).JSON(banlearner)

}

// DeleteBanLearner godoc
//
//	@Summary		Delete a banned learner record
//	@Description	DeleteBanLearner deletes a ban record by ID
//	@Tags			BanLearners
//	@Produce		json
//	@Param			id	path		int		true	"Ban Learner ID"
//	@Success		200	{string}	string	"Successfully deleted ban learner"
//	@Failure		400	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/banlearners/{id} [delete]
func DeleteBanLearner(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var banlearner models.BanDetailsLearner

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	err = findBanLearner(id, &banlearner)
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return c.Status(404).JSON("banlearner not found")
	case err != nil:
		return c.Status(500).JSON(err.Error())
	}

	if err = db.Delete(&banlearner).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}
	return c.Status(200).JSON("Successfully deleted ban learner")
}
