package handlers

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/a2n2k3p4/tutorium-backend/models"
)

func TestIntegration_Review_CRUD(t *testing.T) {
	// preload
	learner, _ := createTestLearner(t)
	teacher, _ := createTestTeacher(t)
	class := createTestClass(t, teacher.ID)

	created := createTestReview(t, learner.ID, class.ID)

	reviews := getJSONResource[[]models.Review](t, "/reviews/", http.StatusOK)
	if len(reviews) == 0 {
		t.Fatalf("expected reviews list to be non-empty")
	}

	fetched := getJSONResource[models.Review](t, fmt.Sprintf("/reviews/%d", created.ID), http.StatusOK)
	if fetched.LearnerID != learner.ID || fetched.ClassID != class.ID {
		t.Fatalf("unexpected learner/class id got %d/%d", fetched.LearnerID, fetched.ClassID)
	}

	jsonRequestExpect(t, http.MethodGet, "/reviews/abc", nil, http.StatusBadRequest, nil)

	updatePayload := map[string]any{
		"rating":  4,
		"comment": "Updated comment: Excellent session!",
	}
	updateJSONResource(t, fmt.Sprintf("/reviews/%d", created.ID), updatePayload, http.StatusOK)

	fetched = getJSONResource[models.Review](t, fmt.Sprintf("/reviews/%d", created.ID), http.StatusOK)
	if fetched.Rating != 4 || fetched.Comment != "Updated comment: Excellent session!" {
		t.Fatalf("expected updated rating/comment, got %d/%s", fetched.Rating, fetched.Comment)
	}

	deleteJSONResource(t, fmt.Sprintf("/reviews/%d", created.ID), http.StatusOK)
	jsonRequestExpect(t, http.MethodGet, fmt.Sprintf("/reviews/%d", created.ID), nil, http.StatusNotFound, nil)
}
