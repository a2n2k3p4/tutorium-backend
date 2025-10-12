package handlers

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/a2n2k3p4/tutorium-backend/models"
)

func TestIntegration_Notification_CRUD(t *testing.T) {
	// preload
	user := createTestUser(t)

	created := createTestNotification(t, user.ID)

	notifications := getJSONResource[[]models.Notification](t, "/notifications/", http.StatusOK)
	if len(notifications) == 0 {
		t.Fatalf("expected notifications list to be non-empty")
	}

	fetched := getJSONResource[models.Notification](t, fmt.Sprintf("/notifications/%d", created.ID), http.StatusOK)
	if fetched.UserID != user.ID {
		t.Fatalf("expected user_id %d, got %d", user.ID, fetched.UserID)
	}

	jsonRequestExpect(t, http.MethodGet, "/notifications/abc", nil, http.StatusBadRequest, nil)

	updatePayload := map[string]any{
		"notification_description": "Updated integration notification",
		"read_flag":                true,
	}
	updateJSONResource(t, fmt.Sprintf("/notifications/%d", created.ID), updatePayload, http.StatusOK)

	fetched = getJSONResource[models.Notification](t, fmt.Sprintf("/notifications/%d", created.ID), http.StatusOK)
	if fetched.NotificationDescription != "Updated integration notification" || !fetched.ReadFlag {
		t.Fatalf("expected updated notification to be marked read with new description")
	}

	deleteJSONResource(t, fmt.Sprintf("/notifications/%d", created.ID), http.StatusOK)
	jsonRequestExpect(t, http.MethodGet, fmt.Sprintf("/notifications/%d", created.ID), nil, http.StatusNotFound, nil)
}
