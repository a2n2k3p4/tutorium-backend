package handlers

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/a2n2k3p4/tutorium-backend/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

/* ------------------ CreateReport ------------------ */

// 201
func TestCreateReport_OK(t *testing.T) {
	table := "reports"
	userID := uint(42)
	userReportedID := uint(50)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, true)(mock)
			ExpInsertReturningID(table, 1)(mock)

			req := jsonBody(models.Report{
				ReportUserID:      userID,
				ReportedUserID:    userReportedID,
				ReportType:        "test",
				ReportDescription: "Lorem ipsum dolor sit amet, consectetur adipiscing elit",
				ReportPictureURL:  "",
				ReportDate:        time.Now(),
			})
			*payload = req
			*uID = userID
		},
		http.StatusCreated,
		http.MethodPost,
		"/reports/",
	)
}

// 400
func TestCreateReport_BadRequest(t *testing.T) {
	userID := uint(42)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, true)(mock)
			*uID = userID
		},
		http.StatusBadRequest,
		http.MethodPost,
		"/reports/",
	)
}

// 500
func TestCreateReport_DBError(t *testing.T) {
	table := "reports"
	userID := uint(42)
	userReportedID := uint(50)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, true)(mock)
			ExpInsertError(table, fmt.Errorf("db insert failed"))(mock)

			req := jsonBody(models.Report{
				ReportUserID:      userID,
				ReportedUserID:    userReportedID,
				ReportType:        "test",
				ReportDescription: "Lorem ipsum dolor sit amet, consectetur adipiscing elit",
				ReportPictureURL:  "",
				ReportDate:        time.Now(),
			})
			*payload = req
			*uID = userID
		},
		http.StatusInternalServerError,
		http.MethodPost,
		"/reports/",
	)
}

/* ------------------ GetReports ------------------ */

// 200
func TestGetReports_OK(t *testing.T) {
	userID := uint(42)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, true, false, false)(mock)
			ExpListRows("reports", []string{"id"}, []any{1}, []any{2})(mock)
			*uID = userID
		},
		http.StatusOK,
		http.MethodGet,
		"/reports/",
	)
}

// 500
func TestGetReports_DBError(t *testing.T) {
	userID := uint(42)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, true, false, false)(mock)
			ExpListError("reports", fmt.Errorf("select failed"))(mock)
			*uID = userID
		},
		http.StatusInternalServerError,
		http.MethodGet,
		"/reports/",
	)
}

/* ------------------ GetReport ------------------ */

// 200
func TestGetReport_OK(t *testing.T) {
	table := "reports"
	userID := uint(42)
	reportID := uint(7)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, true, false, false)(mock)
			ExpSelectByIDFound(table, reportID, []string{"id"}, []any{reportID})(mock)
			*uID = userID
		},
		http.StatusOK,
		http.MethodGet,
		fmt.Sprintf("/reports/%d", reportID),
	)
}

// 404
func TestGetReport_NotFound(t *testing.T) {
	table := "reports"
	userID := uint(42)
	reportID := uint(999)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, true, false, false)(mock)
			ExpSelectByIDEmpty(table, reportID)(mock)
			*uID = userID
		},
		http.StatusNotFound,
		http.MethodGet,
		fmt.Sprintf("/reports/%d", reportID),
	)
}

// 500
func TestGetReport_DBError(t *testing.T) {
	table := "reports"
	userID := uint(42)
	reportID := uint(7)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, true, false, false)(mock)
			ExpSelectByIDError(table, reportID, fmt.Errorf("select failed"))(mock)
			*uID = userID
		},
		http.StatusInternalServerError,
		http.MethodGet,
		fmt.Sprintf("/reports/%d", reportID),
	)
}

// 400
func TestGetReport_BadRequest(t *testing.T) {
	userID := uint(42)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, true, false, false)(mock)
			*uID = userID
		},
		http.StatusBadRequest,
		http.MethodGet,
		"/reports/not-an-int",
	)
}

/* ------------------ UpdateReport ------------------ */

// 200
func TestUpdateReport_OK(t *testing.T) {
	table := "reports"
	preloadTable := "users"
	userID := uint(42)
	userReportedID := uint(50)
	reportID := uint(1)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, true, false, false)(mock)
			ExpSelectByIDFound(table, reportID,
				[]string{"id", "report_user_id", "reported_user_id", "report_type", "report_description", "report_picture_url", "report_date"},
				[]any{reportID, userID, userReportedID, "old_type", "Lorem", "", time.Now()},
			)(mock)

			ExpPreloadField(preloadTable, []string{"id"}, []any{userReportedID})(mock)
			ExpPreloadField(preloadTable, []string{"id"}, []any{userID})(mock)
			ExpUpdateOK(table)(mock)

			req := jsonBody(models.Report{
				ReportUserID:      userID,
				ReportedUserID:    userReportedID,
				ReportType:        "test",
				ReportDescription: "Lorem ipsum dolor sit amet, consectetur adipiscing elit",
				ReportPictureURL:  "",
				ReportDate:        time.Now(),
			})
			*payload = req
			*uID = userID
		},
		http.StatusOK,
		http.MethodPut,
		fmt.Sprintf("/reports/%d", reportID),
	)
}

// 404
func TestUpdateReport_NotFound(t *testing.T) {
	table := "reports"
	userID := uint(42)
	reportID := uint(12345)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, true, false, false)(mock)
			ExpSelectByIDEmpty(table, reportID)(mock)
			*uID = userID
		},
		http.StatusNotFound,
		http.MethodPut,
		fmt.Sprintf("/reports/%d", reportID),
	)
}

// 500
func TestUpdateReport_DBError(t *testing.T) {
	table := "reports"
	userID := uint(42)
	userReportedID := uint(50)
	reportID := uint(1)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, true, false, false)(mock)
			ExpSelectByIDFound(table, reportID, []string{"id"}, []any{reportID})(mock)
			ExpUpdateError(table, fmt.Errorf("update failed"))(mock)

			req := jsonBody(models.Report{
				ReportUserID:      userID,
				ReportedUserID:    userReportedID,
				ReportType:        "test",
				ReportDescription: "Lorem ipsum dolor sit amet, consectetur adipiscing elit",
				ReportPictureURL:  "",
				ReportDate:        time.Now(),
			})
			*payload = req
			*uID = userID
		},
		http.StatusInternalServerError,
		http.MethodPut,
		fmt.Sprintf("/reports/%d", reportID),
	)
}

// 400
func TestUpdateReport_BadRequest(t *testing.T) {
	userID := uint(42)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, true, false, false)(mock)
			*uID = userID
		},
		http.StatusBadRequest,
		http.MethodPut,
		"/reports/not-an-int",
	)
}

/* ------------------ DeleteReport ------------------ */

// 200
func TestDeleteReport_OK_SoftDelete(t *testing.T) {
	table := "reports"
	userID := uint(42)
	reportID := uint(5)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, true, false, false)(mock)
			ExpSelectByIDFound(table, reportID, []string{"id"}, []any{reportID})(mock)
			ExpSoftDeleteOK(table)(mock)
			*uID = userID
		},
		http.StatusOK,
		http.MethodDelete,
		fmt.Sprintf("/reports/%d", reportID),
	)
}

// 404
func TestDeleteReport_NotFound(t *testing.T) {
	table := "reports"
	userID := uint(42)
	reportID := uint(12345)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, true, false, false)(mock)
			ExpSelectByIDEmpty(table, reportID)(mock)
			*uID = userID
		},
		http.StatusNotFound,
		http.MethodDelete,
		fmt.Sprintf("/reports/%d", reportID),
	)
}

// 500
func TestDeleteReport_DBError(t *testing.T) {
	table := "reports"
	userID := uint(42)
	reportID := uint(5)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, true, false, false)(mock)
			ExpSelectByIDFound(table, reportID, []string{"id"}, []any{reportID})(mock)
			ExpSoftDeleteError(table, fmt.Errorf("update failed"))(mock)
			*uID = userID
		},
		http.StatusInternalServerError,
		http.MethodDelete,
		fmt.Sprintf("/reports/%d", reportID),
	)
}

// 400
func TestDeleteReport_BadRequest(t *testing.T) {
	userID := uint(42)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, true, false, false)(mock)
			*uID = userID
		},
		http.StatusBadRequest,
		http.MethodDelete,
		"/reports/not-an-int",
	)
}
