package handlers

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/a2n2k3p4/tutorium-backend/models"
)

func TestIntegration_ClassSession_CRUD(t *testing.T) {
	// preload
	teacher, _ := createTestTeacher(t)
	class := createTestClass(t, teacher.ID)

	created := createTestClassSession(t, class.ID)

	sessions := getJSONResource[[]models.ClassSession](t, "/class_sessions/", http.StatusOK)
	if len(sessions) == 0 {
		t.Fatalf("expected non-empty class sessions list")
	}

	fetched := getJSONResource[models.ClassSession](t, fmt.Sprintf("/class_sessions/%d", created.ID), http.StatusOK)
	if fetched.ClassID != class.ID {
		t.Fatalf("expected class_id %d, got %d", class.ID, fetched.ClassID)
	}

	jsonRequestExpect(t, http.MethodGet, "/class_sessions/abc", nil, http.StatusBadRequest, nil)

	updatePayload := map[string]any{
		"description": "Updated integration session",
		"price":       1500,
	}
	updateJSONResource(t, fmt.Sprintf("/class_sessions/%d", created.ID), updatePayload, http.StatusOK)

	fetched = getJSONResource[models.ClassSession](t, fmt.Sprintf("/class_sessions/%d", created.ID), http.StatusOK)
	if fetched.Description != "Updated integration session" || fetched.Price != 1500 {
		t.Fatalf("expected updated description and price, got desc=%s price=%f", fetched.Description, fetched.Price)
	}

	deleteJSONResource(t, fmt.Sprintf("/class_sessions/%d", created.ID), http.StatusOK)
	jsonRequestExpect(t, http.MethodGet, fmt.Sprintf("/class_sessions/%d", created.ID), nil, http.StatusNotFound, nil)
}
