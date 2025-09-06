package handlers

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/a2n2k3p4/tutorium-backend/models"
)

/* ------------------ CreateReport ------------------ */

// 201
func TestCreateReport_OK(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)

		table := "reports"
		userID := uint(42)
		userReportedID := uint(50)

		ExpAuthUser(userID, false, false, true)(mock)
		ExpInsertReturningID(table, 1)(mock)

		app := setupApp(gdb)

		payload := models.Report{
			ReportUserID:      userID,
			ReportedUserID:    userReportedID,
			ReportType:        "test",
			ReportDescription: "Lorem ipsum dolor sit amet, consectetur adipiscing elit",
			ReportPictureURL:  "",
			ReportDate:        time.Now(),
		}

		resp := runHTTP(t, app, httpInput{
			Method:      http.MethodPost,
			Path:        "/reports/",
			Body:        jsonBody(payload),
			ContentType: "application/json",
			UserID:      &userID,
		})
		wantStatus(t, resp, http.StatusCreated)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

// 400
func TestCreateReport_BadRequest(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)
		userID := uint(42)

		ExpAuthUser(userID, false, false, true)(mock)

		app := setupApp(gdb)
		resp := runHTTP(t, app, httpInput{
			Method: http.MethodPost, Path: "/reports/",
			Body: []byte(`{invalid-json}`), ContentType: "application/json", UserID: &userID,
		})
		wantStatus(t, resp, http.StatusBadRequest)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

// 500
func TestCreateReport_DBError(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)

		table := "reports"
		userID := uint(42)
		userReportedID := uint(50)

		ExpAuthUser(userID, false, false, true)(mock)
		ExpInsertError(table, fmt.Errorf("db insert failed"))(mock)

		payload := models.Report{
			ReportUserID:      userID,
			ReportedUserID:    userReportedID,
			ReportType:        "test",
			ReportDescription: "Lorem ipsum dolor sit amet, consectetur adipiscing elit",
			ReportPictureURL:  "",
			ReportDate:        time.Now(),
		}
		app := setupApp(gdb)
		resp := runHTTP(t, app, httpInput{
			Method: http.MethodPost, Path: "/reports/",
			Body: jsonBody(payload), ContentType: "application/json", UserID: &userID,
		})
		wantStatus(t, resp, http.StatusInternalServerError)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

/* ------------------ GetReports ------------------ */
// 200
func TestGetReports_OK(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)
		userID := uint(42)

		ExpAuthUser(userID, true, false, false)(mock)
		ExpListRows("reports", []string{"id"}, []any{1}, []any{2})(mock)

		app := setupApp(gdb)
		resp := runHTTP(t, app, httpInput{
			Method: http.MethodGet, Path: "/reports/", UserID: &userID,
		})
		wantStatus(t, resp, http.StatusOK)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

// 500
func TestGetReports_DBError(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)
		userID := uint(42)

		ExpAuthUser(userID, true, false, false)(mock)
		ExpListError("reports", fmt.Errorf("select failed"))(mock)

		app := setupApp(gdb)
		resp := runHTTP(t, app, httpInput{
			Method: http.MethodGet, Path: "/reports/", UserID: &userID,
		})
		wantStatus(t, resp, http.StatusInternalServerError)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

/* ------------------ GetReport ------------------ */
// 200
func TestGetReport_OK(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)

		table := "reports"
		userID := uint(42)
		reportID := uint(7)

		ExpAuthUser(userID, true, false, false)(mock)
		ExpSelectByIDFound(table, reportID, []string{"id"}, []any{reportID})(mock)

		app := setupApp(gdb)
		resp := runHTTP(t, app, httpInput{
			Method: http.MethodGet, Path: fmt.Sprintf("/reports/%d", reportID), UserID: &userID,
		})
		wantStatus(t, resp, http.StatusOK)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

// 404
func TestGetReport_NotFound(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)

		table := "reports"
		userID := uint(42)
		reportID := uint(999)

		ExpAuthUser(userID, true, false, false)(mock)
		ExpSelectByIDEmpty(table, reportID)(mock)

		app := setupApp(gdb)
		resp := runHTTP(t, app, httpInput{
			Method: http.MethodGet, Path: fmt.Sprintf("/reports/%d", reportID), UserID: &userID,
		})
		wantStatus(t, resp, http.StatusNotFound)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

// 500
func TestGetReport_DBError(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)

		table := "reports"
		userID := uint(42)
		reportID := uint(7)

		ExpAuthUser(userID, true, false, false)(mock)
		ExpSelectByIDError(table, reportID, fmt.Errorf("select failed"))(mock)

		app := setupApp(gdb)
		resp := runHTTP(t, app, httpInput{
			Method: http.MethodGet, Path: fmt.Sprintf("/reports/%d", reportID), UserID: &userID,
		})
		wantStatus(t, resp, http.StatusInternalServerError)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

// 400
func TestGetReport_BadRequest(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)
		userID := uint(42)

		ExpAuthUser(userID, true, false, false)(mock)

		app := setupApp(gdb)
		resp := runHTTP(t, app, httpInput{
			Method: http.MethodGet, Path: "/reports/not-an-int", UserID: &userID,
		})
		wantStatus(t, resp, http.StatusBadRequest)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

/* ------------------ UpdateReport ------------------ */

// 200
func TestUpdateReport_OK(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)

		table := "reports"
		preloadTable := "users"

		userID := uint(42)
		userReportedID := uint(50)
		reportID := uint(1)

		ExpAuthUser(userID, true, false, false)(mock)
		ExpSelectByIDFound(table, reportID,
			[]string{"id", "report_user_id", "reported_user_id", "report_type", "report_description", "report_picture_url", "report_date"},
			[]any{reportID, userID, userReportedID, "old_type", "Lorem", "", time.Now()},
		)(mock)

		ExpPreloadField(preloadTable, []string{"id"}, []any{userReportedID})(mock)
		ExpPreloadField(preloadTable, []string{"id"}, []any{userID})(mock)

		ExpUpdateOK(table)(mock)

		app := setupApp(gdb)
		payload := models.Report{
			ReportUserID:      userID,
			ReportedUserID:    userReportedID,
			ReportType:        "test",
			ReportDescription: "Lorem ipsum dolor sit amet, consectetur adipiscing elit",
			ReportPictureURL:  "",
			ReportDate:        time.Now(),
		}
		resp := runHTTP(t, app, httpInput{
			Method: http.MethodPut, Path: fmt.Sprintf("/reports/%d", reportID),
			Body: jsonBody(payload), ContentType: "application/json", UserID: &userID,
		})
		wantStatus(t, resp, http.StatusOK)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

// 404
func TestUpdateReport_NotFound(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)

		table := "reports"
		userID := uint(42)
		reportID := uint(12345)

		ExpAuthUser(userID, true, false, false)(mock)
		ExpSelectByIDEmpty(table, reportID)(mock)

		app := setupApp(gdb)
		resp := runHTTP(t, app, httpInput{
			Method: http.MethodPut, Path: fmt.Sprintf("/reports/%d", reportID), UserID: &userID,
		})
		wantStatus(t, resp, http.StatusNotFound)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

// 500
func TestUpdateReport_DBError(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)

		table := "reports"

		userID := uint(42)
		userReportedID := uint(50)
		reportID := uint(1)

		ExpAuthUser(userID, true, false, false)(mock)
		ExpSelectByIDFound(table, reportID, []string{"id"}, []any{reportID})(mock)
		ExpUpdateError(table, fmt.Errorf("update failed"))(mock)

		app := setupApp(gdb)
		payload := models.Report{
			ReportUserID:      userID,
			ReportedUserID:    userReportedID,
			ReportType:        "test",
			ReportDescription: "Lorem ipsum dolor sit amet, consectetur adipiscing elit",
			ReportPictureURL:  "",
			ReportDate:        time.Now(),
		}
		resp := runHTTP(t, app, httpInput{
			Method: http.MethodPut, Path: fmt.Sprintf("/reports/%d", reportID),
			Body: jsonBody(payload), ContentType: "application/json", UserID: &userID,
		})
		wantStatus(t, resp, http.StatusInternalServerError)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

// 400
func TestUpdateReport_BadRequest(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)
		userID := uint(42)

		ExpAuthUser(userID, true, false, false)(mock)

		app := setupApp(gdb)
		resp := runHTTP(t, app, httpInput{
			Method: http.MethodPut, Path: "/reports/not-an-int", UserID: &userID,
		})
		wantStatus(t, resp, http.StatusBadRequest)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

/* ------------------ DeleteReport ------------------ */

// 200
func TestDeleteReport_OK_SoftDelete(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)

		table := "reports"
		userID := uint(42)
		reportID := uint(5)

		ExpAuthUser(userID, true, false, false)(mock)
		ExpSelectByIDFound(table, reportID, []string{"id"}, []any{reportID})(mock)
		ExpSoftDeleteOK(table)(mock)

		app := setupApp(gdb)
		resp := runHTTP(t, app, httpInput{
			Method: http.MethodDelete, Path: fmt.Sprintf("/reports/%d", reportID), UserID: &userID,
		})
		wantStatus(t, resp, http.StatusOK)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

// 404
func TestDeleteReport_NotFound(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)

		table := "reports"
		userID := uint(42)
		reportID := uint(12345)

		ExpAuthUser(userID, true, false, false)(mock)
		ExpSelectByIDEmpty(table, reportID)(mock)

		app := setupApp(gdb)
		resp := runHTTP(t, app, httpInput{
			Method: http.MethodDelete, Path: fmt.Sprintf("/reports/%d", reportID), UserID: &userID,
		})
		wantStatus(t, resp, http.StatusNotFound)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

// 500
func TestDeleteReport_DBError(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)

		table := "reports"
		userID := uint(42)
		reportID := uint(5)

		ExpAuthUser(userID, true, false, false)(mock)
		ExpSelectByIDFound(table, reportID, []string{"id"}, []any{reportID})(mock)
		ExpSoftDeleteError(table, fmt.Errorf("update failed"))(mock)

		app := setupApp(gdb)
		resp := runHTTP(t, app, httpInput{
			Method: http.MethodDelete, Path: fmt.Sprintf("/reports/%d", reportID), UserID: &userID,
		})
		wantStatus(t, resp, http.StatusInternalServerError)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

// 400
func TestDeleteReport_BadRequest(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)
		userID := uint(42)

		ExpAuthUser(userID, true, false, false)(mock)

		app := setupApp(gdb)
		resp := runHTTP(t, app, httpInput{
			Method: http.MethodDelete, Path: "/reports/not-an-int", UserID: &userID,
		})
		wantStatus(t, resp, http.StatusBadRequest)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}
