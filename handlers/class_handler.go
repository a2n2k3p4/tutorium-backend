package handlers

import (
	"errors"

	"github.com/a2n2k3p4/tutorium-backend/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ClassRoutes(app *fiber.App) {
	app.Post("/class", CreateClass)
	app.Get("/classes", GetClasses)
	app.Get("/class/:id", GetClass)
	app.Put("/class/:id", UpdateClass)
	app.Delete("/class/:id", DeleteClass)
}

func CreateClass(c *fiber.Ctx) error {
	var class models.Class

	if err := c.BodyParser(&class); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if err := db.Create(&class).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}
	return c.Status(200).JSON(class)
}

func GetClasses(c *fiber.Ctx) error {
	classes := []models.Class{}
	if err := db.Preload("Teacher").Preload("Categories").Find(&classes).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(200).JSON(classes)
}

func findClass(id int, class *models.Class) error {
	return db.Preload("Teacher").Preload("Categories").First(class, "id = ?", id).Error
}

func GetClass(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var class models.Class

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	err = findClass(id, &class)
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return c.Status(404).JSON("class not found")
	case err != nil:
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(200).JSON(class)
}

func UpdateClass(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var class models.Class

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	err = findClass(id, &class)
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

func DeleteClass(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var class models.Class

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	err = findClass(id, &class)
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
