package handlers

import (
	"errors"

	"github.com/a2n2k3p4/tutorium-backend/middlewares"
	"github.com/a2n2k3p4/tutorium-backend/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func ReviewRoutes(app *fiber.App) {
	review := app.Group("/reviews", middlewares.ProtectedMiddleware())
	review.Get("/", GetReviews)
	review.Get("/:id", GetReview)

	reviewLearner := review.Group("/", middlewares.LearnerRequired())
	reviewLearner.Post("/", CreateReview)
	reviewLearner.Put("/:id", UpdateReview)
	reviewLearner.Delete("/:id", DeleteReview)
}

// CreateReview godoc
//
//	@Summary		Create a new review
//	@Description	CreateReview creates a new Review record with rating validation
//	@Tags			Reviews
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			review	body		models.ReviewDoc	true	"Review payload"
//	@Success		201		{object}	models.ReviewDoc
//	@Failure		400		{object}	map[string]string	"Invalid input or rating out of range"
//	@Failure		500		{object}	map[string]string	"Server error"
//	@Router			/reviews [post]
func CreateReview(c *fiber.Ctx) error {
	var review models.Review

	if err := c.BodyParser(&review); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if review.Rating < 1 || review.Rating > 5 {
		return c.Status(400).JSON("Rating must be between 1 and 5")
	}
	db, err := middlewares.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(err.Error())
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
//	@Security		BearerAuth
//	@Produce		json
//	@Success		200	{array}		models.ReviewDoc
//	@Failure		500	{object}	map[string]string	"Server error"
//	@Router			/reviews [get]
func GetReviews(c *fiber.Ctx) error {
	var reviews []models.Review
	db, err := middlewares.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	if err := db.Preload("Learner").Preload("Class").Find(&reviews).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}
	return c.Status(200).JSON(reviews)
}

func findReview(db *gorm.DB, id int, review *models.Review) error {
	return db.Preload("Learner").Preload("Class").First(review, "id = ?", id).Error
}

// GetReview godoc
//
//	@Summary		Get review by ID
//	@Description	GetReview retrieves a single Review by its ID, including related Learner and Class
//	@Tags			Reviews
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id	path		int	true	"Review ID"
//	@Success		200	{object}	models.ReviewDoc
//	@Failure		400	{object}	map[string]string	"Invalid ID"
//	@Failure		404	{object}	map[string]string	"Review not found"
//	@Failure		500	{object}	map[string]string	"Server error"
//	@Router			/reviews/{id} [get]
func GetReview(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var review models.Review

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}
	db, err := middlewares.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	err = findReview(db, id, &review)
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
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int					true	"Review ID"
//	@Param			review	body		models.ReviewDoc	true	"Updated review payload"
//	@Success		200		{object}	models.ReviewDoc
//	@Failure		400		{object}	map[string]string	"Invalid input or rating out of range"
//	@Failure		404		{object}	map[string]string	"Review not found"
//	@Failure		500		{object}	map[string]string	"Server error"
//	@Router			/reviews/{id} [put]
func UpdateReview(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var review models.Review

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}
	db, err := middlewares.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	err = findReview(db, id, &review)
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

	if err := db.Model(&review).
		Omit(clause.Associations).
		Updates(review_updated).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(200).JSON(review)
}

// DeleteReview godoc
//
//	@Summary		Delete a review by ID
//	@Description	DeleteReview removes a Review record by its ID
//	@Tags			Reviews
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id	path		int					true	"Review ID"
//	@Success		200	{string}	string				"Successfully deleted review"
//	@Failure		400	{object}	map[string]string	"Invalid ID"
//	@Failure		404	{object}	map[string]string	"Review not found"
//	@Failure		500	{object}	map[string]string	"Server error"
//	@Router			/reviews/{id} [delete]
func DeleteReview(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var review models.Review

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}
	db, err := middlewares.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	err = findReview(db, id, &review)
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
