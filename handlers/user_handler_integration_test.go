package handlers

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/a2n2k3p4/tutorium-backend/models"
)

func TestIntegration_UserCRUD(t *testing.T) {
	created := createTestUser(t)

	users := getJSONResource[[]models.User](t, "/users/", http.StatusOK)
	if len(users) == 0 {
		t.Fatalf("expected users list to have entries")
	}

	fetched := getJSONResource[models.User](t, fmt.Sprintf("/users/%d", created.ID), http.StatusOK)
	if fetched.StudentID != created.StudentID {
		t.Fatalf("expected student_id %s, got %s", created.StudentID, fetched.StudentID)
	}

	updatePayload := map[string]any{
		"phone_number": "+66111111111",
		"ban_count":    1,
	}
	updateJSONResource(t, fmt.Sprintf("/users/%d", created.ID), updatePayload, http.StatusOK)

	fetched = getJSONResource[models.User](t, fmt.Sprintf("/users/%d", created.ID), http.StatusOK)
	if fetched.PhoneNumber != "+66111111111" || fetched.BanCount != 1 {
		t.Fatalf("update failed, got phone=%s ban=%d", fetched.PhoneNumber, fetched.BanCount)
	}

	deleteJSONResource(t, fmt.Sprintf("/users/%d", created.ID), http.StatusOK)
	jsonRequestExpect(t, http.MethodGet, fmt.Sprintf("/users/%d", created.ID), nil, http.StatusNotFound, nil)
}
