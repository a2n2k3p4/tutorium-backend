package handlers

import (
	"testing"

	"github.com/a2n2k3p4/tutorium-backend/models"
)

func TestIntegration_Report_CRUD(t *testing.T) {
	reporter, _ := createTestUser(t)
	reported, _ := createTestUser(t)
	teacher := createTestTeacher(t, reported.ID)
	class := createTestClass(t, teacher.ID)
	session := createTestClassSession(t, class.ID)
	updatedStatus := "resolved"
	updatedResult := "Case closed by admin"

	runCRUDTest(t, crudTestCase[models.Report]{
		ResourceName: "reports",
		BasePath:     "/reports/",
		Create: func(t *testing.T) models.Report {
			return createTestReport(t, reporter.ID, reported.ID, session.ID)
		},
		GetID: func(r models.Report) uint { return r.ID },
		UpdatePayload: map[string]any{
			"report_status": updatedStatus,
			"report_result": updatedResult,
		},
		AssertUpdated: func(t *testing.T, updated models.Report) {
			if updated.ReportStatus != updatedStatus || updated.ReportResult != updatedResult {
				t.Fatalf("expected report status %q with result %q, got status=%q result=%q", updatedStatus, updatedResult, updated.ReportStatus, updated.ReportResult)
			}
		},
	})
}
