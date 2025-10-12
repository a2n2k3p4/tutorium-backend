package handlers

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/a2n2k3p4/tutorium-backend/models"
)

func TestIntegration_ClassCategory_CRUD(t *testing.T) {
	createPayload := map[string]any{"class_category": "Integration Category"}
	created := createJSONResource[models.ClassCategory](t, "/class_categories/", createPayload, http.StatusCreated)

	cat := getJSONResource[[]models.ClassCategory](t, "/class_categories/", http.StatusOK)
	if len(cat) == 0 {
		t.Fatalf("expected Categories list to be non-empty")
	}

	fetched := getJSONResource[models.ClassCategory](t, fmt.Sprintf("/class_categories/%d", created.ID), http.StatusOK)
	if fetched.ClassCategory != "Integration Category" {
		t.Fatalf("expected Integration Category, got %s", fetched.ClassCategory)
	}

	jsonRequestExpect(t, http.MethodGet, "/class_categories/abc", nil, http.StatusBadRequest, nil)

	updatePayload := map[string]any{"class_category": "Updated Integration Category"}
	updateJSONResource(t, fmt.Sprintf("/class_categories/%d", created.ID), updatePayload, http.StatusOK)

	fetched = getJSONResource[models.ClassCategory](t, fmt.Sprintf("/class_categories/%d", created.ID), http.StatusOK)
	if fetched.ClassCategory != "Updated Integration Category" {
		t.Fatalf("expected updated category name, got %s", fetched.ClassCategory)
	}

	deleteJSONResource(t, fmt.Sprintf("/class_categories/%d", created.ID), http.StatusOK)
	jsonRequestExpect(t, http.MethodGet, fmt.Sprintf("/class_categories/%d", created.ID), nil, http.StatusNotFound, nil)
}
