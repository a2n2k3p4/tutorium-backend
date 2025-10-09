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

func setupBanLearnerIntegrationApp(t *testing.T) *fiber.App {
	app := fiber.New()

	// connect PostgreSQL DB
	cfg := config.NewConfig()
	db, err := config.ConnectDB(cfg)
	if err != nil {
		t.Fatalf("failed to connect db: %v", err)
	}

	// migrate tables
	if err := db.AutoMigrate(&models.BanDetailsLearner{}, &models.Learner{}); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}

	// inject db
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("db", db)
		return c.Next()
	})

	// register ban learner routes
	BanLearnerRoutes(app)

	return app
}

func TestIntegration_BanLearner_CRUD(t *testing.T) {
	app := setupBanLearnerIntegrationApp(t)

	// Create
	now := time.Now()
	bl := models.BanDetailsLearner{
		LearnerID:      1,
		BanStart:       now,
		BanEnd:         now.Add(24 * time.Hour),
		BanDescription: "spamming",
	}
	body, _ := json.Marshal(bl)
	req := httptest.NewRequest(http.MethodPost, "/banlearners/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var created models.BanDetailsLearner
	err := json.NewDecoder(resp.Body).Decode(&created)
	assert.NoError(t, err)
	assert.NotZero(t, created.ID)

	// Get All
	req = httptest.NewRequest(http.MethodGet, "/banlearners/", nil)
	resp, _ = app.Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Get by ID
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/banlearners/%d", created.ID), nil)
	resp, _ = app.Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Bad Request
	req = httptest.NewRequest(http.MethodGet, "/banlearners/abc", nil)
	resp, _ = app.Test(req)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Update
	update := map[string]string{"ban_description": "updated reason"}
	upBody, _ := json.Marshal(update)
	req = httptest.NewRequest(http.MethodPut, fmt.Sprintf("/banlearners/%d", created.ID), bytes.NewBuffer(upBody))
	req.Header.Set("Content-Type", "application/json")
	resp, _ = app.Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Delete
	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/banlearners/%d", created.ID), nil)
	resp, _ = app.Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Not Found
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/banlearners/%d", created.ID), nil)
	resp, _ = app.Test(req)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}
