package handlers

import (
	"errors"

	"github.com/a2n2k3p4/tutorium-backend/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func BanTeacherRoutes(app *fiber.App) {
	app.Post("/banteacher", CreateBanTeacher)
	app.Get("/banteachers", GetBanTeachers)
	app.Get("/banteacher/:id", GetBanTeacher)
	app.Put("/banteacher/:id", UpdateBanTeacher)
	app.Delete("/banteacher/:id", DeleteBanTeacher)
}

func CreateBanTeacher(c *fiber.Ctx) error {
	var banteacher models.BanDetailsTeacher

	if err := c.BodyParser(&banteacher); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if err := db.Create(&banteacher).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}
	return c.Status(201).JSON(banteacher)
}

func GetBanTeachers(c *fiber.Ctx) error {
	banteachers := []models.BanDetailsTeacher{}
	if err := db.Preload("Teacher").Find(&banteachers).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(200).JSON(banteachers)
}

func findBanTeacher(id int, banteacher *models.BanDetailsTeacher) error {
	return db.Preload("Teacher").First(banteacher, "id = ?", id).Error
}

func GetBanTeacher(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var banteacher models.BanDetailsTeacher

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	err = findBanTeacher(id, &banteacher)
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return c.Status(404).JSON("banteacher not found")
	case err != nil:
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(200).JSON(banteacher)
}

func UpdateBanTeacher(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var banteacher models.BanDetailsTeacher

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	err = findBanTeacher(id, &banteacher)
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return c.Status(404).JSON("banteacher not found")
	case err != nil:
		return c.Status(500).JSON(err.Error())
	}

	var banteacher_update models.BanDetailsTeacher
	if err := c.BodyParser(&banteacher_update); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if err := db.Model(&banteacher).Updates(banteacher_update).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(200).JSON(banteacher)
}

func DeleteBanTeacher(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var banteacher models.BanDetailsTeacher

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	err = findBanTeacher(id, &banteacher)
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return c.Status(404).JSON("banteacher not found")
	case err != nil:
		return c.Status(500).JSON(err.Error())
	}

	if err = db.Delete(&banteacher).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}
	return c.Status(200).JSON("Successfully deleted ban teacher")
}
