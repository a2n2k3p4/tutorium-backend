package handlers

import (
	"testing"

	"github.com/a2n2k3p4/tutorium-backend/models"
)

func TestIntegration_BanTeacher_CRUD(t *testing.T) {
	user, _ := createTestUser(t)
	teacher := createTestTeacher(t, user.ID)
	updatedDescription := "updated reason"

	runCRUDTest(t, crudTestCase[models.BanDetailsTeacher]{
		ResourceName: "ban teachers",
		BasePath:     "/banteachers/",
		Create: func(t *testing.T) models.BanDetailsTeacher {
			return createTestBanTeacher(t, teacher.ID)
		},
		GetID:         func(b models.BanDetailsTeacher) uint { return b.ID },
		UpdatePayload: map[string]any{"ban_description": updatedDescription},
		AssertUpdated: func(t *testing.T, updated models.BanDetailsTeacher) {
			if updated.BanDescription != updatedDescription {
				t.Fatalf("expected updated description %q, got %q", updatedDescription, updated.BanDescription)
			}
		},
	})
}
