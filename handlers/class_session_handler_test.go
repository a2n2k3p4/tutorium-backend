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

/* ------------------ CreateClassSession ------------------ */

// 201
func TestCreateClassSession_OK(t *testing.T) {
	table := "class_sessions"
	userID := uint(42)
	classID := uint(50)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, true, false)(mock)
			ExpListRows("classes", []string{"teacher_id"}, []any{classID})(mock)
			ExpInsertReturningID(table, 1)(mock)

			req := jsonBody(models.ClassSession{
				ClassID:            classID,
				Description:        "Lorem Ipsum",
				LearnerLimit:       40,
				EnrollmentDeadline: time.Now().Add(72 * time.Hour),
				ClassStart:         time.Now().Add(108 * time.Hour),
				ClassFinish:        time.Now().Add(110 * time.Hour),
				ClassStatus:        "Available",
				ClassURL:           "https://meet.jit.si/KUtutorium-test",
			})
			*payload = req
			*uID = userID
		},
		http.StatusCreated,
		http.MethodPost,
		"/class_sessions/",
	)
}

// 400
func TestCreateClassSession_BadRequest(t *testing.T) {
	userID := uint(42)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, true, false)(mock)
			*uID = userID
		},
		http.StatusBadRequest,
		http.MethodPost,
		"/class_sessions/",
	)
}

// 500
func TestCreateClassSession_DBError(t *testing.T) {
	table := "class_sessions"
	userID := uint(42)
	classID := uint(50)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, true, false)(mock)
			ExpListRows("classes", []string{"teacher_id"}, []any{classID})(mock)
			ExpInsertError(table, fmt.Errorf("db insert failed"))(mock)

			req := jsonBody(models.ClassSession{
				ClassID:            classID,
				Description:        "Lorem Ipsum",
				LearnerLimit:       40,
				EnrollmentDeadline: time.Now().Add(72 * time.Hour),
				ClassStart:         time.Now().Add(108 * time.Hour),
				ClassFinish:        time.Now().Add(110 * time.Hour),
				ClassStatus:        "Available",
				ClassURL:           "https://meet.jit.si/KUtutorium-test",
			})
			*payload = req
			*uID = userID
		},
		http.StatusInternalServerError,
		http.MethodPost,
		"/class_sessions/",
	)
}

/* ------------------ GetClassSessions ------------------ */

// 200
func TestGetClassSessions_OK(t *testing.T) {
	userID := uint(42)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			ExpListRows("class_sessions", []string{"id"}, []any{1}, []any{2})(mock)
			*uID = userID
		},
		http.StatusOK,
		http.MethodGet,
		"/class_sessions/",
	)
}

// 500
func TestGetClassSessions_DBError(t *testing.T) {
	userID := uint(42)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			ExpListError("class_sessions", fmt.Errorf("select failed"))(mock)
			*uID = userID
		},
		http.StatusInternalServerError,
		http.MethodGet,
		"/class_sessions/",
	)
}

/* ------------------ GetClassSession ------------------ */

// 200
func TestGetClassSession_OK(t *testing.T) {
	table := "class_sessions"
	userID := uint(42)
	classSessionID := uint(7)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			ExpSelectByIDFound(table, classSessionID, []string{"id"}, []any{classSessionID})(mock)
			*uID = userID
		},
		http.StatusOK,
		http.MethodGet,
		fmt.Sprintf("/class_sessions/%d", classSessionID),
	)
}

// 404
func TestGetClassSession_NotFound(t *testing.T) {
	table := "class_sessions"
	userID := uint(42)
	classSessionID := uint(999)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			ExpSelectByIDEmpty(table, classSessionID)(mock)
			*uID = userID
		},
		http.StatusNotFound,
		http.MethodGet,
		fmt.Sprintf("/class_sessions/%d", classSessionID),
	)
}

// 500
func TestGetClassSession_DBError(t *testing.T) {
	table := "class_sessions"
	userID := uint(42)
	classSessionID := uint(7)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			ExpSelectByIDError(table, classSessionID, fmt.Errorf("select failed"))(mock)
			*uID = userID
		},
		http.StatusInternalServerError,
		http.MethodGet,
		fmt.Sprintf("/class_sessions/%d", classSessionID),
	)
}

// 400
func TestGetClassSession_BadRequest(t *testing.T) {
	userID := uint(42)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			*uID = userID
		},
		http.StatusBadRequest,
		http.MethodGet,
		"/class_sessions/not-an-int",
	)
}

/* ------------------ UpdateClassSession ------------------ */

// 200
func TestUpdateClassSession_OK(t *testing.T) {
	table := "class_sessions"
	userID := uint(42)
	classID := uint(50)
	classSessionID := uint(1)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, true, false)(mock)
			ExpSelectByIDFound(table, classSessionID,
				[]string{"id", "class_id", "description", "learner_limit", "enrollment_deadline", "class_start", "class_finish", "class_status", "class_url"},
				[]any{classSessionID, classID, "Lorem", 40, time.Now().Add(72 * time.Hour), time.Now().Add(108 * time.Hour), time.Now().Add(110 * time.Hour), "pending", "https://meet.jit.si/KUtutorium-test"},
			)(mock)

			ExpPreloadField("classes", []string{"id"}, []any{classID})(mock)
			ExpUpdateOK(table)(mock)

			req := jsonBody(models.ClassSession{
				ClassID:            classID,
				Description:        "Lorem Ipsum",
				LearnerLimit:       40,
				EnrollmentDeadline: time.Now().Add(72 * time.Hour),
				ClassStart:         time.Now().Add(108 * time.Hour),
				ClassFinish:        time.Now().Add(110 * time.Hour),
				ClassStatus:        "Available",
				ClassURL:           "https://meet.jit.si/KUtutorium-test",
			})
			*payload = req
			*uID = userID
		},
		http.StatusOK,
		http.MethodPut,
		fmt.Sprintf("/class_sessions/%d", classSessionID),
	)
}

// 404
func TestUpdateClassSession_NotFound(t *testing.T) {
	table := "class_sessions"
	userID := uint(42)
	classSessionID := uint(12345)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, true, false)(mock)
			ExpSelectByIDEmpty(table, classSessionID)(mock)
			*uID = userID
		},
		http.StatusNotFound,
		http.MethodPut,
		fmt.Sprintf("/class_sessions/%d", classSessionID),
	)
}

// 500
func TestUpdateClassSession_DBError(t *testing.T) {
	table := "class_sessions"
	userID := uint(42)
	classID := uint(50)
	classSessionID := uint(1)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, true, false)(mock)
			ExpSelectByIDFound(table, classSessionID, []string{"id"}, []any{classSessionID})(mock)
			ExpUpdateError(table, fmt.Errorf("update failed"))(mock)

			req := jsonBody(models.ClassSession{
				ClassID:            classID,
				Description:        "Lorem Ipsum",
				LearnerLimit:       40,
				EnrollmentDeadline: time.Now().Add(72 * time.Hour),
				ClassStart:         time.Now().Add(108 * time.Hour),
				ClassFinish:        time.Now().Add(110 * time.Hour),
				ClassStatus:        "Available",
				ClassURL:           "https://meet.jit.si/KUtutorium-test",
			})
			*payload = req
			*uID = userID
		},
		http.StatusInternalServerError,
		http.MethodPut,
		fmt.Sprintf("/class_sessions/%d", classSessionID),
	)
}

// 400
func TestUpdateClassSession_BadRequest(t *testing.T) {
	userID := uint(42)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, true, false)(mock)
			*uID = userID
		},
		http.StatusBadRequest,
		http.MethodPut,
		"/class_sessions/not-an-int",
	)
}

/* ------------------ DeleteClassSession ------------------ */

// 200
func TestDeleteClassSession_OK_SoftDelete(t *testing.T) {
	table := "class_sessions"
	userID := uint(42)
	classSessionID := uint(5)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, true, false)(mock)
			ExpSelectByIDFound(table, classSessionID, []string{"id"}, []any{classSessionID})(mock)
			ExpSoftDeleteOK(table)(mock)
			*uID = userID
		},
		http.StatusOK,
		http.MethodDelete,
		fmt.Sprintf("/class_sessions/%d", classSessionID),
	)
}

// 404
func TestDeleteClassSession_NotFound(t *testing.T) {
	table := "class_sessions"
	userID := uint(42)
	classSessionID := uint(12345)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, true, false)(mock)
			ExpSelectByIDEmpty(table, classSessionID)(mock)
			*uID = userID
		},
		http.StatusNotFound,
		http.MethodDelete,
		fmt.Sprintf("/class_sessions/%d", classSessionID),
	)
}

// 500
func TestDeleteClassSession_DBError(t *testing.T) {
	table := "class_sessions"
	userID := uint(42)
	classSessionID := uint(5)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, true, false)(mock)
			ExpSelectByIDFound(table, classSessionID, []string{"id"}, []any{classSessionID})(mock)
			ExpSoftDeleteError(table, fmt.Errorf("update failed"))(mock)
			*uID = userID
		},
		http.StatusInternalServerError,
		http.MethodDelete,
		fmt.Sprintf("/class_sessions/%d", classSessionID),
	)
}

// 400
func TestDeleteClassSession_BadRequest(t *testing.T) {
	userID := uint(42)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, true, false)(mock)
			*uID = userID
		},
		http.StatusBadRequest,
		http.MethodDelete,
		"/class_sessions/not-an-int",
	)
}
