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

func setupNotificationIntegrationApp(t *testing.T) *fiber.App {
	app := fiber.New()

	cfg := config.NewConfig()
	db, err := config.ConnectDB(cfg)
	if err != nil {
		t.Fatalf("failed to connect db: %v", err)
	}

	// migrate table
	if err := db.AutoMigrate(&models.Notification{}); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}

	// inject db
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("db", db)
		return c.Next()
	})

	NotificationRoutes(app)
	return app
}

func TestIntegration_Notification_CRUD(t *testing.T) {
	app := setupNotificationIntegrationApp(t)

	// Create
	noti := models.Notification{
		UserID:  99900,
		NotificationType:   		"New Feature Update",
		NotificationDescription: 	"Our UI are no longer overlaps on Android devices!!!",
		NotificationDate:  			time.Now(),
		ReadFlag:  					false,
	}
	
	body, _ := json.Marshal(noti)
	req := httptest.NewRequest(http.MethodPost, "/notifications/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var created models.Notification
	err := json.NewDecoder(resp.Body).Decode(&created)
	assert.NoError(t, err)
	assert.NotZero(t, created.ID)

	// Get All
	req = httptest.NewRequest(http.MethodGet, "/notifications/", nil)
	resp, _ = app.Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Get by ID
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/notifications/%d", created.ID), nil)
	resp, _ = app.Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Bad Request (non-integer id)
	req = httptest.NewRequest(http.MethodGet, "/notifications/abc", nil)
	resp, _ = app.Test(req)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Update
	updated := map[string]interface{}{
		"Title":  "Updated Notification",
		"IsRead": true,
	}
	body, _ = json.Marshal(updated)
	req = httptest.NewRequest(http.MethodPut, fmt.Sprintf("/notifications/%d", created.ID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ = app.Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Delete
	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/notifications/%d", created.ID), nil)
	resp, _ = app.Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Not Found
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/notifications/%d", created.ID), nil)
	resp, _ = app.Test(req)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}