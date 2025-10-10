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

func setupBanTeacherIntegrationApp(t *testing.T) *fiber.App {
	app := fiber.New()

	// connect PostgreSQL
	cfg := config.NewConfig()
	db, err := config.ConnectDB(cfg)
	if err != nil {
		t.Fatalf("failed to connect db: %v", err)
	}

	// migrate tables
	if err := db.AutoMigrate(&models.BanDetailsTeacher{}, &models.Teacher{}); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}

	// inject db into ctx
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("db", db)
		return c.Next()
	})

	// register routes
	BanTeacherRoutes(app)

	return app
}

func TestIntegration_BanTeacher_CRUD(t *testing.T) {
	app := setupBanTeacherIntegrationApp(t)

	// Create
	now := time.Now()
	bt := models.BanDetailsTeacher{
		TeacherID:      1,
		BanStart:       now,
		BanEnd:         now.Add(24 * time.Hour),
		BanDescription: "policy violation",
	}
	body, _ := json.Marshal(bt)
	req := httptest.NewRequest(http.MethodPost, "/banteachers/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var created models.BanDetailsTeacher
	err := json.NewDecoder(resp.Body).Decode(&created)
	assert.NoError(t, err)
	assert.NotZero(t, created.ID)

	// Get All
	req = httptest.NewRequest(http.MethodGet, "/banteachers/", nil)
	resp, _ = app.Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Get by ID
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/banteachers/%d", created.ID), nil)
	resp, _ = app.Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Bad Request
	req = httptest.NewRequest(http.MethodGet, "/banteachers/abc", nil)
	resp, _ = app.Test(req)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Update
	update := map[string]string{"ban_description": "updated reason"}
	upBody, _ := json.Marshal(update)
	req = httptest.NewRequest(http.MethodPut, fmt.Sprintf("/banteachers/%d", created.ID), bytes.NewBuffer(upBody))
	req.Header.Set("Content-Type", "application/json")
	resp, _ = app.Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Delete
	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/banteachers/%d", created.ID), nil)
	resp, _ = app.Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Not Found
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/banteachers/%d", created.ID), nil)
	resp, _ = app.Test(req)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}