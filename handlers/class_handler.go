package handlers

import (
	"errors"

	"github.com/a2n2k3p4/tutorium-backend/models"
	"github.com/gofiber/fiber/v2"
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

	db.Create(&class)
	return c.Status(200).JSON(class)
}

func GetClasses(c *fiber.Ctx) error {
	classes := []models.Class{}
	dbErr := db.Find(&classes).Error
	if dbErr != nil {
		return c.Status(404).JSON(dbErr)
	}

	return c.Status(200).JSON(classes)
}

func findclass(id int, class *models.Class) error {
	dbErr := db.Find(&class, "id = ?", id).Error
	if dbErr != nil {
		return errors.New(dbErr.Error())
	}

	if class.ID == 0 {
		return errors.New("class does not exist")
	}
	return nil
}

func GetClass(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var class models.Class

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	if err := findclass(id, &class); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	return c.Status(200).JSON(class)
}

func UpdateClass(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var class models.Class

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	err = findclass(id, &class)

	if err != nil {
		return c.Status(400).JSON(err.Error())
	}

	var class_update models.Class
	if err := c.BodyParser(&class_update); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	db.Model(&class).Updates(class_update)

	return c.Status(200).JSON(class)

}

func DeleteClass(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var class models.Class

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	err = findclass(id, &class)

	if err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if err = db.Delete(&class).Error; err != nil {
		return c.Status(404).JSON(err.Error())
	}
	return c.Status(200).JSON("Successfully deleted class")
}
