package handlers

import (
	"testing"

	"github.com/a2n2k3p4/tutorium-backend/models"
)

func TestIntegration_BanLearner_CRUD(t *testing.T) {
	_, learner := createTestUser(t)
	updatedDescription := "updated reason"

	runCRUDTest(t, crudTestCase[models.BanDetailsLearner]{
		ResourceName: "ban learners",
		BasePath:     "/banlearners/",
		Create: func(t *testing.T) models.BanDetailsLearner {
			return createTestBanLearner(t, learner.ID)
		},
		GetID:         func(b models.BanDetailsLearner) uint { return b.ID },
		UpdatePayload: map[string]any{"ban_description": updatedDescription},
		AssertUpdated: func(t *testing.T, updated models.BanDetailsLearner) {
			if updated.BanDescription != updatedDescription {
				t.Fatalf("expected updated description %q, got %q", updatedDescription, updated.BanDescription)
			}
		},
	})
}
