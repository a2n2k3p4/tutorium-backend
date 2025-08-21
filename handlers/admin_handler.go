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

// CreateAdmin godoc
// @Summary Create a new admin
// @Description Create a new admin with the provided data
// @Tags admins
// @Accept json
// @Produce json
// @Param admin body models.AdminDoc true "Admin data"
// @Success 200 {object} models.AdminDoc
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Router /admin [post]
func CreateAdmin(c *fiber.Ctx) error {
	var admin models.Admin

	if err := c.BodyParser(&admin); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if err := db.Create(&admin).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}
	return c.Status(201).JSON(admin)
}

// GetAdmins godoc
// @Summary Get all admins
// @Description Retrieve a list of all admins
// @Tags admins
// @Accept json
// @Produce json
// @Success 200 {array} models.AdminDoc
// @Failure 404 {object} map[string]interface{} "Admins not found"
// @Router /admins [get]
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

// GetAdmin godoc
// @Summary Get admin by ID
// @Description Retrieve a specific admin by their ID
// @Tags admins
// @Accept json
// @Produce json
// @Param id path int true "Admin ID"
// @Success 200 {object} models.AdminDoc
// @Failure 400 {object} map[string]interface{} "Bad request - Invalid ID or admin not found"
// @Router /admin/{id} [get]
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

// DeleteAdmin godoc
// @Summary Delete admin by ID
// @Description Delete a specific admin by their ID
// @Tags admins
// @Accept json
// @Produce json
// @Param id path int true "Admin ID"
// @Success 200 {object} map[string]interface{} "Successfully deleted admin"
// @Failure 400 {object} map[string]interface{} "Bad request - Invalid ID or admin not found"
// @Failure 500 {object} map[string]interface{} "Internal server error during deletion"
// @Router /admin/{id} [delete]
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
