package handlers

import (
	"errors"

	"github.com/a2n2k3p4/tutorium-backend/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ReportRoutes(app *fiber.App) {
	app.Post("/report", CreateReport)
	app.Get("/reports", GetReports)
	app.Get("/report/:id", GetReport)
	app.Put("/report/:id", UpdateReport)
	app.Delete("/report/:id", DeleteReport)
}

func CreateReport(c *fiber.Ctx) error {
	var report models.Report

	if err := c.BodyParser(&report); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if err := db.Create(&report).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(201).JSON(report)
}

func GetReports(c *fiber.Ctx) error {
	reports := []models.Report{}
	if err := db.Preload("Reporter").Preload("Reported").Find(&reports).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}
	return c.Status(200).JSON(reports)
}

func findReport(id int, report *models.Report) error {
	return db.Preload("Reporter").Preload("Reported").First(report, "id = ?", id).Error
}

func GetReport(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var report models.Report

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	err = findReport(id, &report)
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return c.Status(404).JSON("report not found")
	case err != nil:
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(200).JSON(report)
}

func UpdateReport(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var report models.Report

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	err = findReport(id, &report)
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

func DeleteReport(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var report models.Report

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	err = findReport(id, &report)
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
