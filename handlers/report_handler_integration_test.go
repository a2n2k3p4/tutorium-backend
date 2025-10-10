package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/a2n2k3p4/tutorium-backend/config"
	"github.com/a2n2k3p4/tutorium-backend/models"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func setupReportIntegrationApp(t *testing.T) *fiber.App {
	app := fiber.New()

	cfg := config.NewConfig()
	db, err := config.ConnectDB(cfg)
	if err != nil {
		t.Fatalf("failed to connect db: %v", err)
	}

	if err := db.AutoMigrate(&models.Report{}); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("db", db)
		return c.Next()
	})

	ReportRoutes(app)
	return app
}

func TestIntegration_Report_CRUD(t *testing.T) {
	app := setupReportIntegrationApp(t)

	// Create
	report := models.Report{
		ReportUserID:      99999,
		ReportedUserID:    88888,
		ReportType:        "Inappropriate behavior",
		ReportDescription: "Reported user used the most offensive english word in class.",
		ReportPictureURL:  "",
		ReportDate:        time.Now(),
	}
	body, _ := json.Marshal(report)
	req := httptest.NewRequest(http.MethodPost, "/reports/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var created models.Report
	err := json.NewDecoder(resp.Body).Decode(&created)
	assert.NoError(t, err)
	assert.NotZero(t, created.ID)

	// Get All
	req = httptest.NewRequest(http.MethodGet, "/reports/", nil)
	resp, _ = app.Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Get by ID
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/reports/%d", created.ID), nil)
	resp, _ = app.Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Bad Request
	req = httptest.NewRequest(http.MethodGet, "/reports/abc", nil)
	resp, _ = app.Test(req)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Update
	updated := map[string]interface{}{
		"Status": "resolved",
		"Detail": "Case closed by admin.",
	}
	body, _ = json.Marshal(updated)
	req = httptest.NewRequest(http.MethodPut, fmt.Sprintf("/reports/%d", created.ID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ = app.Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Delete
	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/reports/%d", created.ID), nil)
	resp, _ = app.Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Not Found
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/reports/%d", created.ID), nil)
	resp, _ = app.Test(req)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}