package handlers

import (
	"errors"

	"github.com/a2n2k3p4/tutorium-backend/models"
	"github.com/gofiber/fiber/v2"
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

	db.Create(&admin)
	return c.Status(200).JSON(admin)
}

func GetAdmins(c *fiber.Ctx) error {
	admins := []models.Admin{}
	dbErr := db.Find(&admins).Error
	if dbErr != nil {
		return c.Status(404).JSON(dbErr)
	}

	return c.Status(200).JSON(admins)
}

func findadmin(id int, admin *models.Admin) error {
	dbErr := db.Find(&admin, "id = ?", id).Error
	if dbErr != nil {
		return errors.New(dbErr.Error())
	}

	if admin.ID == 0 {
		return errors.New("admin does not exist")
	}
	return nil
}

func GetAdmin(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var admin models.Admin

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	if err := findadmin(id, &admin); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	return c.Status(200).JSON(admin)
}

func DeleteAdmin(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var admin models.Admin

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	err = findadmin(id, &admin)

	if err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if err = db.Delete(&admin).Error; err != nil {
		return c.Status(404).JSON(err.Error())
	}
	return c.Status(200).JSON("Successfully deleted admin")
}
