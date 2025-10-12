package handlers

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/a2n2k3p4/tutorium-backend/models"
)

func TestIntegration_BanTeacher_CRUD(t *testing.T) {
	// preload
	teacher, _ := createTestTeacher(t)

	created := createTestBanTeacher(t, teacher.ID)

	bans := getJSONResource[[]models.BanDetailsTeacher](t, "/banteachers/", http.StatusOK)
	if len(bans) == 0 {
		t.Fatalf("expected bans list to be non-empty")
	}

	fetched := getJSONResource[models.BanDetailsTeacher](t, fmt.Sprintf("/banteachers/%d", created.ID), http.StatusOK)
	if fetched.TeacherID != teacher.ID {
		t.Fatalf("expected teacher_id %d, got %d", teacher.ID, fetched.TeacherID)
	}

	jsonRequestExpect(t, http.MethodGet, "/banteachers/abc", nil, http.StatusBadRequest, nil)

	updatePayload := map[string]any{"ban_description": "updated reason"}
	updateJSONResource(t, fmt.Sprintf("/banteachers/%d", created.ID), updatePayload, http.StatusOK)

	fetched = getJSONResource[models.BanDetailsTeacher](t, fmt.Sprintf("/banteachers/%d", created.ID), http.StatusOK)
	if fetched.BanDescription != "updated reason" {
		t.Fatalf("expected updated description, got %s", fetched.BanDescription)
	}

	deleteJSONResource(t, fmt.Sprintf("/banteachers/%d", created.ID), http.StatusOK)
	jsonRequestExpect(t, http.MethodGet, fmt.Sprintf("/banteachers/%d", created.ID), nil, http.StatusNotFound, nil)
}
