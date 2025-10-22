package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/a2n2k3p4/tutorium-backend/middlewares"
	"github.com/a2n2k3p4/tutorium-backend/models"
	"github.com/a2n2k3p4/tutorium-backend/services"
	"github.com/a2n2k3p4/tutorium-backend/storage"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func ReportRoutes(app *fiber.App) {
	report := app.Group("/reports", middlewares.ProtectedMiddleware(), middlewares.BanMiddleware())
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
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			report	body		models.ReportDoc	true	"Report payload"
//	@Success		201		{object}	models.ReportDoc
//	@Failure		400		{string}	string	"Invalid input"
//	@Failure		500		{string}	string	"Server error"
//	@Router			/reports [post]
func CreateReport(c *fiber.Ctx) error {
	var report models.Report

	if err := c.BodyParser(&report); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if err := processReportPicture(c, &report); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	db, err := middlewares.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	if err := db.Create(&report).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	// Notify the reported user
	services.CreateNotification(
		db,
		report.ReportedUserID,
		"system",
		"You have been reported. An admin will review the case shortly.",
	)

	// Notify all admins
	var admins []models.Admin
	db.Find(&admins)
	for _, admin := range admins {
		services.CreateNotification(
			db,
			admin.UserID,
			"system",
			fmt.Sprintf("A new report (ID: %d) has been filed and requires your review.", report.ID),
		)
	}

	return c.Status(201).JSON(report)
}

// GetReports godoc
//
//	@Summary		List all reports
//	@Description	GetReports retrieves all Report records with Reporter and Reported relations
//	@Tags			Reports
//	@Security		BearerAuth
//	@Produce		json
//	@Success		200	{array}		models.ReportDoc
//	@Failure		500	{string}	string	"Server error"
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

	// Generate presigned URL if ReportPictureURL exists
	mc, ok := c.Locals("minio").(*storage.Client)
	if ok {
		for i := range reports {
			if reports[i].ReportPictureURL != "" {
				presignedURL, err := mc.PresignedGetObject(c.Context(), reports[i].ReportPictureURL, 15*time.Minute)
				if err == nil {
					reports[i].ReportPictureURL = presignedURL
				} else {
					reports[i].ReportPictureURL = ""
				}
			}
		}
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
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id	path		int	true	"Report ID"
//	@Success		200	{object}	models.ReportDoc
//	@Failure		400	{string}	string	"Invalid ID"
//	@Failure		404	{string}	string	"Report not found"
//	@Failure		500	{string}	string	"Server error"
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

	// Generate presigned URL if ReportPictureURL exists
	mc, ok := c.Locals("minio").(*storage.Client)
	if ok && report.ReportPictureURL != "" {
		presignedURL, err := mc.PresignedGetObject(c.Context(), report.ReportPictureURL, 15*time.Minute)
		if err == nil {
			report.ReportPictureURL = presignedURL
		} else {
			report.ReportPictureURL = ""
		}
	}

	return c.Status(200).JSON(report)
}

// UpdateReport godoc
//
//	@Summary		Update an existing report
//	@Description	UpdateReport updates a Report record by its ID
//	@Tags			Reports
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int					true	"Report ID"
//	@Param			report	body		models.ReportDoc	true	"Updated report payload"
//	@Success		200		{object}	models.ReportDoc
//	@Failure		400		{string}	string	"Invalid input"
//	@Failure		404		{string}	string	"Report not found"
//	@Failure		500		{string}	string	"Server error"
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

	if report.ReportStatus == "pending" && report_updated.ReportStatus == "resolve" {
		flags := 0
		switch report.ReportReason {
		case "teacher_absent", "poor_teaching", "not_teaching":
			flags = 1
		case "fake_review", "disruption", "disrespected", "harassment", "bullying":
			flags = 2
		}

		if flags > 0 {
			reason := fmt.Sprintf("Flagged from resolved report ID %d: %s", report.ID, report.ReportDescription)
			switch report.ReportType {
			case "learner": // Learner reported a Teacher
				var teacher models.Teacher
				if err := db.Where("user_id = ?", report.ReportedUserID).First(&teacher).Error; err == nil {
					desc := fmt.Sprintf("An admin has reviewed report #%d and issued a warning with %d flag(s).", report.ID, flags)
					services.CreateNotification(db, teacher.UserID, "system", desc)
					services.ApplyTeacherFlags(db, teacher.ID, flags, reason)
				}
			case "teacher": // Teacher reported a Learner
				var learner models.Learner
				if err := db.Where("user_id = ?", report.ReportedUserID).First(&learner).Error; err == nil {
					desc := fmt.Sprintf("An admin has reviewed report #%d and issued a warning with %d flag(s).", report.ID, flags)
					services.CreateNotification(db, learner.UserID, "system", desc)
					services.ApplyLearnerFlags(db, learner.ID, flags, reason)
				}
			}
		} // False Report
	} else if report.ReportStatus == "pending" && report_updated.ReportStatus == "reject" {
		reason := fmt.Sprintf("Flagged for submitting a false report (ID: %d)", report.ID)
		switch report.ReportType {
		case "learner": // The reporter was a Learner
			var learner models.Learner
			if err := db.Where("user_id = ?", report.ReportUserID).First(&learner).Error; err == nil {
				desc := fmt.Sprintf("Report #%d was found to be false. You have received 1 flag as a warning.", report.ID)
				services.CreateNotification(db, learner.UserID, "system", desc)
				services.ApplyLearnerFlags(db, learner.ID, 1, reason)
			}
		case "teacher": // The reporter was a Teacher
			var teacher models.Teacher
			if err := db.Where("user_id = ?", report.ReportUserID).First(&teacher).Error; err == nil {
				desc := fmt.Sprintf("Report #%d was found to be false. You have received 1 flag as a warning.", report.ID)
				services.CreateNotification(db, teacher.UserID, "system", desc)
				services.ApplyTeacherFlags(db, teacher.ID, 1, reason)
			}
		}
	}

	if err := processReportPicture(c, &report_updated); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if err := db.Model(&report).
		Omit(clause.Associations).
		Updates(report_updated).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(200).JSON(report)
}

// DeleteReport godoc
//
//	@Summary		Delete a report by ID
//	@Description	DeleteReport removes a Report record by its ID
//	@Tags			Reports
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id	path		int		true	"Report ID"
//	@Success		200	{string}	string	"Successfully deleted Report"
//	@Failure		400	{string}	string	"Invalid ID"
//	@Failure		404	{string}	string	"Report not found"
//	@Failure		500	{string}	string	"Server error"
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

func processReportPicture(c *fiber.Ctx, report *models.Report) error {
	if report.ReportPictureURL != "" && !strings.HasPrefix(report.ReportPictureURL, "http") {
		b, err := storage.DecodeBase64Image(report.ReportPictureURL)
		if err != nil {
			return fmt.Errorf("invalid base64 image: %w", err)
		}
		if err := validateImageBytes(b); err != nil {
			return fmt.Errorf("invalid image: %w", err)
		}
		mc := c.Locals("minio").(*storage.Client)
		filename := storage.GenerateFilename(http.DetectContentType(b[:min(512, len(b))]))
		uploaded, err := mc.UploadBytes(context.Background(), "reports", filename, b)
		if err != nil {
			return err
		}
		report.ReportPictureURL = uploaded
	}
	return nil
}
