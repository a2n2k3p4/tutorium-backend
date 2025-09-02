package handlers

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/a2n2k3p4/tutorium-backend/models"
)

/* ------------------ CreateAdmin ------------------ */

// 201
func TestCreateAdmin_OK(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)

	table := "admins"
	userID := uint(42)

	ExpAuthUser(userID, false, false, false)(mock)
	ExpInsertReturningID(table, 1)(mock)

	app := setupApp(gdb)

	payload := models.Admin{
		UserID: userID,
	}

	resp := runHTTP(t, app, httpInput{
		Method:      http.MethodPost,
		Path:        "/admins/",
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
func TestCreateAdmin_BadRequest(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)
	userID := uint(42)

	ExpAuthUser(userID, false, false, false)(mock)

	app := setupApp(gdb)
	resp := runHTTP(t, app, httpInput{
		Method: http.MethodPost, Path: "/admins/",
		Body: []byte(`{invalid-json}`), ContentType: "application/json", UserID: &userID,
	})
	wantStatus(t, resp, http.StatusBadRequest)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

// 500
func TestCreateAdmin_DBError(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)

	table := "admins"
	userID := uint(42)

	ExpAuthUser(userID, false, false, false)(mock)
	ExpInsertError(table, fmt.Errorf("db insert failed"))(mock)

	payload := models.Admin{
		UserID: userID,
	}
	app := setupApp(gdb)
	resp := runHTTP(t, app, httpInput{
		Method: http.MethodPost, Path: "/admins/",
		Body: jsonBody(payload), ContentType: "application/json", UserID: &userID,
	})
	wantStatus(t, resp, http.StatusInternalServerError)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

/* ------------------ GetAdmins ------------------ */
// 200
func TestGetAdmins_OK(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)
	userID := uint(42)

	ExpAuthUser(userID, false, false, false)(mock)
	ExpListRows("admins", []string{"id"}, []any{1}, []any{2})(mock)

	app := setupApp(gdb)
	resp := runHTTP(t, app, httpInput{
		Method: http.MethodGet, Path: "/admins/", UserID: &userID,
	})
	wantStatus(t, resp, http.StatusOK)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

// 500
func TestGetAdmins_DBError(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)
	userID := uint(42)

	ExpAuthUser(userID, false, false, false)(mock)
	ExpListError("admins", fmt.Errorf("select failed"))(mock)

	app := setupApp(gdb)
	resp := runHTTP(t, app, httpInput{
		Method: http.MethodGet, Path: "/admins/", UserID: &userID,
	})
	wantStatus(t, resp, http.StatusInternalServerError)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

/* ------------------ GetAdmin ------------------ */
// 200
func TestGetAdmin_OK(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)

	table := "admins"
	userID := uint(42)
	adminID := uint(7)

	ExpAuthUser(userID, false, false, false)(mock)
	ExpSelectByIDFound(table, adminID, []string{"id"}, []any{adminID})(mock)

	app := setupApp(gdb)
	resp := runHTTP(t, app, httpInput{
		Method: http.MethodGet, Path: fmt.Sprintf("/admins/%d", adminID), UserID: &userID,
	})
	wantStatus(t, resp, http.StatusOK)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

// 404
func TestGetAdmin_NotFound(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)

	table := "admins"
	userID := uint(42)
	adminID := uint(999)

	ExpAuthUser(userID, false, false, false)(mock)
	ExpSelectByIDEmpty(table, adminID)(mock)

	app := setupApp(gdb)
	resp := runHTTP(t, app, httpInput{
		Method: http.MethodGet, Path: fmt.Sprintf("/admins/%d", adminID), UserID: &userID,
	})
	wantStatus(t, resp, http.StatusNotFound)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

// 500
func TestGetAdmin_DBError(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)

	table := "admins"
	userID := uint(42)
	adminID := uint(7)

	ExpAuthUser(userID, false, false, false)(mock)
	ExpSelectByIDError(table, adminID, fmt.Errorf("select failed"))(mock)

	app := setupApp(gdb)
	resp := runHTTP(t, app, httpInput{
		Method: http.MethodGet, Path: fmt.Sprintf("/admins/%d", adminID), UserID: &userID,
	})
	wantStatus(t, resp, http.StatusInternalServerError)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

// 400
func TestGetAdmin_BadRequest(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)
	userID := uint(42)

	ExpAuthUser(userID, false, false, false)(mock)

	app := setupApp(gdb)
	resp := runHTTP(t, app, httpInput{
		Method: http.MethodGet, Path: "/admins/not-an-int", UserID: &userID,
	})
	wantStatus(t, resp, http.StatusBadRequest)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

/* ------------------ DeleteAdmin ------------------ */

// 200
func TestDeleteAdmin_OK_SoftDelete(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)

	table := "admins"
	userID := uint(42)
	adminID := uint(5)

	ExpAuthUser(userID, false, false, false)(mock)
	ExpSelectByIDFound(table, adminID, []string{"id"}, []any{adminID})(mock)
	ExpSoftDeleteOK(table)(mock)

	app := setupApp(gdb)
	resp := runHTTP(t, app, httpInput{
		Method: http.MethodDelete, Path: fmt.Sprintf("/admins/%d", adminID), UserID: &userID,
	})
	wantStatus(t, resp, http.StatusOK)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

// 404
func TestDeleteAdmin_NotFound(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)

	table := "admins"
	userID := uint(42)
	adminID := uint(12345)

	ExpAuthUser(userID, false, false, false)(mock)
	ExpSelectByIDEmpty(table, adminID)(mock)

	app := setupApp(gdb)
	resp := runHTTP(t, app, httpInput{
		Method: http.MethodDelete, Path: fmt.Sprintf("/admins/%d", adminID), UserID: &userID,
	})
	wantStatus(t, resp, http.StatusNotFound)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

// 500
func TestDeleteAdmin_DBError(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)

	table := "admins"
	userID := uint(42)
	adminID := uint(5)

	ExpAuthUser(userID, false, false, false)(mock)
	ExpSelectByIDFound(table, adminID, []string{"id"}, []any{adminID})(mock)
	ExpSoftDeleteError(table, fmt.Errorf("update failed"))(mock)

	app := setupApp(gdb)
	resp := runHTTP(t, app, httpInput{
		Method: http.MethodDelete, Path: fmt.Sprintf("/admins/%d", adminID), UserID: &userID,
	})
	wantStatus(t, resp, http.StatusInternalServerError)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

// 400
func TestDeleteAdmin_BadRequest(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)
	userID := uint(42)

	ExpAuthUser(userID, false, false, false)(mock)

	app := setupApp(gdb)
	resp := runHTTP(t, app, httpInput{
		Method: http.MethodDelete, Path: "/admins/not-an-int", UserID: &userID,
	})
	wantStatus(t, resp, http.StatusBadRequest)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
