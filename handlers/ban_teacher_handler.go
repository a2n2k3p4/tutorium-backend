package handlers

import (
	"errors"

	"github.com/a2n2k3p4/tutorium-backend/models"
	"github.com/gofiber/fiber/v2"
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

	db.Create(&banteacher)
	return c.Status(200).JSON(banteacher)
}

func GetBanTeachers(c *fiber.Ctx) error {
	banteachers := []models.BanDetailsTeacher{}
	dbErr := db.Find(&banteachers).Error
	if dbErr != nil {
		return c.Status(404).JSON(dbErr)
	}

	return c.Status(200).JSON(banteachers)
}

func findbanteacher(id int, banteacher *models.BanDetailsTeacher) error {
	dbErr := db.Find(&banteacher, "id = ?", id).Error
	if dbErr != nil {
		return errors.New(dbErr.Error())
	}

	if banteacher.ID == 0 {
		return errors.New("ban teacher does not exist")
	}
	return nil
}

func GetBanTeacher(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var banteacher models.BanDetailsTeacher

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	if err := findbanteacher(id, &banteacher); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	return c.Status(200).JSON(banteacher)
}

func UpdateBanTeacher(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var banteacher models.BanDetailsTeacher

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	err = findbanteacher(id, &banteacher)
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}

	var banteacher_update models.BanDetailsTeacher
	if err := c.BodyParser(&banteacher_update); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	db.Model(&banteacher).Updates(banteacher_update)

	return c.Status(200).JSON(banteacher)
}

func DeleteBanTeacher(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var banteacher models.BanDetailsTeacher

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	err = findbanteacher(id, &banteacher)
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if err = db.Delete(&banteacher).Error; err != nil {
		return c.Status(404).JSON(err.Error())
	}
	return c.Status(200).JSON("Successfully deleted ban teacher")
}
