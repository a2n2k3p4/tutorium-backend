package handlers

import (
	"testing"

	"github.com/a2n2k3p4/tutorium-backend/models"
)

func TestIntegration_Review_CRUD(t *testing.T) {
	_, learner := createTestUser(t)
	user, _ := createTestUser(t)
	teacher := createTestTeacher(t, user.ID)
	class := createTestClass(t, teacher.ID)
	updatedRating := 4
	updatedComment := "Updated comment: Excellent session!"

	runCRUDTest(t, crudTestCase[models.Review]{
		ResourceName: "reviews",
		BasePath:     "/reviews/",
		Create: func(t *testing.T) models.Review {
			return createTestReview(t, learner.ID, class.ID)
		},
		GetID: func(r models.Review) uint { return r.ID },
		UpdatePayload: map[string]any{
			"rating":  updatedRating,
			"comment": updatedComment,
		},
		AssertUpdated: func(t *testing.T, updated models.Review) {
			if updated.Rating != updatedRating || updated.Comment != updatedComment {
				t.Fatalf("expected updated rating/comment %d/%q, got %d/%q", updatedRating, updatedComment, updated.Rating, updated.Comment)
			}
		},
	})
}
