package handlers

import (
	"testing"

	"github.com/a2n2k3p4/tutorium-backend/models"
)

func TestIntegration_Teacher_CRUD(t *testing.T) {
	user, _ := createTestUser(t)
	updatedDescription := "Updated bio: now teaching advanced algebra."
	updatedFlagCount := 3

	runCRUDTest(t, crudTestCase[models.Teacher]{
		ResourceName: "teachers",
		BasePath:     "/teachers/",
		Create: func(t *testing.T) models.Teacher {
			created := createTestTeacher(t, user.ID)
			return created
		},
		GetID: func(teach models.Teacher) uint { return teach.ID },
		UpdatePayload: map[string]any{
			"description": updatedDescription,
			"flag_count":  updatedFlagCount,
		},
		AssertUpdated: func(t *testing.T, updated models.Teacher) {
			if updated.Description != updatedDescription || updated.FlagCount != updatedFlagCount {
				t.Fatalf("expected updated description/flag_count %q/%d, got %q/%d", updatedDescription, updatedFlagCount, updated.Description, updated.FlagCount)
			}
		},
	})
}
