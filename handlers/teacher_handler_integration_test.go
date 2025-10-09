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

func setupTeacherIntegrationApp(t *testing.T) *fiber.App {
	app := fiber.New()

	cfg := config.NewConfig()
	db, err := config.ConnectDB(cfg)
	if err != nil {
		t.Fatalf("failed to connect db: %v", err)
	}

	// migrate teacher table
	if err := db.AutoMigrate(&models.Teacher{}); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}

	// inject db into fiber context
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("db", db)
		return c.Next()
	})

	TeacherRoutes(app)
	return app
}

func TestIntegration_Teacher_CRUD(t *testing.T) {
	app := setupTeacherIntegrationApp(t)

	// Create
	teacher := models.Teacher{
		Name:     "John Doe",
		Email:    "john.doe@example.com",
		Subject:  "Mathematics",
		Biography: "Experienced teacher with 10 years of teaching calculus.",
	}
	body, _ := json.Marshal(teacher)
	req := httptest.NewRequest(http.MethodPost, "/teachers/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var created models.Teacher
	err := json.NewDecoder(resp.Body).Decode(&created)
	assert.NoError(t, err)
	assert.NotZero(t, created.ID)
	assert.Equal(t, teacher.Email, created.Email)

	// Get All
	req = httptest.NewRequest(http.MethodGet, "/teachers/", nil)
	resp, _ = app.Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Get by ID
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/teachers/%d", created.ID), nil)
	resp, _ = app.Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Invalid ID
	req = httptest.NewRequest(http.MethodGet, "/teachers/abc", nil)
	resp, _ = app.Test(req)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Update
	updatePayload := map[string]interface{}{
		"Biography": "Updated bio: now teaching advanced algebra.",
	}
	body, _ = json.Marshal(updatePayload)
	req = httptest.NewRequest(http.MethodPut, fmt.Sprintf("/teachers/%d", created.ID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ = app.Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Delete
	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/teachers/%d", created.ID), nil)
	resp, _ = app.Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Not Found
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/teachers/%d", created.ID), nil)
	resp, _ = app.Test(req)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}