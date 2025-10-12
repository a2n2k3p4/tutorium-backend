package handlers

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/a2n2k3p4/tutorium-backend/models"
)

func TestIntegration_Class_CRUD(t *testing.T) {
	// preload
	teacher, _ := createTestTeacher(t)

	created := createTestClass(t, teacher.ID)

	class := getJSONResource[[]models.Class](t, "/classes/", http.StatusOK)
	if len(class) == 0 {
		t.Fatalf("expected classes list to be non-empty")
	}

	fetched := getJSONResource[models.Class](t, fmt.Sprintf("/classes/%d", created.ID), http.StatusOK)
	if fetched.ClassName != created.ClassName {
		t.Fatalf("expected class name %s, got %s", created.ClassName, fetched.ClassName)
	}

	jsonRequestExpect(t, http.MethodGet, "/classes/abc", nil, http.StatusBadRequest, nil)

	updatePayload := map[string]any{
		"class_description": "Updated integration class description",
		"rating":            4.2,
	}
	updateJSONResource(t, fmt.Sprintf("/classes/%d", created.ID), updatePayload, http.StatusOK)

	fetched = getJSONResource[models.Class](t, fmt.Sprintf("/classes/%d", created.ID), http.StatusOK)
	if fetched.ClassDescription != "Updated integration class description" || fetched.Rating != 4.2 {
		t.Fatalf("expected updated description and rating, got desc=%s rating=%f", fetched.ClassDescription, fetched.Rating)
	}

	deleteJSONResource(t, fmt.Sprintf("/classes/%d", created.ID), http.StatusOK)
	jsonRequestExpect(t, http.MethodGet, fmt.Sprintf("/classes/%d", created.ID), nil, http.StatusNotFound, nil)
}
