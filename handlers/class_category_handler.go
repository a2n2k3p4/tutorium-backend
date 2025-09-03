package handlers

import (
	"errors"

	"github.com/a2n2k3p4/tutorium-backend/middlewares"
	"github.com/a2n2k3p4/tutorium-backend/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ClassCategoryRoutes(app *fiber.App) {
	classCategory := app.Group("/class_categories", middlewares.ProtectedMiddleware())
	classCategory.Get("/", GetClassCategories)
	classCategory.Get("/:id", GetClassCategory)

	classCategoryProtected := classCategory.Group("/", middlewares.AdminRequired())
	classCategoryProtected.Post("/", CreateClassCategory)
	classCategoryProtected.Put("/:id", UpdateClassCategory)
	classCategoryProtected.Delete("/:id", DeleteClassCategory)
}

// CreateClassCategory godoc
//
//		@Summary		Create a new class category
//		@Description	CreateClassCategory creates a new ClassCategory record
//		@Tags			ClassCategories
//	 @Security 		BearerAuth
//		@Accept			json
//		@Produce		json
//		@Param			class_category	body		models.ClassCategoryDoc	true	"ClassCategory payload"
//		@Success		201				{object}	models.ClassCategoryDoc
//		@Failure		400				{object}	map[string]string	"Invalid input"
//		@Failure		500				{object}	map[string]string	"Server error"
//		@Router			/class_categories [post]
func CreateClassCategory(c *fiber.Ctx) error {
	var class_category models.ClassCategory

	if err := c.BodyParser(&class_category); err != nil {
		return c.Status(400).JSON(err.Error())
	}
	db, err := middlewares.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	if err := db.Create(&class_category).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(201).JSON(class_category)
}

// GetClassCategories godoc
//
//		@Summary		List all class categories
//		@Description	GetClassCategories retrieves all ClassCategory records with Classes relation
//		@Tags			ClassCategories
//	 @Security 		BearerAuth
//		@Produce		json
//		@Success		200	{array}		models.ClassCategoryDoc
//		@Failure		500	{object}	map[string]string	"Server error"
//		@Router			/class_categories [get]
func GetClassCategories(c *fiber.Ctx) error {
	class_categories := []models.ClassCategory{}
	db, err := middlewares.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	if err := db.Find(&class_categories).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(200).JSON(class_categories)
}

func findClassCategory(db *gorm.DB, id int, class_category *models.ClassCategory) error {
	return db.First(class_category, "id = ?", id).Error
}

// GetClassCategory godoc
//
//		@Summary		Get class category by ID
//		@Description	GetClassCategory retrieves a single ClassCategory by its ID, including Classes
//		@Tags			ClassCategories
//	 @Security 		BearerAuth
//		@Produce		json
//		@Param			id	path		int	true	"ClassCategory ID"
//		@Success		200	{object}	models.ClassCategoryDoc
//		@Failure		400	{object}	map[string]string	"Invalid ID"
//		@Failure		404	{object}	map[string]string	"ClassCategory not found"
//		@Failure		500	{object}	map[string]string	"Server error"
//		@Router			/class_categories/{id} [get]
func GetClassCategory(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var class_category models.ClassCategory

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}
	db, err := middlewares.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	err = db.Preload("Classes").First(class_category, "id = ?", id).Error
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return c.Status(404).JSON("class_category not found")
	case err != nil:
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(200).JSON(class_category)
}

// UpdateClassCategory godoc
//
//		@Summary		Update an existing class category
//		@Description	UpdateClassCategory updates a ClassCategory record by its ID
//		@Tags			ClassCategories
//	 @Security 		BearerAuth
//		@Accept			json
//		@Produce		json
//		@Param			id				path		int						true	"ClassCategory ID"
//		@Param			class_category	body		models.ClassCategoryDoc	true	"Updated payload"
//		@Success		200				{object}	models.ClassCategoryDoc
//		@Failure		400				{object}	map[string]string	"Invalid input"
//		@Failure		404				{object}	map[string]string	"ClassCategory not found"
//		@Failure		500				{object}	map[string]string	"Server error"
//		@Router			/class_categories/{id} [put]
func UpdateClassCategory(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var class_category models.ClassCategory

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}
	db, err := middlewares.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	err = findClassCategory(db, id, &class_category)
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return c.Status(404).JSON("class_category not found")
	case err != nil:
		return c.Status(500).JSON(err.Error())
	}

	var class_category_update models.ClassCategory
	if err := c.BodyParser(&class_category_update); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if err := db.Model(&class_category).Omit("Classes").Updates(class_category_update).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(200).JSON(class_category)

}

// DeleteClassCategory godoc
//
//		@Summary		Delete a class category by ID
//		@Description	DeleteClassCategory removes a ClassCategory record by its ID
//		@Tags			ClassCategories
//	 @Security 		BearerAuth
//		@Produce		json
//		@Param			id	path		int					true	"ClassCategory ID"
//		@Success		200	{string}	string				"Successfully deleted class category"
//		@Failure		400	{object}	map[string]string	"Invalid ID"
//		@Failure		404	{object}	map[string]string	"ClassCategory not found"
//		@Failure		500	{object}	map[string]string	"Server error"
//		@Router			/class_categories/{id} [delete]
func DeleteClassCategory(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var class_category models.ClassCategory

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}
	db, err := middlewares.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	err = findClassCategory(db, id, &class_category)
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return c.Status(404).JSON("class_category not found")
	case err != nil:
		return c.Status(500).JSON(err.Error())
	}

	if err = db.Delete(&class_category).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}
	return c.Status(200).JSON("Successfully deleted class category")
}
