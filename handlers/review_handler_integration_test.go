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

func setupReviewIntegrationApp(t *testing.T) *fiber.App {
	app := fiber.New()

	cfg := config.NewConfig()
	db, err := config.ConnectDB(cfg)
	if err != nil {
		t.Fatalf("failed to connect db: %v", err)
	}

	// migrate review table
	if err := db.AutoMigrate(&models.Review{}); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}

	// inject db into fiber context
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("db", db)
		return c.Next()
	})

	ReviewRoutes(app)
	return app
}

func TestIntegration_Review_CRUD(t *testing.T) {
	app := setupReviewIntegrationApp(t)

	// Create
	review := models.Review{
		LearnerID: 1,
		ClassID:   101,
		Rating:    4,
		Comment:   "Thank you so much!!! Now I can do it by myself",
	}
	body, _ := json.Marshal(review)
	req := httptest.NewRequest(http.MethodPost, "/reviews/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var created models.Review
	err := json.NewDecoder(resp.Body).Decode(&created)
	assert.NoError(t, err)
	assert.NotZero(t, created.ID)

	// Get All
	req = httptest.NewRequest(http.MethodGet, "/reviews/", nil)
	resp, _ = app.Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Get by ID
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/reviews/%d", created.ID), nil)
	resp, _ = app.Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Bad Request
	req = httptest.NewRequest(http.MethodGet, "/reviews/abc", nil)
	resp, _ = app.Test(req)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Update
	updatePayload := map[string]interface{}{
		"Rating":  5,
		"Comment": "Updated comment: Excellent session!",
	}
	body, _ = json.Marshal(updatePayload)
	req = httptest.NewRequest(http.MethodPut, fmt.Sprintf("/reviews/%d", created.ID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ = app.Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Delete
	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/reviews/%d", created.ID), nil)
	resp, _ = app.Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Not Found
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/reviews/%d", created.ID), nil)
	resp, _ = app.Test(req)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}