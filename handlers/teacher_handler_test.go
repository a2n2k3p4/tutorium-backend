package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

/* ------------------ CreateTeacher ------------------ */

// code 201
func TestCreateTeacher_OK(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()

	mock.MatchExpectationsInOrder(false)
	userID := 42

	// Auth: user + preloads
	preloadUserForAuth(mock, uint(userID), false, false, false)

	// Handler: insert teachers
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "teachers".*RETURNING "id"`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	app := setupApp(gdb)
	token := makeJWT(t, fileSecret, uint(userID))

	req := httptest.NewRequest(http.MethodPost, "/teachers/", bytes.NewReader([]byte(`{}`)))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusCreated)
	}

	var got map[string]any
	_ = json.Unmarshal(readBody(t, resp.Body), &got)

	if _, ok := got["ID"]; !ok {
		t.Fatalf("response missing ID; got: %v", got)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

// code 400
func TestCreateTeacher_BadRequest(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)
	userID := 42

	preloadUserForAuth(mock, uint(userID), false, false, false)

	app := setupApp(gdb)
	token := makeJWT(t, fileSecret, uint(userID))

	req := httptest.NewRequest(http.MethodPost, "/teachers/", bytes.NewBufferString(`{invalid-json}`))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusBadRequest)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

// code 500
func TestCreateTeacher_DBError(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)
	userID := 42

	preloadUserForAuth(mock, uint(userID), false, false, false)

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "teachers".*RETURNING "id"`).
		WillReturnError(fmt.Errorf("db insert failed"))
	mock.ExpectRollback()

	app := setupApp(gdb)
	token := makeJWT(t, fileSecret, uint(userID))

	req := httptest.NewRequest(http.MethodPost, "/teachers/", bytes.NewReader([]byte(`{}`)))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusInternalServerError)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

/* ------------------ GetTeachers ------------------ */
//code 200
func TestGetTeachers_OK(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()

	// Handler: list teachers
	mock.ExpectQuery(`SELECT .* FROM "teachers".*`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1).AddRow(2))

	app := setupApp(gdb)

	req := httptest.NewRequest(http.MethodGet, "/teachers/", nil)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	var arr []map[string]any
	_ = json.Unmarshal(readBody(t, resp.Body), &arr)
	if len(arr) != 2 {
		t.Fatalf("expected 2 teachers, got %d (%v)", len(arr), arr)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

// code 500
func TestGetTeachers_DBError(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()

	mock.ExpectQuery(`SELECT .* FROM "teachers".*`).
		WillReturnError(fmt.Errorf("select failed"))

	app := setupApp(gdb)

	req := httptest.NewRequest(http.MethodGet, "/teachers/", nil)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusInternalServerError)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

/* ------------------ GetTeacher ------------------ */
//code 200
func TestGetTeacher_OK(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()

	teacherID := 7

	// Handler: find by id
	mock.ExpectQuery(`SELECT \* FROM "teachers" WHERE id = \$1 AND "teachers"\."deleted_at" IS NULL ORDER BY "teachers"\."id" LIMIT .*`).
		WithArgs(teacherID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(teacherID))

	app := setupApp(gdb)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/teachers/%d", teacherID), nil)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	var got map[string]any
	_ = json.Unmarshal(readBody(t, resp.Body), &got)
	if _, ok := got["ID"]; !ok {
		t.Fatalf("response missing ID; got: %v", got)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

// code 404
func TestGetTeacher_NotFound(t *testing.T) {

	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()

	teacherID := 999

	// Handler: not found (empty rowset)
	mock.ExpectQuery(`SELECT \* FROM "teachers" WHERE id = \$1 AND "teachers"\."deleted_at" IS NULL ORDER BY "teachers"\."id" LIMIT .*`).
		WithArgs(teacherID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	app := setupApp(gdb)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/teachers/%d", teacherID), nil)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusNotFound)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

// code 500
func TestGetTeacher_DBError(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	teacherID := 7

	mock.ExpectQuery(`SELECT \* FROM "teachers" WHERE id = \$1 AND "teachers"\."deleted_at" IS NULL ORDER BY "teachers"\."id" LIMIT .*`).
		WithArgs(teacherID, 1).
		WillReturnError(fmt.Errorf("select failed"))

	app := setupApp(gdb)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/teachers/%d", teacherID), nil)

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusInternalServerError)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

// code 400
func TestGetTeacher_BadRequest(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()

	app := setupApp(gdb)

	req := httptest.NewRequest(http.MethodGet, "/teachers/not-an-int", nil)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusBadRequest)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

/* ------------------ DeleteTeachers ------------------ */
//code 200
func TestDeleteTeacher_OK_SoftDelete(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()

	mock.MatchExpectationsInOrder(false)
	userID := 42
	teacherID := 5

	// Auth: user + preloads
	preloadUserForAuth(mock, uint(userID), false, false, false)

	// Handler: find then soft-delete
	mock.ExpectQuery(`SELECT \* FROM "teachers" WHERE id = \$1 AND "teachers"\."deleted_at" IS NULL ORDER BY "teachers"\."id" LIMIT .*`).
		WithArgs(teacherID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(teacherID))

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "teachers" SET "deleted_at"=`).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	app := setupApp(gdb)
	token := makeJWT(t, fileSecret, uint(userID))

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/teachers/%d", teacherID), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

// code 404
func TestDeleteTeacher_NotFound(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()

	mock.MatchExpectationsInOrder(false)
	userID := 42
	teacherID := 12345

	// Auth: user + preloads
	preloadUserForAuth(mock, uint(userID), false, false, false)

	// Handler: not found on SELECT
	mock.ExpectQuery(`SELECT \* FROM "teachers" WHERE id = \$1 AND "teachers"\."deleted_at" IS NULL ORDER BY "teachers"\."id" LIMIT .*`).
		WithArgs(teacherID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))
	app := setupApp(gdb)
	token := makeJWT(t, fileSecret, uint(userID))

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/teachers/%d", teacherID), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusNotFound)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

// code 500
func TestDeleteTeacher_DBError(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)
	userID := 42
	teacherID := 5

	preloadUserForAuth(mock, uint(userID), false, false, false)

	// Found, then UPDATE fails -> rollback
	mock.ExpectQuery(`SELECT \* FROM "teachers" WHERE id = \$1 AND "teachers"\."deleted_at" IS NULL ORDER BY "teachers"\."id" LIMIT .*`).
		WithArgs(teacherID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(5))

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "teachers" SET "deleted_at"=`).
		WillReturnError(fmt.Errorf("update failed"))
	mock.ExpectRollback()

	app := setupApp(gdb)
	token := makeJWT(t, fileSecret, uint(userID))

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/teachers/%d", teacherID), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusInternalServerError)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

// code 400
func TestDeleteTeacher_BadRequest(t *testing.T) {

	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)
	userID := 42

	preloadUserForAuth(mock, uint(userID), false, false, false)

	app := setupApp(gdb)
	token := makeJWT(t, fileSecret, uint(userID))

	req := httptest.NewRequest(http.MethodDelete, "/teachers/not-an-int", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusBadRequest)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
