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

func setupEnrollmentIntegrationApp(t *testing.T) *fiber.App {
	app := fiber.New()

	// connect PostgreSQL
	cfg := config.NewConfig()
	db, err := config.ConnectDB(cfg)
	if err != nil {
		t.Fatalf("failed to connect db: %v", err)
	}

	// migrate dependencies
	if err := db.AutoMigrate(
		&models.Enrollment{},
		&models.Learner{},
		&models.ClassSession{},
		&models.Class{},
	); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}

	// inject db into context
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("db", db)
		return c.Next()
	})

	EnrollmentRoutes(app)
	return app
}

func TestIntegration_Enrollment_CRUD(t *testing.T) {
	app := setupEnrollmentIntegrationApp(t)

	// seed dependencies
	db := config.MustConnectTestDB(t)
	db.Create(&models.Learner{Model: models.BaseModel{ID: 1}})
	db.Create(&models.ClassSession{Model: models.BaseModel{ID: 1}})

	// Create
	reqBody := models.Enrollment{
		LearnerID:      1,
		ClassSessionID: 1,
		PaymentStatus:  "paid",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/enrollments/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var created models.Enrollment
	err := json.NewDecoder(resp.Body).Decode(&created)
	assert.NoError(t, err)
	assert.NotZero(t, created.ID)

	// Get All
	req = httptest.NewRequest(http.MethodGet, "/enrollments/", nil)
	resp, _ = app.Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Get By ID
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/enrollments/%d", created.ID), nil)
	resp, _ = app.Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Invalid ID
	req = httptest.NewRequest(http.MethodGet, "/enrollments/abc", nil)
	resp, _ = app.Test(req)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Update
	update := map[string]interface{}{
		"payment_status": "refunded",
	}
	upBody, _ := json.Marshal(update)
	req = httptest.NewRequest(http.MethodPut, fmt.Sprintf("/enrollments/%d", created.ID), bytes.NewBuffer(upBody))
	req.Header.Set("Content-Type", "application/json")
	resp, _ = app.Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Delete
	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/enrollments/%d", created.ID), nil)
	resp, _ = app.Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Not Found
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/enrollments/%d", created.ID), nil)
	resp, _ = app.Test(req)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}
