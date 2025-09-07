package handlers

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/a2n2k3p4/tutorium-backend/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

/* ------------------ CreateEnrollment ------------------ */

// 201
func TestCreateEnrollment_OK(t *testing.T) {
	table := "enrollments"
	userID := uint(42)
	learnerID := uint(5)
	classSessionID := uint(10)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, true)(mock)
			ExpInsertReturningID(table, 1)(mock)

			req := jsonBody(models.Enrollment{
				LearnerID:        learnerID,
				ClassSessionID:   classSessionID,
				EnrollmentStatus: "Success",
			})
			*payload = req
			*uID = userID
		},
		http.StatusCreated,
		http.MethodPost,
		"/enrollments/",
	)
}

// 400
func TestCreateEnrollment_BadRequest(t *testing.T) {
	userID := uint(42)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, true)(mock)
			*uID = userID
		},
		http.StatusBadRequest,
		http.MethodPost,
		"/enrollments/",
	)
}

// 500
func TestCreateEnrollment_DBError(t *testing.T) {
	table := "enrollments"
	userID := uint(42)
	learnerID := uint(5)
	classSessionID := uint(10)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, true)(mock)
			ExpInsertError(table, fmt.Errorf("db insert failed"))(mock)
			req := jsonBody(models.Enrollment{
				LearnerID:        learnerID,
				ClassSessionID:   classSessionID,
				EnrollmentStatus: "Success",
			})
			*payload = req
			*uID = userID
		},
		http.StatusInternalServerError,
		http.MethodPost,
		"/enrollments/",
	)
}

/* ------------------ GetEnrollments ------------------ */

// 200
func TestGetEnrollments_OK(t *testing.T) {
	userID := uint(42)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, true)(mock)
			ExpListRows("enrollments", []string{"id"}, []any{1}, []any{2})(mock)
			*uID = userID
		},
		http.StatusOK,
		http.MethodGet,
		"/enrollments/",
	)
}

// 500
func TestGetEnrollments_DBError(t *testing.T) {
	userID := uint(42)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, true)(mock)
			ExpListError("enrollments", fmt.Errorf("select failed"))(mock)
			*uID = userID
		},
		http.StatusInternalServerError,
		http.MethodGet,
		"/enrollments/",
	)
}

/* ------------------ GetEnrollment ------------------ */

// 200
func TestGetEnrollment_OK(t *testing.T) {
	table := "enrollments"
	userID := uint(42)
	enrollmentID := uint(7)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, true)(mock)
			ExpSelectByIDFound(table, enrollmentID, []string{"id"}, []any{enrollmentID})(mock)
			*uID = userID
		},
		http.StatusOK,
		http.MethodGet,
		fmt.Sprintf("/enrollments/%d", enrollmentID),
	)
}

// 404
func TestGetEnrollment_NotFound(t *testing.T) {
	table := "enrollments"
	userID := uint(42)
	enrollmentID := uint(999)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, true)(mock)
			ExpSelectByIDEmpty(table, enrollmentID)(mock)
			*uID = userID
		},
		http.StatusNotFound,
		http.MethodGet,
		fmt.Sprintf("/enrollments/%d", enrollmentID),
	)
}

// 500
func TestGetEnrollment_DBError(t *testing.T) {
	table := "enrollments"
	userID := uint(42)
	enrollmentID := uint(7)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, true)(mock)
			ExpSelectByIDError(table, enrollmentID, fmt.Errorf("select failed"))(mock)
			*uID = userID
		},
		http.StatusInternalServerError,
		http.MethodGet,
		fmt.Sprintf("/enrollments/%d", enrollmentID),
	)
}

// 400
func TestGetEnrollment_BadRequest(t *testing.T) {
	userID := uint(42)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, true)(mock)
			*uID = userID
		},
		http.StatusBadRequest,
		http.MethodGet,
		"/enrollments/not-an-int",
	)
}

/* ------------------ UpdateEnrollment ------------------ */

// 200
func TestUpdateEnrollment_OK(t *testing.T) {
	table := "enrollments"
	preloadTable1 := "learners"
	preloadTable2 := "class_sessions"
	userID := uint(42)
	enrollmentID := uint(1)
	learnerID := uint(5)
	classSessionID := uint(10)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, true)(mock)
			ExpSelectByIDFound(table, enrollmentID,
				[]string{"id", "learner_id", "class_session_id", "enrollment_status"},
				[]any{enrollmentID, learnerID, classSessionID, "pending"},
			)(mock)

			ExpPreloadField(preloadTable1, []string{"id"}, []any{learnerID})(mock)
			ExpPreloadField(preloadTable2, []string{"id"}, []any{classSessionID})(mock)
			ExpUpdateOK(table)(mock)

			req := jsonBody(models.Enrollment{
				LearnerID:        learnerID,
				ClassSessionID:   classSessionID,
				EnrollmentStatus: "Success",
			})
			*payload = req
			*uID = userID
		},
		http.StatusOK,
		http.MethodPut,
		fmt.Sprintf("/enrollments/%d", enrollmentID),
	)
}

// 404
func TestUpdateEnrollment_NotFound(t *testing.T) {
	table := "enrollments"
	userID := uint(42)
	enrollmentID := uint(12345)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, true)(mock)
			ExpSelectByIDEmpty(table, enrollmentID)(mock)
			*uID = userID
		},
		http.StatusNotFound,
		http.MethodPut,
		fmt.Sprintf("/enrollments/%d", enrollmentID),
	)
}

// 500
func TestUpdateEnrollment_DBError(t *testing.T) {
	table := "enrollments"
	userID := uint(42)
	enrollmentID := uint(1)
	learnerID := uint(5)
	classSessionID := uint(10)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, true)(mock)
			ExpSelectByIDFound(table, enrollmentID, []string{"id"}, []any{enrollmentID})(mock)
			ExpUpdateError(table, fmt.Errorf("update failed"))(mock)

			req := jsonBody(models.Enrollment{
				LearnerID:        learnerID,
				ClassSessionID:   classSessionID,
				EnrollmentStatus: "Success",
			})
			*payload = req
			*uID = userID
		},
		http.StatusInternalServerError,
		http.MethodPut,
		fmt.Sprintf("/enrollments/%d", enrollmentID),
	)
}

// 400
func TestUpdateEnrollment_BadRequest(t *testing.T) {
	userID := uint(42)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, true)(mock)
			*uID = userID
		},
		http.StatusBadRequest,
		http.MethodPut,
		"/enrollments/not-an-int",
	)
}

/* ------------------ DeleteEnrollment ------------------ */

// 200 (soft delete)
func TestDeleteEnrollment_OK_SoftDelete(t *testing.T) {
	table := "enrollments"
	userID := uint(42)
	enrollmentID := uint(5)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, true)(mock)
			ExpSelectByIDFound(table, enrollmentID, []string{"id"}, []any{enrollmentID})(mock)
			ExpSoftDeleteOK(table)(mock)
			*uID = userID
		},
		http.StatusOK,
		http.MethodDelete,
		fmt.Sprintf("/enrollments/%d", enrollmentID),
	)
}

// 404
func TestDeleteEnrollment_NotFound(t *testing.T) {
	table := "enrollments"
	userID := uint(42)
	enrollmentID := uint(12345)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, true)(mock)
			ExpSelectByIDEmpty(table, enrollmentID)(mock)
			*uID = userID
		},
		http.StatusNotFound,
		http.MethodDelete,
		fmt.Sprintf("/enrollments/%d", enrollmentID),
	)
}

// 500
func TestDeleteEnrollment_DBError(t *testing.T) {
	table := "enrollments"
	userID := uint(42)
	enrollmentID := uint(5)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, true)(mock)
			ExpSelectByIDFound(table, enrollmentID, []string{"id"}, []any{enrollmentID})(mock)
			ExpSoftDeleteError(table, fmt.Errorf("update failed"))(mock)
			*uID = userID
		},
		http.StatusInternalServerError,
		http.MethodDelete,
		fmt.Sprintf("/enrollments/%d", enrollmentID),
	)
}

// 400
func TestDeleteEnrollment_BadRequest(t *testing.T) {
	userID := uint(42)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, true)(mock)
			*uID = userID
		},
		http.StatusBadRequest,
		http.MethodDelete,
		"/enrollments/not-an-int",
	)
}
