package handlers

import (
	"testing"

	"github.com/a2n2k3p4/tutorium-backend/models"
)

func TestIntegration_Enrollment_CRUD(t *testing.T) {
	_, learner := createTestUser(t)
	user, _ := createTestUser(t)
	teacher := createTestTeacher(t, user.ID)
	class := createTestClass(t, teacher.ID)
	session := createTestClassSession(t, class.ID)
	updatedStatus := "inactive"

	runCRUDTest(t, crudTestCase[models.Enrollment]{
		ResourceName: "enrollments",
		BasePath:     "/enrollments/",
		Create: func(t *testing.T) models.Enrollment {
			return createTestEnrollment(t, learner.ID, session.ID)
		},
		GetID:         func(e models.Enrollment) uint { return e.ID },
		UpdatePayload: map[string]any{"enrollment_status": updatedStatus},
		AssertUpdated: func(t *testing.T, updated models.Enrollment) {
			if updated.EnrollmentStatus != updatedStatus {
				t.Fatalf("expected enrollment_status %q, got %q", updatedStatus, updated.EnrollmentStatus)
			}
		},
	})
}
