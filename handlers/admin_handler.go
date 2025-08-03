package handlers

import (
	"errors"

	"github.com/a2n2k3p4/tutorium-backend/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func AdminRoutes(app *fiber.App) {
	app.Post("/admin", CreateAdmin)
	app.Get("/admins", GetAdmins)
	app.Get("/admin/:id", GetAdmin)
	// app.Put("/admin/:id", UpdateAdmin) No application logic for updating admin
	app.Delete("/admin/:id", DeleteAdmin)
}

func CreateAdmin(c *fiber.Ctx) error {
	var admin models.Admin

	if err := c.BodyParser(&admin); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if err := db.Create(&admin).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}
	return c.Status(200).JSON(admin)
}

func GetAdmins(c *fiber.Ctx) error {
	admins := []models.Admin{}
	if err := db.Find(&admins).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(200).JSON(admins)
}

func findAdmin(id int, admin *models.Admin) error {
	return db.First(admin, "id = ?", id).Error
}

func GetAdmin(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var admin models.Admin

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	err = findAdmin(id, &admin)
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return c.Status(404).JSON("admin not found")
	case err != nil:
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(200).JSON(admin)
}

func DeleteAdmin(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var admin models.Admin

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	err = findAdmin(id, &admin)
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return c.Status(404).JSON("admin not found")
	case err != nil:
		return c.Status(500).JSON(err.Error())
	}

	if err = db.Delete(&admin).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}
	return c.Status(200).JSON("Successfully deleted admin")
}
