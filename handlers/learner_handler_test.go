package handlers

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/a2n2k3p4/tutorium-backend/models"
)

/* ------------------ CreateLearner ------------------ */

// 201
func TestCreateLearner_OK(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)

	table := "learners"
	userID := uint(42)

	ExpAuthUser(userID, false, false, false)(mock)
	ExpInsertReturningID(table, 1)(mock)

	app := setupApp(gdb)

	payload := models.Learner{
		UserID: userID,
	}

	resp := runHTTP(t, app, httpInput{
		Method:      http.MethodPost,
		Path:        "/learners/",
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
func TestCreateLearner_BadRequest(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)
	userID := uint(42)

	ExpAuthUser(userID, false, false, false)(mock)

	app := setupApp(gdb)
	resp := runHTTP(t, app, httpInput{
		Method: http.MethodPost, Path: "/learners/",
		Body: []byte(`{invalid-json}`), ContentType: "application/json", UserID: &userID,
	})
	wantStatus(t, resp, http.StatusBadRequest)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

// 500
func TestCreateLearner_DBError(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)

	table := "learners"
	userID := uint(42)

	ExpAuthUser(userID, false, false, false)(mock)
	ExpInsertError(table, fmt.Errorf("db insert failed"))(mock)

	payload := models.Learner{
		UserID: userID,
	}
	app := setupApp(gdb)
	resp := runHTTP(t, app, httpInput{
		Method: http.MethodPost, Path: "/learners/",
		Body: jsonBody(payload), ContentType: "application/json", UserID: &userID,
	})
	wantStatus(t, resp, http.StatusInternalServerError)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

/* ------------------ GetLearners ------------------ */
// 200
func TestGetLearners_OK(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)
	userID := uint(42)

	ExpAuthUser(userID, false, false, false)(mock)
	ExpListRows("learners", []string{"id"}, []any{1}, []any{2})(mock)

	app := setupApp(gdb)
	resp := runHTTP(t, app, httpInput{
		Method: http.MethodGet, Path: "/learners/", UserID: &userID,
	})
	wantStatus(t, resp, http.StatusOK)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

// 500
func TestGetLearners_DBError(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)
	userID := uint(42)

	ExpAuthUser(userID, false, false, false)(mock)
	ExpListError("learners", fmt.Errorf("select failed"))(mock)

	app := setupApp(gdb)
	resp := runHTTP(t, app, httpInput{
		Method: http.MethodGet, Path: "/learners/", UserID: &userID,
	})
	wantStatus(t, resp, http.StatusInternalServerError)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

/* ------------------ GetLearner ------------------ */
// 200
func TestGetLearner_OK(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)

	table := "learners"
	userID := uint(42)
	learnerID := uint(7)

	ExpAuthUser(userID, false, false, false)(mock)
	ExpSelectByIDFound(table, learnerID, []string{"id"}, []any{learnerID})(mock)

	app := setupApp(gdb)
	resp := runHTTP(t, app, httpInput{
		Method: http.MethodGet, Path: fmt.Sprintf("/learners/%d", learnerID), UserID: &userID,
	})
	wantStatus(t, resp, http.StatusOK)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

// 404
func TestGetLearner_NotFound(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)

	table := "learners"
	userID := uint(42)
	learnerID := uint(999)

	ExpAuthUser(userID, false, false, false)(mock)
	ExpSelectByIDEmpty(table, learnerID)(mock)

	app := setupApp(gdb)
	resp := runHTTP(t, app, httpInput{
		Method: http.MethodGet, Path: fmt.Sprintf("/learners/%d", learnerID), UserID: &userID,
	})
	wantStatus(t, resp, http.StatusNotFound)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

// 500
func TestGetLearner_DBError(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)

	table := "learners"
	userID := uint(42)
	learnerID := uint(7)

	ExpAuthUser(userID, false, false, false)(mock)
	ExpSelectByIDError(table, learnerID, fmt.Errorf("select failed"))(mock)

	app := setupApp(gdb)
	resp := runHTTP(t, app, httpInput{
		Method: http.MethodGet, Path: fmt.Sprintf("/learners/%d", learnerID), UserID: &userID,
	})
	wantStatus(t, resp, http.StatusInternalServerError)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

// 400
func TestGetLearner_BadRequest(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)
	userID := uint(42)

	ExpAuthUser(userID, false, false, false)(mock)

	app := setupApp(gdb)
	resp := runHTTP(t, app, httpInput{
		Method: http.MethodGet, Path: "/learners/not-an-int", UserID: &userID,
	})
	wantStatus(t, resp, http.StatusBadRequest)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

/* ------------------ DeleteLearner ------------------ */

// 200
func TestDeleteLearner_OK_SoftDelete(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)

	table := "learners"
	userID := uint(42)
	learnerID := uint(5)

	ExpAuthUser(userID, false, false, false)(mock)
	ExpSelectByIDFound(table, learnerID, []string{"id"}, []any{learnerID})(mock)
	ExpSoftDeleteOK(table)(mock)

	app := setupApp(gdb)
	resp := runHTTP(t, app, httpInput{
		Method: http.MethodDelete, Path: fmt.Sprintf("/learners/%d", learnerID), UserID: &userID,
	})
	wantStatus(t, resp, http.StatusOK)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

// 404
func TestDeleteLearner_NotFound(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)

	table := "learners"
	userID := uint(42)
	learnerID := uint(12345)

	ExpAuthUser(userID, false, false, false)(mock)
	ExpSelectByIDEmpty(table, learnerID)(mock)

	app := setupApp(gdb)
	resp := runHTTP(t, app, httpInput{
		Method: http.MethodDelete, Path: fmt.Sprintf("/learners/%d", learnerID), UserID: &userID,
	})
	wantStatus(t, resp, http.StatusNotFound)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

// 500
func TestDeleteLearner_DBError(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)

	table := "learners"
	userID := uint(42)
	learnerID := uint(5)

	ExpAuthUser(userID, false, false, false)(mock)
	ExpSelectByIDFound(table, learnerID, []string{"id"}, []any{learnerID})(mock)
	ExpSoftDeleteError(table, fmt.Errorf("update failed"))(mock)

	app := setupApp(gdb)
	resp := runHTTP(t, app, httpInput{
		Method: http.MethodDelete, Path: fmt.Sprintf("/learners/%d", learnerID), UserID: &userID,
	})
	wantStatus(t, resp, http.StatusInternalServerError)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

// 400
func TestDeleteLearner_BadRequest(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)
	userID := uint(42)

	ExpAuthUser(userID, false, false, false)(mock)

	app := setupApp(gdb)
	resp := runHTTP(t, app, httpInput{
		Method: http.MethodDelete, Path: "/learners/not-an-int", UserID: &userID,
	})
	wantStatus(t, resp, http.StatusBadRequest)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
