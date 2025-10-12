package handlers

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/a2n2k3p4/tutorium-backend/models"
)

func TestIntegration_BanLearner_CRUD(t *testing.T) {
	// preload
	learner, _ := createTestLearner(t)

	created := createTestBanLearner(t, learner.ID)

	bans := getJSONResource[[]models.BanDetailsLearner](t, "/banlearners/", http.StatusOK)
	if len(bans) == 0 {
		t.Fatalf("expected ban list to contain at least one record")
	}

	fetched := getJSONResource[models.BanDetailsLearner](t, fmt.Sprintf("/banlearners/%d", created.ID), http.StatusOK)
	if fetched.LearnerID != learner.ID {
		t.Fatalf("expected learner_id %d, got %d", learner.ID, fetched.LearnerID)
	}

	jsonRequestExpect(t, http.MethodGet, "/banlearners/abc", nil, http.StatusBadRequest, nil)

	updatePayload := map[string]any{"ban_description": "updated reason"}
	updateJSONResource(t, fmt.Sprintf("/banlearners/%d", created.ID), updatePayload, http.StatusOK)

	fetched = getJSONResource[models.BanDetailsLearner](t, fmt.Sprintf("/banlearners/%d", created.ID), http.StatusOK)
	if fetched.BanDescription != "updated reason" {
		t.Fatalf("expected updated description, got %s", fetched.BanDescription)
	}

	deleteJSONResource(t, fmt.Sprintf("/banlearners/%d", created.ID), http.StatusOK)
	jsonRequestExpect(t, http.MethodGet, fmt.Sprintf("/banlearners/%d", created.ID), nil, http.StatusNotFound, nil)
}
