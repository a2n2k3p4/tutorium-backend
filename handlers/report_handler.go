package handlers

import (
	"errors"

	"github.com/a2n2k3p4/tutorium-backend/middlewares"
	"github.com/a2n2k3p4/tutorium-backend/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ReportRoutes(app *fiber.App) {
	report := app.Group("/reports", middlewares.ProtectedMiddleware())
	report.Post("/", CreateReport)

	reportAdmin := report.Group("/", middlewares.AdminRequired())
	reportAdmin.Get("/", GetReports)
	reportAdmin.Get("/:id", GetReport)
	reportAdmin.Put("/:id", UpdateReport)
	reportAdmin.Delete("/:id", DeleteReport)
}

// CreateReport godoc
//
//	@Summary		Create a new report
//	@Description	CreateReport creates a new Report record
//	@Tags			Reports
//	@Accept			json
//	@Produce		json
//	@Param			report	body		models.ReportDoc	true	"Report payload"
//	@Success		201		{object}	models.ReportDoc
//	@Failure		400		{object}	map[string]string	"Invalid input"
//	@Failure		500		{object}	map[string]string	"Server error"
//	@Router			/reports [post]
func CreateReport(c *fiber.Ctx) error {
	var report models.Report

	if err := c.BodyParser(&report); err != nil {
		return c.Status(400).JSON(err.Error())
	}
	db, err := middlewares.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	if err := db.Create(&report).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(201).JSON(report)
}

// GetReports godoc
//
//	@Summary		List all reports
//	@Description	GetReports retrieves all Report records with Reporter and Reported relations
//	@Tags			Reports
//	@Produce		json
//	@Success		200	{array}		models.ReportDoc
//	@Failure		500	{object}	map[string]string	"Server error"
//	@Router			/reports [get]
func GetReports(c *fiber.Ctx) error {
	reports := []models.Report{}
	db, err := middlewares.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	if err := db.Preload("Reporter").Preload("Reported").Find(&reports).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}
	return c.Status(200).JSON(reports)
}

func findReport(db *gorm.DB, id int, report *models.Report) error {
	return db.Preload("Reporter").Preload("Reported").First(report, "id = ?", id).Error
}

// GetReport godoc
//
//	@Summary		Get report by ID
//	@Description	GetReport retrieves a single Report by its ID, including Reporter and Reported
//	@Tags			Reports
//	@Produce		json
//	@Param			id	path		int	true	"Report ID"
//	@Success		200	{object}	models.ReportDoc
//	@Failure		400	{object}	map[string]string	"Invalid ID"
//	@Failure		404	{object}	map[string]string	"Report not found"
//	@Failure		500	{object}	map[string]string	"Server error"
//	@Router			/reports/{id} [get]
func GetReport(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var report models.Report

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}
	db, err := middlewares.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	err = findReport(db, id, &report)
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return c.Status(404).JSON("report not found")
	case err != nil:
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(200).JSON(report)
}

// UpdateReport godoc
//
//	@Summary		Update an existing report
//	@Description	UpdateReport updates a Report record by its ID
//	@Tags			Reports
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int					true	"Report ID"
//	@Param			report	body		models.ReportDoc	true	"Updated report payload"
//	@Success		200		{object}	models.ReportDoc
//	@Failure		400		{object}	map[string]string	"Invalid input"
//	@Failure		404		{object}	map[string]string	"Report not found"
//	@Failure		500		{object}	map[string]string	"Server error"
//	@Router			/reports/{id} [put]
func UpdateReport(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var report models.Report

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}
	db, err := middlewares.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	err = findReport(db, id, &report)
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return c.Status(404).JSON("report not found")
	case err != nil:
		return c.Status(500).JSON(err.Error())
	}

	var report_updated models.Report
	if err := c.BodyParser(&report_updated); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if err := db.Model(&report).Updates(report_updated).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(200).JSON(report)
}

// DeleteReport godoc
//
//	@Summary		Delete a report by ID
//	@Description	DeleteReport removes a Report record by its ID
//	@Tags			Reports
//	@Produce		json
//	@Param			id	path		int					true	"Report ID"
//	@Success		200	{string}	string				"Successfully deleted Report"
//	@Failure		400	{object}	map[string]string	"Invalid ID"
//	@Failure		404	{object}	map[string]string	"Report not found"
//	@Failure		500	{object}	map[string]string	"Server error"
//	@Router			/reports/{id} [delete]
func DeleteReport(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var report models.Report

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}
	db, err := middlewares.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	err = findReport(db, id, &report)
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return c.Status(404).JSON("report not found")
	case err != nil:
		return c.Status(500).JSON(err.Error())
	}

	if err = db.Delete(&report).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}
	return c.Status(200).JSON("Successfully deleted Report")
}
