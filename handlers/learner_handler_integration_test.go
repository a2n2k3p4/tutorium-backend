package handlers

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/a2n2k3p4/tutorium-backend/models"
)

func TestIntegration_Learner_CRUD(t *testing.T) {
	created, user := createTestLearner(t)

	learners := getJSONResource[[]models.Learner](t, "/learners/", http.StatusOK)
	if len(learners) == 0 {
		t.Fatalf("expected learners list to be non-empty")
	}

	fetched := getJSONResource[models.Learner](t, fmt.Sprintf("/learners/%d", created.ID), http.StatusOK)
	if fetched.UserID != user.ID {
		t.Fatalf("expected user_id %d, got %d", user.ID, fetched.UserID)
	}

	jsonRequestExpect(t, http.MethodGet, "/learners/abc", nil, http.StatusBadRequest, nil)

	deleteJSONResource(t, fmt.Sprintf("/learners/%d", created.ID), http.StatusOK)
	jsonRequestExpect(t, http.MethodGet, fmt.Sprintf("/learners/%d", created.ID), nil, http.StatusNotFound, nil)
}
