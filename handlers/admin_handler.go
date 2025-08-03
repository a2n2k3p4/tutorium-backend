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

// CreateAdmin godoc
// @Summary Create a new admin
// @Description Create a new admin with the provided data
// @Tags admins
// @Accept json
// @Produce json
// @Param admin body models.Admin true "Admin data"
// @Success 200 {object} models.Admin
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Router /admin [post]
func CreateAdmin(c *fiber.Ctx) error {
	var admin models.Admin

	if err := c.BodyParser(&admin); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	db.Create(&admin)
	return c.Status(200).JSON(admin)
}

// GetAdmins godoc
// @Summary Get all admins
// @Description Retrieve a list of all admins
// @Tags admins
// @Accept json
// @Produce json
// @Success 200 {array} models.Admin
// @Failure 404 {object} map[string]interface{} "Admins not found"
// @Router /admins [get]
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

// GetAdmin godoc
// @Summary Get admin by ID
// @Description Retrieve a specific admin by their ID
// @Tags admins
// @Accept json
// @Produce json
// @Param id path int true "Admin ID"
// @Success 200 {object} models.Admin
// @Failure 400 {object} map[string]interface{} "Bad request - Invalid ID or admin not found"
// @Router /admin/{id} [get]
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

// DeleteAdmin godoc
// @Summary Delete admin by ID
// @Description Delete a specific admin by their ID
// @Tags admins
// @Accept json
// @Produce json
// @Param id path int true "Admin ID"
// @Success 200 {object} map[string]interface{} "Successfully deleted admin"
// @Failure 400 {object} map[string]interface{} "Bad request - Invalid ID or admin not found"
// @Failure 404 {object} map[string]interface{} "Database error during deletion"
// @Router /admin/{id} [delete]
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
