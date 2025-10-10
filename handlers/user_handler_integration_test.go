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

type MockMinioClient struct{}

func (m *MockMinioClient) UploadBytes(ctx any, bucket, filename string, b []byte) (string, error) {
	return fmt.Sprintf("mock://%s/%s", bucket, filename), nil
}
func (m *MockMinioClient) PresignedGetObject(ctx any, objectKey string, expireMinutes int64) (string, error) {
	return fmt.Sprintf("mock://signed/%s", objectKey), nil
}

func setupUserIntegrationApp(t *testing.T) *fiber.App {
	app := fiber.New()

	cfg := config.NewConfig()
	db, err := config.ConnectDB(cfg)
	if err != nil {
		t.Fatalf("failed to connect db: %v", err)
	}

	// migrate user & learner
	if err := db.AutoMigrate(&models.User{}, &models.Learner{}); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}

	// inject DB & mock Minio
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("db", db)
		c.Locals("minio", &MockMinioClient{})
		return c.Next()
	})

	UserRoutes(app)
	return app
}

func TestIntegration_User_CRUD(t *testing.T) {
	app := setupUserIntegrationApp(t)

	// Create
	user := models.User{
		StudentID:         "99990",
		ProfilePictureURL: "",
		FirstName:         "John",
		LastName:          "Doe",
		Gender:            "Male",
		PhoneNumber:       "",
		Balance:           0,
	}
	body, _ := json.Marshal(user)
	req := httptest.NewRequest(http.MethodPost, "/users/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var created models.User
	err := json.NewDecoder(resp.Body).Decode(&created)
	assert.NoError(t, err)
	assert.NotZero(t, created.ID)
	assert.NotNil(t, created.Learner)

	// Get All (admin only)
	req = httptest.NewRequest(http.MethodGet, "/users/", nil)
	resp, _ = app.Test(req)
	// Simulate being ADMIN without TOKEN
	assert.True(t, resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusUnauthorized)

	// Get by ID
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/users/%d", created.ID), nil)
	resp, _ = app.Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	
	// Update
	updatePayload := map[string]interface{}{
		"Name": "Alice Updated",
	}
	body, _ = json.Marshal(updatePayload)
	req = httptest.NewRequest(http.MethodPut, fmt.Sprintf("/users/%d", created.ID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ = app.Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Invalid ID
	req = httptest.NewRequest(http.MethodGet, "/users/abc", nil)
	resp, _ = app.Test(req)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Delete
	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/users/%d", created.ID), nil)
	resp, _ = app.Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Not Found
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/users/%d", created.ID), nil)
	resp, _ = app.Test(req)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}