package handlers

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/a2n2k3p4/tutorium-backend/models"
)

func TestIntegration_Report_CRUD(t *testing.T) {
	_, reporter := createTestLearner(t)
	teacher, reported := createTestTeacher(t)
	class := createTestClass(t, teacher.ID)
	session := createTestClassSession(t, class.ID)

	created := createTestReport(t, reporter.ID, reported.ID, session.ID)

	reports := getJSONResource[[]models.Report](t, "/reports/", http.StatusOK)
	if len(reports) == 0 {
		t.Fatalf("expected reports list to be non-empty")
	}

	fetched := getJSONResource[models.Report](t, fmt.Sprintf("/reports/%d", created.ID), http.StatusOK)
	if fetched.ReportUserID != reporter.ID || fetched.ReportedUserID != reported.ID {
		t.Fatalf("unexpected reporter/reported ids got %d/%d", fetched.ReportUserID, fetched.ReportedUserID)
	}

	jsonRequestExpect(t, http.MethodGet, "/reports/abc", nil, http.StatusBadRequest, nil)

	updatePayload := map[string]any{
		"report_status": "resolved",
		"report_result": "Case closed by admin",
	}
	updateJSONResource(t, fmt.Sprintf("/reports/%d", created.ID), updatePayload, http.StatusOK)

	fetched = getJSONResource[models.Report](t, fmt.Sprintf("/reports/%d", created.ID), http.StatusOK)
	if fetched.ReportStatus != "resolved" || fetched.ReportResult != "Case closed by admin" {
		t.Fatalf("expected report resolved with result, got status=%s result=%s", fetched.ReportStatus, fetched.ReportResult)
	}

	deleteJSONResource(t, fmt.Sprintf("/reports/%d", created.ID), http.StatusOK)
	jsonRequestExpect(t, http.MethodGet, fmt.Sprintf("/reports/%d", created.ID), nil, http.StatusNotFound, nil)
}
