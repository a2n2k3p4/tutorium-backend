package handlers

import (
	"errors"

	"github.com/a2n2k3p4/tutorium-backend/middlewares"
	"github.com/a2n2k3p4/tutorium-backend/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ClassRoutes(app *fiber.App) {
	class := app.Group("/classes", middlewares.ProtectedMiddleware())
	class.Get("/", GetClasses)
	class.Get("/:id", GetClass)

	classProtected := class.Group("/", middlewares.TeacherRequired())
	classProtected.Post("/", CreateClass)
	classProtected.Put("/:id", UpdateClass)
	classProtected.Delete("/:id", DeleteClass)
}

// CreateClass godoc
//
//		@Summary		Create a new class
//		@Description	CreateClass creates a new Class record
//		@Tags			Classes
//	 @Security 		BearerAuth
//		@Accept			json
//		@Produce		json
//		@Param			class	body		models.ClassDoc	true	"Class payload"
//		@Success		201		{object}	models.ClassDoc
//		@Failure		400		{object}	map[string]string	"Invalid input"
//		@Failure		500		{object}	map[string]string	"Server error"
//		@Router			/classes [post]
func CreateClass(c *fiber.Ctx) error {
	var class models.Class

	if err := c.BodyParser(&class); err != nil {
		return c.Status(400).JSON(err.Error())
	}
	db, err := middlewares.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	if err := db.Create(&class).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}
	return c.Status(201).JSON(class)
}

// GetClasses godoc
//
//		@Summary		List all classes
//		@Description	GetClasses retrieves all Class records with Teacher and Categories relations
//		@Tags			Classes
//	 @Security 		BearerAuth
//		@Produce		json
//		@Success		200	{array}		models.ClassDoc
//		@Failure		500	{object}	map[string]string	"Server error"
//		@Router			/classes [get]
func GetClasses(c *fiber.Ctx) error {
	classes := []models.Class{}
	db, err := middlewares.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	if err := db.Preload("Teacher").Preload("Categories").Find(&classes).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(200).JSON(classes)
}

func findClass(db *gorm.DB, id int, class *models.Class) error {
	return db.Preload("Teacher").Preload("Categories").First(class, "id = ?", id).Error
}

// GetClass godoc
//
//		@Summary		Get class by ID
//		@Description	GetClass retrieves a single Class by its ID, including Teacher and Categories
//		@Tags			Classes
//	 @Security 		BearerAuth
//		@Produce		json
//		@Param			id	path		int	true	"Class ID"
//		@Success		200	{object}	models.ClassDoc
//		@Failure		400	{object}	map[string]string	"Invalid ID"
//		@Failure		404	{object}	map[string]string	"Class not found"
//		@Failure		500	{object}	map[string]string	"Server error"
//		@Router			/classes/{id} [get]
func GetClass(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var class models.Class

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}
	db, err := middlewares.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	err = findClass(db, id, &class)
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return c.Status(404).JSON("class not found")
	case err != nil:
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(200).JSON(class)
}

// UpdateClass godoc
//
//		@Summary		Update an existing class
//		@Description	UpdateClass updates a Class record by its ID
//		@Tags			Classes
//	 @Security 		BearerAuth
//		@Accept			json
//		@Produce		json
//		@Param			id		path		int				true	"Class ID"
//		@Param			class	body		models.ClassDoc	true	"Updated class payload"
//		@Success		200		{object}	models.ClassDoc
//		@Failure		400		{object}	map[string]string	"Invalid input"
//		@Failure		404		{object}	map[string]string	"Class not found"
//		@Failure		500		{object}	map[string]string	"Server error"
//		@Router			/classes/{id} [put]
func UpdateClass(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var class models.Class

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}
	db, err := middlewares.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	err = findClass(db, id, &class)
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return c.Status(404).JSON("class not found")
	case err != nil:
		return c.Status(500).JSON(err.Error())
	}

	var class_update models.Class
	if err := c.BodyParser(&class_update); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if err := db.Model(&class).Updates(class_update).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(200).JSON(class)

}

// DeleteClass godoc
//
//		@Summary		Delete a class by ID
//		@Description	DeleteClass removes a Class record by its ID
//		@Tags			Classes
//	 @Security 		BearerAuth
//		@Produce		json
//		@Param			id	path		int					true	"Class ID"
//		@Success		200	{string}	string				"Successfully deleted class"
//		@Failure		400	{object}	map[string]string	"Invalid ID"
//		@Failure		404	{object}	map[string]string	"Class not found"
//		@Failure		500	{object}	map[string]string	"Server error"
//		@Router			/classes/{id} [delete]
func DeleteClass(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var class models.Class

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}
	db, err := middlewares.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	err = findClass(db, id, &class)
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return c.Status(404).JSON("class not found")
	case err != nil:
		return c.Status(500).JSON(err.Error())
	}

	if err = db.Delete(&class).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}
	return c.Status(200).JSON("Successfully deleted class")
}
