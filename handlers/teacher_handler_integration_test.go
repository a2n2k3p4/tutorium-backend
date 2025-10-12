package handlers

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/a2n2k3p4/tutorium-backend/models"
)

func TestIntegration_Teacher_CRUD(t *testing.T) {
	created, _ := createTestTeacher(t)

	teachers := getJSONResource[[]models.Teacher](t, "/teachers/", http.StatusOK)
	if len(teachers) == 0 {
		t.Fatalf("expected teachers list to have entries")
	}

	fetched := getJSONResource[models.Teacher](t, fmt.Sprintf("/teachers/%d", created.ID), http.StatusOK)
	if fetched.Email != created.Email {
		t.Fatalf("expected teacher email %s, got %s", created.Email, fetched.Email)
	}

	jsonRequestExpect(t, http.MethodGet, "/teachers/abc", nil, http.StatusBadRequest, nil)

	updatePayload := map[string]any{
		"description": "Updated bio: now teaching advanced algebra.",
		"flag_count":  3,
	}
	updateJSONResource(t, fmt.Sprintf("/teachers/%d", created.ID), updatePayload, http.StatusOK)

	fetched = getJSONResource[models.Teacher](t, fmt.Sprintf("/teachers/%d", created.ID), http.StatusOK)
	if fetched.Description != "Updated bio: now teaching advanced algebra." || fetched.FlagCount != 3 {
		t.Fatalf("expected updated description/flag_count, got %s/%d", fetched.Description, fetched.FlagCount)
	}

	deleteJSONResource(t, fmt.Sprintf("/teachers/%d", created.ID), http.StatusOK)
	jsonRequestExpect(t, http.MethodGet, fmt.Sprintf("/teachers/%d", created.ID), nil, http.StatusNotFound, nil)
}
