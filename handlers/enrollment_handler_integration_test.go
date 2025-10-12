package handlers

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/a2n2k3p4/tutorium-backend/models"
)

func TestIntegration_Enrollment_CRUD(t *testing.T) {
	// preload
	learner, _ := createTestLearner(t)
	teacher, _ := createTestTeacher(t)
	class := createTestClass(t, teacher.ID)
	session := createTestClassSession(t, class.ID)

	created := createTestEnrollment(t, learner.ID, session.ID)

	enrollments := getJSONResource[[]models.Enrollment](t, "/enrollments/", http.StatusOK)
	if len(enrollments) == 0 {
		t.Fatalf("expected enrollments list to be non-empty")
	}

	fetched := getJSONResource[models.Enrollment](t, fmt.Sprintf("/enrollments/%d", created.ID), http.StatusOK)
	if fetched.LearnerID != learner.ID || fetched.ClassSessionID != session.ID {
		t.Fatalf("expected learner_id %d and class_session_id %d, got %d/%d", learner.ID, session.ID, fetched.LearnerID, fetched.ClassSessionID)
	}

	jsonRequestExpect(t, http.MethodGet, "/enrollments/abc", nil, http.StatusBadRequest, nil)

	updatePayload := map[string]any{"enrollment_status": "inactive"}
	updateJSONResource(t, fmt.Sprintf("/enrollments/%d", created.ID), updatePayload, http.StatusOK)

	fetched = getJSONResource[models.Enrollment](t, fmt.Sprintf("/enrollments/%d", created.ID), http.StatusOK)
	if fetched.EnrollmentStatus != "inactive" {
		t.Fatalf("expected enrollment_status inactive, got %s", fetched.EnrollmentStatus)
	}

	deleteJSONResource(t, fmt.Sprintf("/enrollments/%d", created.ID), http.StatusOK)
	jsonRequestExpect(t, http.MethodGet, fmt.Sprintf("/enrollments/%d", created.ID), nil, http.StatusNotFound, nil)
}
