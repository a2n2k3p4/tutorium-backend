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

func setupClassSessionIntegrationApp(t *testing.T) *fiber.App {
	app := fiber.New()

	// connect PostgreSQL
	cfg := config.NewConfig()
	db, err := config.ConnectDB(cfg)
	if err != nil {
		t.Fatalf("failed to connect db: %v", err)
	}

	// migrate dependencies
	if err := db.AutoMigrate(&models.ClassSession{}, &models.Class{}); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}

	// inject db
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("db", db)
		return c.Next()
	})

	ClassSessionRoutes(app)
	return app
}

func TestIntegration_ClassSession_CRUD(t *testing.T) {
	app := setupClassSessionIntegrationApp(t)

	// Create
	reqBody := models.CreateClassSessionRequest{
		ClassID:            1,
		Description:        "Go Programming",
		Price:              100,
		LearnerLimit:       60,
		EnrollmentDeadline: time.Now().Add(24 * time.Hour),
		ClassStart:         time.Now().Add(48 * time.Hour),
		ClassFinish:        time.Now().Add(72 * time.Hour),
		ClassStatus:        "upcoming",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/class_sessions/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var created models.ClassSession
	err := json.NewDecoder(resp.Body).Decode(&created)
	assert.NoError(t, err)
	assert.NotZero(t, created.ID)
	assert.Contains(t, created.MeetingUrl, "KUtutorium_Go_Programming")

	// Get All
	req = httptest.NewRequest(http.MethodGet, "/class_sessions/", nil)
	resp, _ = app.Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Get By ID
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/class_sessions/%d", created.ID), nil)
	resp, _ = app.Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Bad Request
	req = httptest.NewRequest(http.MethodGet, "/class_sessions/abc", nil)
	resp, _ = app.Test(req)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Update
	update := map[string]interface{}{
		"description": "Go Programming Advanced",
		"price":       150,
	}
	upBody, _ := json.Marshal(update)
	req = httptest.NewRequest(http.MethodPut, fmt.Sprintf("/class_sessions/%d", created.ID), bytes.NewBuffer(upBody))
	req.Header.Set("Content-Type", "application/json")
	resp, _ = app.Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Delete
	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/class_sessions/%d", created.ID), nil)
	resp, _ = app.Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Not Found
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/class_sessions/%d", created.ID), nil)
	resp, _ = app.Test(req)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}