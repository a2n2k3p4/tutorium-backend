package handlers

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/a2n2k3p4/tutorium-backend/models"
)

func TestIntegration_Admin_CRUD(t *testing.T) {
	// preload
	user := createTestUser(t)

	payload := map[string]any{"user_id": user.ID}
	created := createJSONResource[models.Admin](t, "/admins/", payload, http.StatusCreated)

	admins := getJSONResource[[]models.Admin](t, "/admins/", http.StatusOK)
	if len(admins) == 0 {
		t.Fatalf("expected at least one admin in list")
	}

	fetched := getJSONResource[models.Admin](t, fmt.Sprintf("/admins/%d", created.ID), http.StatusOK)
	if fetched.UserID != user.ID {
		t.Fatalf("expected admin user_id %d, got %d", user.ID, fetched.UserID)
	}

	jsonRequestExpect(t, http.MethodGet, "/admins/abc", nil, http.StatusBadRequest, nil)

	deleteJSONResource(t, fmt.Sprintf("/admins/%d", created.ID), http.StatusOK)
	jsonRequestExpect(t, http.MethodGet, fmt.Sprintf("/admins/%d", created.ID), nil, http.StatusNotFound, nil)
}
