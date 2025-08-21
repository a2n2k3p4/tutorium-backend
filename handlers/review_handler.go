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

// CreateReview godoc
//
//	@Summary		Create a new review
//	@Description	CreateReview creates a new Review record with rating validation
//	@Tags			Reviews
//	@Accept			json
//	@Produce		json
//	@Param			review	body		models.ReviewDoc	true	"Review payload"
//	@Success		201		{object}	models.ReviewDoc
//	@Failure		400		{object}	map[string]string	"Invalid input or rating out of range"
//	@Failure		500		{object}	map[string]string	"Server error"
//	@Router			/review [post]
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

// GetReviews godoc
//
//	@Summary		List all reviews
//	@Description	GetReviews retrieves all Review records with related Learner and Class
//	@Tags			Reviews
//	@Produce		json
//	@Success		200	{array}		models.ReviewDoc
//	@Failure		500	{object}	map[string]string	"Server error"
//	@Router			/reviews [get]
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

// GetReview godoc
//
//	@Summary		Get review by ID
//	@Description	GetReview retrieves a single Review by its ID, including related Learner and Class
//	@Tags			Reviews
//	@Produce		json
//	@Param			id	path		int	true	"Review ID"
//	@Success		200	{object}	models.ReviewDoc
//	@Failure		400	{object}	map[string]string	"Invalid ID"
//	@Failure		404	{object}	map[string]string	"Review not found"
//	@Failure		500	{object}	map[string]string	"Server error"
//	@Router			/review/{id} [get]
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

// UpdateReview godoc
//
//	@Summary		Update an existing review
//	@Description	UpdateReview updates a Review record by its ID; only Rating and Comment fields
//	@Tags			Reviews
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int					true	"Review ID"
//	@Param			review	body		models.ReviewDoc	true	"Updated review payload"
//	@Success		200		{object}	models.ReviewDoc
//	@Failure		400		{object}	map[string]string	"Invalid input or rating out of range"
//	@Failure		404		{object}	map[string]string	"Review not found"
//	@Failure		500		{object}	map[string]string	"Server error"
//	@Router			/review/{id} [put]
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

	type ReviewUpdate struct {
		Rating  *int
		Comment *string
	}

	var review_updated ReviewUpdate
	if err := c.BodyParser(&review_updated); err != nil {
		return c.Status(400).JSON(err.Error())
	}
	if review_updated.Rating != nil {
		if *review_updated.Rating < 1 || *review_updated.Rating > 5 {
			return c.Status(400).JSON("Rating must be between 1 and 5")
		}
	}

	if err := db.Model(&review).Updates(review_updated).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(200).JSON(review)
}

// DeleteReview godoc
//
//	@Summary		Delete a review by ID
//	@Description	DeleteReview removes a Review record by its ID
//	@Tags			Reviews
//	@Produce		json
//	@Param			id	path		int					true	"Review ID"
//	@Success		200	{string}	string				"Successfully deleted review"
//	@Failure		400	{object}	map[string]string	"Invalid ID"
//	@Failure		404	{object}	map[string]string	"Review not found"
//	@Failure		500	{object}	map[string]string	"Server error"
//	@Router			/review/{id} [delete]
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
