package handlers

import (
	"testing"

	"github.com/a2n2k3p4/tutorium-backend/models"
)

func TestIntegration_Class_CRUD(t *testing.T) {
	user, _ := createTestUser(t)
	teacher := createTestTeacher(t, user.ID)
	updatedDescription := "Updated integration class description"
	updatedRating := 4.2

	runCRUDTest(t, crudTestCase[models.Class]{
		ResourceName: "classes",
		BasePath:     "/classes/",
		Create: func(t *testing.T) models.Class {
			return createTestClass(t, teacher.ID)
		},
		GetID: func(c models.Class) uint { return c.ID },
		UpdatePayload: map[string]any{
			"class_description": updatedDescription,
			"rating":            updatedRating,
		},
		AssertUpdated: func(t *testing.T, updated models.Class) {
			if updated.ClassDescription != updatedDescription || updated.Rating != updatedRating {
				t.Fatalf("expected updated description %q and rating %.1f, got desc=%q rating=%.1f", updatedDescription, updatedRating, updated.ClassDescription, updated.Rating)
			}
		},
	})
}
