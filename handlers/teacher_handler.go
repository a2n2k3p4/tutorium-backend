package handlers

import (
	"errors"

	"github.com/a2n2k3p4/tutorium-backend/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func TeacherRoutes(app *fiber.App) {
	app.Post("/teacher", CreateTeacher)
	app.Get("/teachers", GetTeachers)
	app.Get("/teacher/:id", GetTeacher)
	app.Put("/teacher/:id", UpdateTeacher)
	app.Delete("/teacher/:id", DeleteTeacher)
}

func CreateTeacher(c *fiber.Ctx) error {
	var teacher models.Teacher

	if err := c.BodyParser(&teacher); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if err := db.Create(&teacher).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(200).JSON(teacher)
}

func GetTeachers(c *fiber.Ctx) error {
	teachers := []models.Teacher{}
	if err := db.Find(&teachers).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(200).JSON(teachers)
}

func findTeacher(id int, teacher *models.Teacher) error {
	if err := db.First(teacher, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("teacher does not exist")
		}
		return err
	}
	return nil
}

func GetTeacher(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var teacher models.Teacher

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	if err := findTeacher(id, &teacher); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	return c.Status(200).JSON(teacher)
}

func UpdateTeacher(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var teacher models.Teacher

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	if err := findTeacher(id, &teacher); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	var teacher_update models.Teacher
	if err := c.BodyParser(&teacher_update); err != nil {
		return c.Status(500).JSON(err.Error())
	}

	if err := db.Model(&teacher).Updates(teacher_update).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(200).JSON(teacher)
}

func DeleteTeacher(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var teacher models.Teacher

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	if err := findTeacher(id, &teacher); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if err = db.Delete(&teacher).Error; err != nil {
		return c.Status(404).JSON(err.Error())
	}
	return c.Status(200).JSON("Successfully deleted Teacher")
}
