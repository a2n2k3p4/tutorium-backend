package handlers

import (
	"errors"
	"fmt"

	"github.com/a2n2k3p4/tutorium-backend/middlewares"
	"github.com/a2n2k3p4/tutorium-backend/models"
	"github.com/a2n2k3p4/tutorium-backend/services"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func AdminRoutes(app *fiber.App) {
	admin := app.Group("/admins", middlewares.ProtectedMiddleware())

	admin.Post("/", CreateAdmin)
	admin.Get("/", GetAdmins)
	admin.Get("/:id", GetAdmin)
	// admin.Put("/admin/:id", UpdateAdmin) No application logic for updating admin
	admin.Delete("/:id", DeleteAdmin)

	admin.Post("/flags/learners", middlewares.AdminRequired(), AddLearnerFlag)
	admin.Post("/flags/teachers", middlewares.AdminRequired(), AddTeacherFlag)
}

// CreateAdmin godoc
//
//	@Summary		Create a new admin
//	@Description	Create a new admin with the provided data
//	@Tags			Admins
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			admin	body		models.AdminDoc	true	"Admin data"
//	@Success		201		{object}	models.AdminDoc
//	@Failure		400		{string}	string	"Invalid request body"
//	@Failure		500		{string}	string	"Internal server error"
//	@Router			/admins [post]
func CreateAdmin(c *fiber.Ctx) error {
	var admin models.Admin

	if err := c.BodyParser(&admin); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	db, err := middlewares.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	if err := db.Create(&admin).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}
	return c.Status(201).JSON(admin)
}

// GetAdmins godoc
//
//	@Summary		Get all admins
//	@Description	Retrieve a list of all admins
//	@Tags			Admins
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		models.AdminDoc
//	@Failure		500	{string}	string	"Internal server error"
//	@Router			/admins [get]
func GetAdmins(c *fiber.Ctx) error {
	admins := []models.Admin{}
	db, err := middlewares.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	if err := db.Find(&admins).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(200).JSON(admins)
}

func findAdmin(db *gorm.DB, id int, admin *models.Admin) error {
	return db.First(admin, "id = ?", id).Error
}

// GetAdmin godoc
//
//	@Summary		Get admin by ID
//	@Description	Retrieve a specific admin by their ID
//	@Tags			Admins
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Admin ID"
//	@Success		200	{object}	models.AdminDoc
//	@Failure		400	{string}	string	"Invalid admin ID"
//	@Failure		404	{string}	string	"admin not found"
//	@Failure		500	{string}	string	"Internal server error"
//	@Router			/admins/{id} [get]
func GetAdmin(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var admin models.Admin

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	db, err := middlewares.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	err = findAdmin(db, id, &admin)
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return c.Status(404).JSON("admin not found")
	case err != nil:
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(200).JSON(admin)
}

// DeleteAdmin godoc
//
//	@Summary		Delete admin by ID
//	@Description	Delete a specific admin by their ID
//	@Tags			Admins
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int		true	"Admin ID"
//	@Success		200	{string}	string	"Successfully deleted admin"
//	@Failure		400	{string}	string	"Invalid admin ID"
//	@Failure		404	{string}	string	"admin not found"
//	@Failure		500	{string}	string	"Internal server error"
//	@Router			/admins/{id} [delete]
func DeleteAdmin(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var admin models.Admin

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	db, err := middlewares.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	err = findAdmin(db, id, &admin)
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

type FlagRequest struct {
	ID         uint   `json:"id"` // This would be LearnerID or TeacherID
	FlagsToAdd int    `json:"flags_to_add"`
	Reason     string `json:"reason"`
}

// AddLearnerFlag godoc
//
//	@Summary		Flag a learner
//	@Description	Allows an admin to apply flags directly to a learner
//	@Tags			Admins
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			flag	body		FlagRequest	true	"Flag details"
//	@Success		200		{object}	map[string]string	"Flags applied successfully"
//	@Failure		400		{object}	map[string]string	"Invalid request body"
//	@Failure		500		{object}	map[string]string	"Failed to apply flags"
//	@Router			/admins/flags/learner [post]
func AddLearnerFlag(c *fiber.Ctx) error {
	var req FlagRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	db, err := middlewares.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	var learner models.Learner
	if err := db.First(&learner, req.ID).Error; err == nil {
		desc := fmt.Sprintf("An administrator has issued a warning with %d flag(s). Reason: %s", req.FlagsToAdd, req.Reason)
		services.CreateNotification(db, learner.UserID, "system", desc)
	}

	if err := services.ApplyLearnerFlags(db, req.ID, req.FlagsToAdd, req.Reason); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to apply flags", "details": err.Error()})
	}
	return c.Status(200).JSON(fiber.Map{"message": "Flags applied successfully to learner"})
}

// AddTeacherFlag godoc
//
//	@Summary		Flag a teacher
//	@Description	Allows an admin to apply flags directly to a teacher
//	@Tags			Admins
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			flag	body		FlagRequest	true	"Flag details"
//	@Success		200		{object}	map[string]string	"Flags applied successfully"
//	@Failure		400		{object}	map[string]string	"Invalid request body"
//	@Failure		500		{object}	map[string]string	"Failed to apply flags"
//	@Router			/admins/flags/teacher [post]
func AddTeacherFlag(c *fiber.Ctx) error {
	var req FlagRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	db, err := middlewares.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	var teacher models.Teacher
	if err := db.First(&teacher, req.ID).Error; err == nil {
		desc := fmt.Sprintf("An administrator has issued a warning with %d flag(s). Reason: %s", req.FlagsToAdd, req.Reason)
		services.CreateNotification(db, teacher.UserID, "system", desc)
	}

	if err := services.ApplyTeacherFlags(db, req.ID, req.FlagsToAdd, req.Reason); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to apply flags", "details": err.Error()})
	}
	return c.Status(200).JSON(fiber.Map{"message": "Flags applied successfully to teacher"})
}
