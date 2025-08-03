package handlers

import (
	"errors"

	"github.com/a2n2k3p4/tutorium-backend/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ReviewRoutes(app *fiber.App) {
	app.Post("/review", CreateReview)
	app.Get("/reviews", GetReviews)
	app.Get("/review/:id", GetReview)
	app.Put("/review/:id", UpdateReview)
	app.Delete("/review/:id", DeleteReview)
}

func CreateReview(c *fiber.Ctx) error {
	var review models.Review

	if err := c.BodyParser(&review); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if review.Rating < 1 || review.Rating > 5 {
		return c.Status(400).JSON("Rating must be between 1 and 5")
	}

	if err := db.Create(&review).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(201).JSON(review)
}

func GetReviews(c *fiber.Ctx) error {
	var reviews []models.Review
	if err := db.Preload("Learner").Preload("Class").Find(&reviews).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}
	return c.Status(200).JSON(reviews)
}

func findReview(id int, review *models.Review) error {
	return db.Preload("Learner").Preload("Class").First(review, "id = ?", id).Error
}

func GetReview(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var review models.Review

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	err = findReview(id, &review)
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return c.Status(404).JSON("review not found")
	case err != nil:
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(200).JSON(review)
}

func UpdateReview(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var review models.Review

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	err = findReview(id, &review)
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return c.Status(404).JSON("review not found")
	case err != nil:
		return c.Status(500).JSON(err.Error())
	}

	var review_updated models.Review
	if err := c.BodyParser(&review_updated); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if review_updated.Rating < 1 || review_updated.Rating > 5 {
		return c.Status(400).JSON("Rating must be between 1 and 5")
	}

	if err := db.Model(&review).Updates(review_updated).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(200).JSON(review)
}

func DeleteReview(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var review models.Review

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	err = findReview(id, &review)
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return c.Status(404).JSON("review not found")
	case err != nil:
		return c.Status(500).JSON(err.Error())
	}

	if err := db.Delete(&review).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(200).JSON("Successfully deleted review")
}
