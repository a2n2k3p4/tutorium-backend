package handlers

import (
	"errors"
	"strings"

	"github.com/a2n2k3p4/tutorium-backend/middlewares"
	"github.com/a2n2k3p4/tutorium-backend/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func EnrollmentRoutes(app *fiber.App) {
	enrollment := app.Group("/enrollments", middlewares.ProtectedMiddleware(), middlewares.BanMiddleware(), middlewares.LearnerRequired())

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
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			enrollment	body		models.EnrollmentDoc	true	"Enrollment payload"
//	@Success		201			{object}	models.EnrollmentDoc
//	@Failure		400			{string}	string	"Invalid input"
//	@Failure		500			{string}	string	"Server error"
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
//	@Security		BearerAuth
//	@Produce		json
//	@Param			session_ids	query	[]string	false	"Filter by one or more class session IDs (comma-separated or repeated query param)"
//	@Param			include		query	[]string	false	"Include related entities: learner, class_session, user (comma-separated or repeated query param)"
//	@Success		200	{array}		models.EnrollmentDoc
//	@Failure		500	{string}	string	"Server error"
//	@Router			/enrollments [get]
func GetEnrollments(c *fiber.Ctx) error {
	type EnrollmentQueryParams struct {
		SessionIDs []string `query:"session_ids"`
		Include    []string `query:"include"`
	}

	var params EnrollmentQueryParams

	if err := c.QueryParser(&params); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	// This splits a single string "1,3" into multiple values ["1", "3"]
	if len(params.SessionIDs) == 1 && strings.Contains(params.SessionIDs[0], ",") {
		params.SessionIDs = strings.Split(params.SessionIDs[0], ",")
		for i := range params.SessionIDs {
			params.SessionIDs[i] = strings.TrimSpace(params.SessionIDs[i])
		}
	}

	// This splits a single string "learner,user" into multiple values ["learner", "user"]
	if len(params.Include) == 1 && strings.Contains(params.Include[0], ",") {
		params.Include = strings.Split(params.Include[0], ",")
		for i := range params.Include {
			params.Include[i] = strings.TrimSpace(params.Include[i])
		}
	}

	enrollments := []models.Enrollment{}
	db, err := middlewares.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	query := db.Model(&models.Enrollment{})

	// include logic
	includeMap := make(map[string]bool)
	for _, inc := range params.Include {
		includeMap[inc] = true
	}

	if includeMap["learner"] {
		query = query.Preload("Learner")
	}

	if includeMap["class_session"] {
		query = query.Preload("ClassSession")
	}

	// Only include enrollments for the given session IDs
	if len(params.SessionIDs) > 0 {
		query = query.Where("class_session_id IN (?)", params.SessionIDs)
	}

	if err := query.Find(&enrollments).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	responses := make([]models.EnrollmentResponse, len(enrollments))

	if includeMap["user"] {
		// collect all user IDs from enrollments
		userIDs := make([]uint, len(enrollments))
		for i, e := range enrollments {
			userIDs[i] = e.Learner.UserID
		}

		users := []models.User{}
		if err := db.Where("id IN (?)", userIDs).Find(&users).Error; err != nil {
			return c.Status(500).JSON(err.Error())
		}

		userMap := make(map[uint]*models.User, len(users))
		for i := range users {
			userMap[users[i].ID] = &users[i]
		}

		for i, e := range enrollments {
			responses[i].Enrollment = e
			responses[i].User = userMap[e.Learner.UserID]
		}
	} else {
		for i, e := range enrollments {
			responses[i].Enrollment = e
		}
	}

	return c.Status(200).JSON(responses)
}

func findEnrollment(db *gorm.DB, id int, enrollment *models.Enrollment) error {
	return db.Preload("Learner").Preload("ClassSession").First(enrollment, "id = ?", id).Error
}

// GetEnrollment godoc
//
//	@Summary		Get enrollment by ID
//	@Description	GetEnrollment retrieves a single Enrollment by its ID, including related Learner and Class
//	@Tags			Enrollments
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id	path		int	true	"Enrollment ID"
//	@Success		200	{object}	models.EnrollmentDoc
//	@Failure		400	{string}	string	"Invalid ID"
//	@Failure		404	{string}	string	"Enrollment not found"
//	@Failure		500	{string}	string	"Server error"
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
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id			path		int						true	"Enrollment ID"
//	@Param			enrollment	body		models.EnrollmentDoc	true	"Updated enrollment payload"
//	@Success		200			{object}	models.EnrollmentDoc
//	@Failure		400			{string}	string	"Invalid input"
//	@Failure		404			{string}	string	"Enrollment not found"
//	@Failure		500			{string}	string	"Server error"
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

	if err := db.Model(&enrollment).
		Omit(clause.Associations).
		Updates(enrollment_update).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(200).JSON(enrollment)

}

// DeleteEnrollment godoc
//
//	@Summary		Delete an enrollment by ID
//	@Description	DeleteEnrollment removes an Enrollment record by its ID
//	@Tags			Enrollments
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id	path		int		true	"Enrollment ID"
//	@Success		200	{string}	string	"Successfully deleted enrollment"
//	@Failure		400	{string}	string	"Invalid ID"
//	@Failure		404	{string}	string	"Enrollment not found"
//	@Failure		500	{string}	string	"Server error"
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
