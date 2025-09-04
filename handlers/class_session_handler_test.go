package handlers

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/a2n2k3p4/tutorium-backend/models"
)

/* ------------------ CreateClassSession ------------------ */

// 201
func TestCreateClassSession_OK(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)

	table := "class_sessions"
	userID := uint(42)
	classID := uint(50)

	ExpAuthUser(userID, false, true, false)(mock)
	ExpInsertReturningID(table, 1)(mock)

	app := setupApp(gdb)

	payload := models.ClassSession{
		ClassID:            classID,
		Description:        "Lorem Ipsum",
		LearnerLimit:       40,
		EnrollmentDeadline: time.Now().Add(time.Hour * 72),
		ClassStart:         time.Now().Add(time.Hour * 108),
		ClassFinish:        time.Now().Add(time.Hour * 110),
		ClassStatus:        "Available",
	}

	resp := runHTTP(t, app, httpInput{
		Method:      http.MethodPost,
		Path:        "/class_sessions/",
		Body:        jsonBody(payload),
		ContentType: "application/json",
		UserID:      &userID,
	})
	wantStatus(t, resp, http.StatusCreated)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

// 400
func TestCreateClassSession_BadRequest(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)
	userID := uint(42)

	ExpAuthUser(userID, false, true, false)(mock)

	app := setupApp(gdb)
	resp := runHTTP(t, app, httpInput{
		Method: http.MethodPost, Path: "/class_sessions/",
		Body: []byte(`{invalid-json}`), ContentType: "application/json", UserID: &userID,
	})
	wantStatus(t, resp, http.StatusBadRequest)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

// 500
func TestCreateClassSession_DBError(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)

	table := "class_sessions"
	userID := uint(42)
	classID := uint(50)

	ExpAuthUser(userID, false, true, false)(mock)
	ExpInsertError(table, fmt.Errorf("db insert failed"))(mock)

	payload := models.ClassSession{
		ClassID:            classID,
		Description:        "Lorem Ipsum",
		LearnerLimit:       40,
		EnrollmentDeadline: time.Now().Add(time.Hour * 72),
		ClassStart:         time.Now().Add(time.Hour * 108),
		ClassFinish:        time.Now().Add(time.Hour * 110),
		ClassStatus:        "Available",
	}
	app := setupApp(gdb)
	resp := runHTTP(t, app, httpInput{
		Method: http.MethodPost, Path: "/class_sessions/",
		Body: jsonBody(payload), ContentType: "application/json", UserID: &userID,
	})
	wantStatus(t, resp, http.StatusInternalServerError)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

/* ------------------ GetClassSessions ------------------ */
// 200
func TestGetClassSessions_OK(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)
	userID := uint(42)

	ExpAuthUser(userID, false, false, false)(mock)
	ExpListRows("class_sessions", []string{"id"}, []any{1}, []any{2})(mock)

	app := setupApp(gdb)
	resp := runHTTP(t, app, httpInput{
		Method: http.MethodGet, Path: "/class_sessions/", UserID: &userID,
	})
	wantStatus(t, resp, http.StatusOK)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

// 500
func TestGetClassSessions_DBError(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)
	userID := uint(42)

	ExpAuthUser(userID, false, false, false)(mock)
	ExpListError("class_sessions", fmt.Errorf("select failed"))(mock)

	app := setupApp(gdb)
	resp := runHTTP(t, app, httpInput{
		Method: http.MethodGet, Path: "/class_sessions/", UserID: &userID,
	})
	wantStatus(t, resp, http.StatusInternalServerError)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

/* ------------------ GetClassSession ------------------ */
// 200
func TestGetClassSession_OK(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)

	table := "class_sessions"
	userID := uint(42)
	classSessionID := uint(7)

	ExpAuthUser(userID, false, false, false)(mock)
	ExpSelectByIDFound(table, classSessionID, []string{"id"}, []any{classSessionID})(mock)

	app := setupApp(gdb)
	resp := runHTTP(t, app, httpInput{
		Method: http.MethodGet, Path: fmt.Sprintf("/class_sessions/%d", classSessionID), UserID: &userID,
	})
	wantStatus(t, resp, http.StatusOK)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

// 404
func TestGetClassSession_NotFound(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)

	table := "class_sessions"
	userID := uint(42)
	classSessionID := uint(999)

	ExpAuthUser(userID, false, false, false)(mock)
	ExpSelectByIDEmpty(table, classSessionID)(mock)

	app := setupApp(gdb)
	resp := runHTTP(t, app, httpInput{
		Method: http.MethodGet, Path: fmt.Sprintf("/class_sessions/%d", classSessionID), UserID: &userID,
	})
	wantStatus(t, resp, http.StatusNotFound)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

// 500
func TestGetClassSession_DBError(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)

	table := "class_sessions"
	userID := uint(42)
	classSessionID := uint(7)

	ExpAuthUser(userID, false, false, false)(mock)
	ExpSelectByIDError(table, classSessionID, fmt.Errorf("select failed"))(mock)

	app := setupApp(gdb)
	resp := runHTTP(t, app, httpInput{
		Method: http.MethodGet, Path: fmt.Sprintf("/class_sessions/%d", classSessionID), UserID: &userID,
	})
	wantStatus(t, resp, http.StatusInternalServerError)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

// 400
func TestGetClassSession_BadRequest(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)
	userID := uint(42)

	ExpAuthUser(userID, false, false, false)(mock)

	app := setupApp(gdb)
	resp := runHTTP(t, app, httpInput{
		Method: http.MethodGet, Path: "/class_sessions/not-an-int", UserID: &userID,
	})
	wantStatus(t, resp, http.StatusBadRequest)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

/* ------------------ UpdateClassSession ------------------ */

// 200
func TestUpdateClassSession_OK(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)

	table := "class_sessions"

	userID := uint(42)
	classID := uint(50)
	classSessionID := uint(1)

	ExpAuthUser(userID, false, true, false)(mock)
	ExpSelectByIDFound(table, classSessionID,
		[]string{"id", "class_id", "description", "learner_limit", "enrollment_deadline", "class_start", "class_finish", "class_status"},
		[]any{classSessionID, classID, "Lorem", 40, time.Now().Add(time.Hour * 72), time.Now().Add(time.Hour * 108), time.Now().Add(time.Hour * 110), "pending"},
	)(mock)

	ExpPreloadField("classes", []string{"id"}, []any{classID})(mock)

	ExpUpdateOK(table)(mock)

	app := setupApp(gdb)
	payload := models.ClassSession{
		ClassID:            classID,
		Description:        "Lorem Ipsum",
		LearnerLimit:       40,
		EnrollmentDeadline: time.Now().Add(time.Hour * 72),
		ClassStart:         time.Now().Add(time.Hour * 108),
		ClassFinish:        time.Now().Add(time.Hour * 110),
		ClassStatus:        "Available",
	}
	resp := runHTTP(t, app, httpInput{
		Method: http.MethodPut, Path: fmt.Sprintf("/class_sessions/%d", classSessionID),
		Body: jsonBody(payload), ContentType: "application/json", UserID: &userID,
	})
	wantStatus(t, resp, http.StatusOK)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

// 404
func TestUpdateClassSession_NotFound(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)

	table := "class_sessions"
	userID := uint(42)
	classSessionID := uint(12345)

	ExpAuthUser(userID, false, true, false)(mock)
	ExpSelectByIDEmpty(table, classSessionID)(mock)

	app := setupApp(gdb)
	resp := runHTTP(t, app, httpInput{
		Method: http.MethodPut, Path: fmt.Sprintf("/class_sessions/%d", classSessionID), UserID: &userID,
	})
	wantStatus(t, resp, http.StatusNotFound)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

// 500
func TestUpdateClassSession_DBError(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)

	table := "class_sessions"

	userID := uint(42)
	classID := uint(50)
	classSessionID := uint(1)

	ExpAuthUser(userID, false, true, false)(mock)
	ExpSelectByIDFound(table, classSessionID, []string{"id"}, []any{classSessionID})(mock)
	ExpUpdateError(table, fmt.Errorf("update failed"))(mock)

	app := setupApp(gdb)
	payload := models.ClassSession{
		ClassID:            classID,
		Description:        "Lorem Ipsum",
		LearnerLimit:       40,
		EnrollmentDeadline: time.Now().Add(time.Hour * 72),
		ClassStart:         time.Now().Add(time.Hour * 108),
		ClassFinish:        time.Now().Add(time.Hour * 110),
		ClassStatus:        "Available",
	}
	resp := runHTTP(t, app, httpInput{
		Method: http.MethodPut, Path: fmt.Sprintf("/class_sessions/%d", classSessionID),
		Body: jsonBody(payload), ContentType: "application/json", UserID: &userID,
	})
	wantStatus(t, resp, http.StatusInternalServerError)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

// 400
func TestUpdateClassSession_BadRequest(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)
	userID := uint(42)

	ExpAuthUser(userID, false, true, false)(mock)

	app := setupApp(gdb)
	resp := runHTTP(t, app, httpInput{
		Method: http.MethodPut, Path: "/class_sessions/not-an-int", UserID: &userID,
	})
	wantStatus(t, resp, http.StatusBadRequest)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

/* ------------------ DeleteClassSession ------------------ */

// 200
func TestDeleteClassSession_OK_SoftDelete(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)

	table := "class_sessions"
	userID := uint(42)
	classSessionID := uint(5)

	ExpAuthUser(userID, false, true, false)(mock)
	ExpSelectByIDFound(table, classSessionID, []string{"id"}, []any{classSessionID})(mock)
	ExpSoftDeleteOK(table)(mock)

	app := setupApp(gdb)
	resp := runHTTP(t, app, httpInput{
		Method: http.MethodDelete, Path: fmt.Sprintf("/class_sessions/%d", classSessionID), UserID: &userID,
	})
	wantStatus(t, resp, http.StatusOK)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

// 404
func TestDeleteClassSession_NotFound(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)

	table := "class_sessions"
	userID := uint(42)
	classSessionID := uint(12345)

	ExpAuthUser(userID, false, true, false)(mock)
	ExpSelectByIDEmpty(table, classSessionID)(mock)

	app := setupApp(gdb)
	resp := runHTTP(t, app, httpInput{
		Method: http.MethodDelete, Path: fmt.Sprintf("/class_sessions/%d", classSessionID), UserID: &userID,
	})
	wantStatus(t, resp, http.StatusNotFound)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

// 500
func TestDeleteClassSession_DBError(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)

	table := "class_sessions"
	userID := uint(42)
	classSessionID := uint(5)

	ExpAuthUser(userID, false, true, false)(mock)
	ExpSelectByIDFound(table, classSessionID, []string{"id"}, []any{classSessionID})(mock)
	ExpSoftDeleteError(table, fmt.Errorf("update failed"))(mock)

	app := setupApp(gdb)
	resp := runHTTP(t, app, httpInput{
		Method: http.MethodDelete, Path: fmt.Sprintf("/class_sessions/%d", classSessionID), UserID: &userID,
	})
	wantStatus(t, resp, http.StatusInternalServerError)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

// 400
func TestDeleteClassSession_BadRequest(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)
	userID := uint(42)

	ExpAuthUser(userID, false, true, false)(mock)

	app := setupApp(gdb)
	resp := runHTTP(t, app, httpInput{
		Method: http.MethodDelete, Path: "/class_sessions/not-an-int", UserID: &userID,
	})
	wantStatus(t, resp, http.StatusBadRequest)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
