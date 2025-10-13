package handlers

import (
	"testing"

	"github.com/a2n2k3p4/tutorium-backend/models"
)

func TestIntegration_Learner_CRUD(t *testing.T) {
	runCRUDTest(t, crudTestCase[models.Learner]{
		ResourceName: "learners",
		BasePath:     "/learners/",
		Create: func(t *testing.T) models.Learner {
			_, learner := createTestUser(t)
			return learner
		},
		GetID: func(l models.Learner) uint { return l.ID },
	})
}
