package handlers

import (
	"errors"

	"github.com/a2n2k3p4/tutorium-backend/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ClassCategoryRoutes(database *gorm.DB, app *fiber.App) {
	db = database

	app.Post("/class_category", CreateClassCategory)
	app.Get("/class_categories", GetClassCategories)
	app.Get("/class_category/:id", GetClassCategory)
	app.Put("/class_category/:id", UpdateClassCategory)
	app.Delete("/class_category/:id", DeleteClassCategory)
}

func CreateClassCategory(c *fiber.Ctx) error {
	var class_category models.ClassCategory

	if err := c.BodyParser(&class_category); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	db.Create(&class_category)
	return c.Status(200).JSON(class_category)
}

func GetClassCategories(c *fiber.Ctx) error {
	class_categories := []models.ClassCategory{}
	db.Find(&class_categories)

	return c.Status(200).JSON(class_categories)
}

func findclasscategory(id int, class_category *models.ClassCategory) error {
	db.Find(&class_category, "id = ?", id)
	if class_category.ID == 0 {
		return errors.New("class category does not exist")
	}
	return nil
}

func GetClassCategory(c *fiber.Ctx) error {
	id, err := c.ParamsInt("ID")

	var class_category models.ClassCategory

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	if err := findclasscategory(id, &class_category); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	return c.Status(200).JSON(class_category)
}

func UpdateClassCategory(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var class_category models.ClassCategory

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	err = findclasscategory(id, &class_category)

	if err != nil {
		return c.Status(400).JSON(err.Error())
	}

	var class_category_update models.ClassCategory
	if err := c.BodyParser(&class_category_update); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	db.Model(&class_category).Updates(class_category_update)

	return c.Status(200).JSON(class_category)

}

func DeleteClassCategory(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var class_category models.ClassCategory

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	err = findclasscategory(id, &class_category)

	if err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if err = db.Delete(&class_category).Error; err != nil {
		return c.Status(404).JSON(err.Error())
	}
	return c.Status(200).JSON("Successfully deleted class category")
}
