package handlers

import (
	"testing"

	"github.com/a2n2k3p4/tutorium-backend/models"
)

func TestIntegration_ClassSession_CRUD(t *testing.T) {
	user, _ := createTestUser(t)
	teacher := createTestTeacher(t, user.ID)
	class := createTestClass(t, teacher.ID)
	updatedDescription := "Updated integration session"
	updatedPrice := 1500.0

	runCRUDTest(t, crudTestCase[models.ClassSession]{
		ResourceName: "class sessions",
		BasePath:     "/class_sessions/",
		Create: func(t *testing.T) models.ClassSession {
			return createTestClassSession(t, class.ID)
		},
		GetID: func(cs models.ClassSession) uint { return cs.ID },
		UpdatePayload: map[string]any{
			"description": updatedDescription,
			"price":       updatedPrice,
		},
		AssertUpdated: func(t *testing.T, updated models.ClassSession) {
			if updated.Description != updatedDescription || updated.Price != updatedPrice {
				t.Fatalf("expected updated description %q and price %.2f, got desc=%q price=%.2f", updatedDescription, updatedPrice, updated.Description, updated.Price)
			}
		},
	})
}
