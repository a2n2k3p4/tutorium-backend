package handlers

import (
	"errors"

	"github.com/a2n2k3p4/tutorium-backend/middlewares"
	"github.com/a2n2k3p4/tutorium-backend/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ClassSessionRoutes(app *fiber.App) {
	classSession := app.Group("/class_sessions", middlewares.ProtectedMiddleware())
	classSession.Get("/", GetClassSessions)
	classSession.Get("/:id", GetClassSession)

	classSessionProtected := classSession.Group("/", middlewares.TeacherRequired())
	classSessionProtected.Post("/", CreateClassSession)
	classSessionProtected.Put("/:id", UpdateClassSession)
	classSessionProtected.Delete("/:id", DeleteClassSession)
}

// CreateClassSession godoc
//
//	@Summary		Create a new class session
//	@Description	CreateClassSession creates a new ClassSession record
//	@Tags			ClassSessions
//  @Security 		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			class_session	body		models.ClassSessionDoc	true	"ClassSession payload"
//	@Success		201				{object}	models.ClassSessionDoc
//	@Failure		400				{object}	map[string]string	"Invalid input"
//	@Failure		500				{object}	map[string]string	"Server error"
//	@Router			/class_sessions [post]
func CreateClassSession(c *fiber.Ctx) error {
	var class_session models.ClassSession

	if err := c.BodyParser(&class_session); err != nil {
		return c.Status(400).JSON(err.Error())
	}
	db, err := middlewares.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	if err := db.Create(&class_session).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(201).JSON(class_session)
}

// GetClassSessions godoc
//
//	@Summary		List all class sessions
//	@Description	GetClassSessions retrieves all ClassSession records with Class relation
//	@Tags			ClassSessions
//  @Security 		BearerAuth
//	@Produce		json
//	@Success		200	{array}		models.ClassSessionDoc
//	@Failure		500	{object}	map[string]string	"Server error"
//	@Router			/class_sessions [get]
func GetClassSessions(c *fiber.Ctx) error {
	class_sessions := []models.ClassSession{}
	db, err := middlewares.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	if err := db.Preload("Class").Find(&class_sessions).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(200).JSON(class_sessions)
}

func findClassSession(db *gorm.DB, id int, class_session *models.ClassSession) error {
	return db.Preload("Class").First(class_session, "id = ?", id).Error
}

// GetClassSession godoc
//
//	@Summary		Get class session by ID
//	@Description	GetClassSession retrieves a single ClassSession by its ID, including Class
//	@Tags			ClassSessions
//  @Security 		BearerAuth
//	@Produce		json
//	@Param			id	path		int	true	"ClassSession ID"
//	@Success		200	{object}	models.ClassSessionDoc
//	@Failure		400	{object}	map[string]string	"Invalid ID"
//	@Failure		404	{object}	map[string]string	"ClassSession not found"
//	@Failure		500	{object}	map[string]string	"Server error"
//	@Router			/class_sessions/{id} [get]
func GetClassSession(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var class_session models.ClassSession

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}
	db, err := middlewares.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	err = findClassSession(db, id, &class_session)
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return c.Status(404).JSON("class_session not found")
	case err != nil:
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(200).JSON(class_session)
}

// UpdateClassSession godoc
//
//	@Summary		Update an existing class session
//	@Description	UpdateClassSession updates a ClassSession record by its ID
//	@Tags			ClassSessions
//  @Security 		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id				path		int						true	"ClassSession ID"
//	@Param			class_session	body		models.ClassSessionDoc	true	"Updated ClassSession payload"
//	@Success		200				{object}	models.ClassSessionDoc
//	@Failure		400				{object}	map[string]string	"Invalid input"
//	@Failure		404				{object}	map[string]string	"ClassSession not found"
//	@Failure		500				{object}	map[string]string	"Server error"
//	@Router			/class_sessions/{id} [put]
func UpdateClassSession(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var class_session models.ClassSession

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}
	db, err := middlewares.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	err = findClassSession(db, id, &class_session)
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return c.Status(404).JSON("class_session not found")
	case err != nil:
		return c.Status(500).JSON(err.Error())
	}

	var class_session_update models.ClassSession
	if err := c.BodyParser(&class_session_update); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if err := db.Model(&class_session).Updates(class_session_update).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(200).JSON(class_session)

}

// DeleteClassSession godoc
//
//	@Summary		Delete a class session by ID
//	@Description	DeleteClassSession removes a ClassSession record by its ID
//	@Tags			ClassSessions
//  @Security 		BearerAuth
//	@Produce		json
//	@Param			id	path		int					true	"ClassSession ID"
//	@Success		200	{string}	string				"Successfully deleted class session"
//	@Failure		400	{object}	map[string]string	"Invalid ID"
//	@Failure		404	{object}	map[string]string	"ClassSession not found"
//	@Failure		500	{object}	map[string]string	"Server error"
//	@Router			/class_sessions/{id} [delete]
func DeleteClassSession(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var class_session models.ClassSession

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}
	db, err := middlewares.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	err = findClassSession(db, id, &class_session)
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return c.Status(404).JSON("class_session not found")
	case err != nil:
		return c.Status(500).JSON(err.Error())
	}

	if err = db.Delete(&class_session).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}
	return c.Status(200).JSON("Successfully deleted class session")
}
