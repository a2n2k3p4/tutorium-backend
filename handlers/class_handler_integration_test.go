package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/a2n2k3p4/tutorium-backend/config"
	"github.com/a2n2k3p4/tutorium-backend/models"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func setupClassIntegrationApp(t *testing.T) *fiber.App {
	app := fiber.New()

	// connect PostgreSQL
	cfg := config.NewConfig()
	db, err := config.ConnectDB(cfg)
	if err != nil {
		t.Fatalf("failed to connect db: %v", err)
	}

	// migrate tables (Class + Category + Teacher + User)
	if err := db.AutoMigrate(&models.Class{}, &models.ClassCategory{}, &models.Teacher{}, &models.User{}); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}

	// inject db
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("db", db)
		return c.Next()
	})

	// mock minio client (disable upload for testing)
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("minio", nil)
		return c.Next()
	})

	ClassRoutes(app)

	return app
}

func TestIntegration_Class_CRUD(t *testing.T) {
	app := setupClassIntegrationApp(t)

	// Create class
	class := models.Class{
		ClassName: "Go Programming",
		Rating:    4,
		TeacherID: 1, // ensure you have a Teacher seeded with ID 1
	}
	body, _ := json.Marshal(class)
	req := httptest.NewRequest(http.MethodPost, "/classes/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var created models.Class
	err := json.NewDecoder(resp.Body).Decode(&created)
	assert.NoError(t, err)
	assert.NotZero(t, created.ID)

	// Get all classes
	req = httptest.NewRequest(http.MethodGet, "/classes/", nil)
	resp, _ = app.Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Get class by ID
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/classes/%d", created.ID), nil)
	resp, _ = app.Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Bad Request (non-int id)
	req = httptest.NewRequest(http.MethodGet, "/classes/abc", nil)
	resp, _ = app.Test(req)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Update class
	update := map[string]interface{}{
		"class_name": "Go Programming Advanced",
		"rating":     4.5,
	}
	upBody, _ := json.Marshal(update)
	req = httptest.NewRequest(http.MethodPut, fmt.Sprintf("/classes/%d", created.ID), bytes.NewBuffer(upBody))
	req.Header.Set("Content-Type", "application/json")
	resp, _ = app.Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Delete class
	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/classes/%d", created.ID), nil)
	resp, _ = app.Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Not Found
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/classes/%d", created.ID), nil)
	resp, _ = app.Test(req)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}