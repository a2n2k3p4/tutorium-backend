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

// CreateBanTeacher godoc
// @Summary      Create a new ban record for a teacher
// @Description  CreateBanTeacher creates a new BanDetailsTeacher entry
// @Tags         BanTeachers
// @Accept       json
// @Produce      json
// @Param        banteacher  body      models.BanDetailsTeacherDoc  true  "BanTeacher payload"
// @Success      201         {object}  models.BanDetailsTeacherDoc
// @Failure      400         {object}  map[string]string         "Invalid input"
// @Failure      500         {object}  map[string]string         "Server error"
// @Router       /banteacher [post]
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

// GetBanTeachers godoc
// @Summary      List all ban records for teachers
// @Description  GetBanTeachers retrieves all BanDetailsTeacher entries with associated Teacher
// @Tags         BanTeachers
// @Produce      json
// @Success      200  {array}   models.BanDetailsTeacherDoc
// @Failure      500  {object}  map[string]string  "Server error"
// @Router       /banteachers [get]
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

// GetBanTeacher godoc
// @Summary      Get ban record by ID
// @Description  GetBanTeacher retrieves a single BanDetailsTeacher by its ID
// @Tags         BanTeachers
// @Produce      json
// @Param        id   path      int  true  "BanTeacher ID"
// @Success      200  {object}  models.BanDetailsTeacherDoc
// @Failure      400  {object}  map[string]string  "Invalid ID"
// @Failure      404  {object}  map[string]string  "BanTeacher not found"
// @Failure      500  {object}  map[string]string  "Server error"
// @Router       /banteacher/{id} [get]
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

// UpdateBanTeacher godoc
// @Summary      Update a ban record by ID
// @Description  UpdateBanTeacher updates an existing BanDetailsTeacher
// @Tags         BanTeachers
// @Accept       json
// @Produce      json
// @Param        id           path      int                     true  "BanTeacher ID"
// @Param        banteacher   body      models.BanDetailsTeacherDoc  true  "Updated payload"
// @Success      200          {object}  models.BanDetailsTeacherDoc
// @Failure      400          {object}  map[string]string         "Invalid input or not found"
// @Failure      404          {object}  map[string]string         "BanTeacher not found"
// @Failure      500          {object}  map[string]string         "Server error"
// @Router       /banteacher/{id} [put]
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

// DeleteBanTeacher godoc
// @Summary      Delete a ban record by ID
// @Description  DeleteBanTeacher removes a BanDetailsTeacher record by its ID
// @Tags         BanTeachers
// @Produce      json
// @Param        id   path      int  true  "BanTeacher ID"
// @Success      200  {string}  string  "Successfully deleted ban teacher"
// @Failure      400  {object}  map[string]string  "Invalid ID"
// @Failure      404  {object}  map[string]string  "BanTeacher not found"
// @Failure      500  {object}  map[string]string  "Server error"
// @Router       /banteacher/{id} [delete]
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
