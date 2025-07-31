package handlers

import (
	"errors"

	"github.com/a2n2k3p4/tutorium-backend/models"
	"github.com/gofiber/fiber/v2"
)

func BanLearnerRoutes(app *fiber.App) {
	app.Post("/banlearner", CreateBanLearner)
	app.Get("/banlearners", GetBanLearners)
	app.Get("/banlearner/:id", GetBanLearner)
	app.Put("/banlearner/:id", UpdateBanLearner)
	app.Delete("/banlearner/:id", DeleteBanLearner)
}

func CreateBanLearner(c *fiber.Ctx) error {
	var banlearner models.BanDetailsLearner

	if err := c.BodyParser(&banlearner); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	db.Create(&banlearner)
	return c.Status(200).JSON(banlearner)
}

func GetBanLearners(c *fiber.Ctx) error {
	banlearners := []models.BanDetailsLearner{}
	db.Find(&banlearners)

	return c.Status(200).JSON(banlearners)
}

func findbanlearner(id int, banlearner *models.BanDetailsLearner) error {
	db.Find(&banlearner, "id = ?", id)
	if banlearner.ID == 0 {
		return errors.New("ban learner does not exist")
	}
	return nil
}

func GetBanLearner(c *fiber.Ctx) error {
	id, err := c.ParamsInt("ID")

	var banlearner models.BanDetailsLearner

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	if err := findbanlearner(id, &banlearner); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	return c.Status(200).JSON(banlearner)
}

func UpdateBanLearner(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var banlearner models.BanDetailsLearner

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	err = findbanlearner(id, &banlearner)

	if err != nil {
		return c.Status(400).JSON(err.Error())
	}

	var banlearner_update models.BanDetailsLearner
	if err := c.BodyParser(&banlearner_update); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	db.Model(&banlearner).Updates(banlearner_update)

	return c.Status(200).JSON(banlearner)

}

func DeleteBanLearner(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var banlearner models.BanDetailsLearner

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	err = findbanlearner(id, &banlearner)

	if err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if err = db.Delete(&banlearner).Error; err != nil {
		return c.Status(404).JSON(err.Error())
	}
	return c.Status(200).JSON("Successfully deleted ban learner")
}
