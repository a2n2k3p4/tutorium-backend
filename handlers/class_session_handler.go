package handlers

import (
	"errors"

	"github.com/a2n2k3p4/tutorium-backend/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ClassSessionRoutes(app *fiber.App) {
	app.Post("/class_session", CreateClassSession)
	app.Get("/class_sessions", GetClassSessions)
	app.Get("/class_session/:id", GetClassSession)
	app.Put("/class_session/:id", UpdateClassSession)
	app.Delete("/class_session/:id", DeleteClassSession)
}

func CreateClassSession(c *fiber.Ctx) error {
	var class_session models.ClassSession

	if err := c.BodyParser(&class_session); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if err := db.Create(&class_session).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(200).JSON(class_session)
}

func GetClassSessions(c *fiber.Ctx) error {
	class_sessions := []models.ClassSession{}
	if err := db.Preload("Class").Find(&class_sessions).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(200).JSON(class_sessions)
}

func findClassSession(id int, class_session *models.ClassSession) error {
	return db.Preload("Class").First(class_session, "id = ?", id).Error
}

func GetClassSession(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var class_session models.ClassSession

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	err = findClassSession(id, &class_session)
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return c.Status(404).JSON("class_session not found")
	case err != nil:
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(200).JSON(class_session)
}

func UpdateClassSession(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var class_session models.ClassSession

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	err = findClassSession(id, &class_session)
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

func DeleteClassSession(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var class_session models.ClassSession

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	err = findClassSession(id, &class_session)
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
