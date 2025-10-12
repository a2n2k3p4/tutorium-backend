package handlers

import (
	"testing"

	"github.com/a2n2k3p4/tutorium-backend/models"
)

func TestIntegration_Notification_CRUD(t *testing.T) {
	user, _ := createTestUser(t)
	updatedDescription := "Updated integration notification"
	updatedReadFlag := true

	runCRUDTest(t, crudTestCase[models.Notification]{
		ResourceName: "notifications",
		BasePath:     "/notifications/",
		Create: func(t *testing.T) models.Notification {
			return createTestNotification(t, user.ID)
		},
		GetID: func(n models.Notification) uint { return n.ID },
		UpdatePayload: map[string]any{
			"notification_description": updatedDescription,
			"read_flag":                updatedReadFlag,
		},
		AssertUpdated: func(t *testing.T, updated models.Notification) {
			if updated.NotificationDescription != updatedDescription || updated.ReadFlag != updatedReadFlag {
				t.Fatalf("expected notification description %q and read_flag %t, got %q/%t", updatedDescription, updatedReadFlag, updated.NotificationDescription, updated.ReadFlag)
			}
		},
	})
}
