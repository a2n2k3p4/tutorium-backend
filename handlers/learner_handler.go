package handlers

import (
	"errors"

	"github.com/a2n2k3p4/tutorium-backend/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func LearnerRoutes(app *fiber.App) {
	app.Post("/learner", CreateLearner)
	app.Get("/learners", GetLearners)
	app.Get("/learner/:id", GetLearner)
	// app.Put("/learner/:id", UpdateLearner) No application logic for updating learner
	app.Delete("/learner/:id", DeleteLearner)
}

func CreateLearner(c *fiber.Ctx) error {
	var learner models.Learner

	if err := c.BodyParser(&learner); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if err := db.Create(&learner).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(200).JSON(learner)
}

func GetLearners(c *fiber.Ctx) error {
	learners := []models.Learner{}
	if err := db.Find(&learners).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(200).JSON(learners)
}

func findLearner(id int, learner *models.Learner) error {
	if err := db.First(learner, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("learner does not exist")
		}
		return err
	}
	return nil
}

func GetLearner(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var learner models.Learner

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	if err := findLearner(id, &learner); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	return c.Status(200).JSON(learner)
}

func DeleteLearner(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var learner models.Learner

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	if err = findLearner(id, &learner); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if err = db.Delete(&learner).Error; err != nil {
		return c.Status(404).JSON(err.Error())
	}
	return c.Status(200).JSON("Successfully deleted Learner")
}
