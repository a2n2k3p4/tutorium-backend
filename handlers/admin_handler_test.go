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

/* ------------------ CreateAdmin ------------------ */

// code 201
func TestCreateAdmin_OK(t *testing.T) {
	cases := []struct {
		name   string
		status bool
	}{
		{"bypass", true},
		{"unbypass", false},
	}

	for _, c := range cases {
		if c.status {
			t.Setenv("STATUS", "development")
		} else {
			t.Setenv("STATUS", "production")
		}

		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)
		userID := 42

		// Auth: user + preloads
		preloadUserForAuth(mock, uint(userID), false, false, false)

		// Handler: insert admins
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "admins".*RETURNING "id"`).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		app := setupApp(gdb)
		token := makeJWT(t, []byte(secretString), uint(userID))

		req := httptest.NewRequest(http.MethodPost, "/admins/", bytes.NewReader([]byte(`{}`)))
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
}

// code 400
func TestCreateAdmin_BadRequest(t *testing.T) {
	cases := []struct {
		name   string
		status bool
	}{
		{"bypass", true},
		{"unbypass", false},
	}

	for _, c := range cases {
		if c.status {
			t.Setenv("STATUS", "development")
		} else {
			t.Setenv("STATUS", "production")
		}
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)
		userID := 42

		preloadUserForAuth(mock, uint(userID), false, false, false)

		app := setupApp(gdb)
		token := makeJWT(t, []byte(secretString), uint(userID))

		req := httptest.NewRequest(http.MethodPost, "/admins/", bytes.NewBufferString(`{invalid-json}`))
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
}

// code 500
func TestCreateAdmin_DBError(t *testing.T) {
	cases := []struct {
		name   string
		status bool
	}{
		{"bypass", true},
		{"unbypass", false},
	}

	for _, c := range cases {
		if c.status {
			t.Setenv("STATUS", "development")
		} else {
			t.Setenv("STATUS", "production")
		}
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)
		userID := 42

		preloadUserForAuth(mock, uint(userID), false, false, false)

		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "admins".*RETURNING "id"`).
			WillReturnError(fmt.Errorf("db insert failed"))
		mock.ExpectRollback()

		app := setupApp(gdb)
		token := makeJWT(t, []byte(secretString), uint(userID))

		req := httptest.NewRequest(http.MethodPost, "/admins/", bytes.NewReader([]byte(`{}`)))
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
}

/* ------------------ GetAdmins ------------------ */
//code 200
func TestGetAdmins_OK(t *testing.T) {
	cases := []struct {
		name   string
		status bool
	}{
		{"bypass", true},
		{"unbypass", false},
	}

	for _, c := range cases {
		if c.status {
			t.Setenv("STATUS", "development")
		} else {
			t.Setenv("STATUS", "production")
		}
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)
		userID := 42

		// Auth: user + preloads
		preloadUserForAuth(mock, uint(userID), false, false, false)

		// Handler: list admins
		mock.ExpectQuery(`SELECT .* FROM "admins".*`).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1).AddRow(2))

		app := setupApp(gdb)
		token := makeJWT(t, []byte(secretString), uint(userID))

		req := httptest.NewRequest(http.MethodGet, "/admins/", nil)
		req.Header.Set("Authorization", "Bearer "+token)

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
			t.Fatalf("expected 2 admins, got %d (%v)", len(arr), arr)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	}
}

// code 500
func TestGetAdmins_DBError(t *testing.T) {
	cases := []struct {
		name   string
		status bool
	}{
		{"bypass", true},
		{"unbypass", false},
	}

	for _, c := range cases {
		if c.status {
			t.Setenv("STATUS", "development")
		} else {
			t.Setenv("STATUS", "production")
		}
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)
		userID := 42

		preloadUserForAuth(mock, uint(userID), false, false, false)

		mock.ExpectQuery(`SELECT .* FROM "admins".*`).
			WillReturnError(fmt.Errorf("select failed"))

		app := setupApp(gdb)
		token := makeJWT(t, []byte(secretString), uint(userID))

		req := httptest.NewRequest(http.MethodGet, "/admins/", nil)
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
}

/* ------------------ GetAdmin ------------------ */
//code 200
func TestGetAdmin_OK(t *testing.T) {
	cases := []struct {
		name   string
		status bool
	}{
		{"bypass", true},
		{"unbypass", false},
	}

	for _, c := range cases {
		if c.status {
			t.Setenv("STATUS", "development")
		} else {
			t.Setenv("STATUS", "production")
		}
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)
		userID := 42

		adminID := 7

		// Auth: user + preloads
		preloadUserForAuth(mock, uint(userID), false, false, false)

		// Handler: find by id
		mock.ExpectQuery(`SELECT \* FROM "admins" WHERE id = \$1 AND "admins"\."deleted_at" IS NULL ORDER BY "admins"\."id" LIMIT .*`).
			WithArgs(adminID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(adminID))

		app := setupApp(gdb)
		token := makeJWT(t, []byte(secretString), uint(userID))

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/admins/%d", adminID), nil)
		req.Header.Set("Authorization", "Bearer "+token)

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
}

// code 404
func TestGetAdmin_NotFound(t *testing.T) {
	cases := []struct {
		name   string
		status bool
	}{
		{"bypass", true},
		{"unbypass", false},
	}

	for _, c := range cases {
		if c.status {
			t.Setenv("STATUS", "development")
		} else {
			t.Setenv("STATUS", "production")
		}
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)
		userID := 42

		adminID := 999

		// Auth: user + preloads
		preloadUserForAuth(mock, uint(userID), false, false, false)

		// Handler: not found (empty rowset)
		mock.ExpectQuery(`SELECT \* FROM "admins" WHERE id = \$1 AND "admins"\."deleted_at" IS NULL ORDER BY "admins"\."id" LIMIT .*`).
			WithArgs(adminID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}))

		app := setupApp(gdb)
		token := makeJWT(t, []byte(secretString), uint(userID))

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/admins/%d", adminID), nil)
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
}

// code 500
func TestGetAdmin_DBError(t *testing.T) {
	cases := []struct {
		name   string
		status bool
	}{
		{"bypass", true},
		{"unbypass", false},
	}

	for _, c := range cases {
		if c.status {
			t.Setenv("STATUS", "development")
		} else {
			t.Setenv("STATUS", "production")
		}
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)
		userID := 42

		adminID := 7

		preloadUserForAuth(mock, uint(userID), false, false, false)

		mock.ExpectQuery(`SELECT \* FROM "admins" WHERE id = \$1 AND "admins"\."deleted_at" IS NULL ORDER BY "admins"\."id" LIMIT .*`).
			WithArgs(adminID, 1).
			WillReturnError(fmt.Errorf("select failed"))

		app := setupApp(gdb)
		token := makeJWT(t, []byte(secretString), uint(userID))

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/admins/%d", adminID), nil)
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
}

// code 400
func TestGetAdmin_BadRequest(t *testing.T) {
	cases := []struct {
		name   string
		status bool
	}{
		{"bypass", true},
		{"unbypass", false},
	}

	for _, c := range cases {
		if c.status {
			t.Setenv("STATUS", "development")
		} else {
			t.Setenv("STATUS", "production")
		}
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)
		userID := 42

		preloadUserForAuth(mock, uint(userID), false, false, false)

		app := setupApp(gdb)
		token := makeJWT(t, []byte(secretString), uint(userID))

		req := httptest.NewRequest(http.MethodGet, "/admins/not-an-int", nil)
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
}

/* ------------------ DeleteAdmins ------------------ */
//code 200
func TestDeleteAdmin_OK_SoftDelete(t *testing.T) {
	cases := []struct {
		name   string
		status bool
	}{
		{"bypass", true},
		{"unbypass", false},
	}

	for _, c := range cases {
		if c.status {
			t.Setenv("STATUS", "development")
		} else {
			t.Setenv("STATUS", "production")
		}
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)
		userID := 42

		adminID := 5

		// Auth: user + preloads
		preloadUserForAuth(mock, uint(userID), false, false, false)

		// Handler: find then soft-delete
		mock.ExpectQuery(`SELECT \* FROM "admins" WHERE id = \$1 AND "admins"\."deleted_at" IS NULL ORDER BY "admins"\."id" LIMIT .*`).
			WithArgs(adminID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(adminID))

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "admins" SET "deleted_at"=`).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		app := setupApp(gdb)
		token := makeJWT(t, []byte(secretString), uint(userID))

		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/admins/%d", adminID), nil)
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
}

// code 404
func TestDeleteAdmin_NotFound(t *testing.T) {
	cases := []struct {
		name   string
		status bool
	}{
		{"bypass", true},
		{"unbypass", false},
	}

	for _, c := range cases {
		if c.status {
			t.Setenv("STATUS", "development")
		} else {
			t.Setenv("STATUS", "production")
		}
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)
		userID := 42

		adminID := 12345

		// Auth: user + preloads
		preloadUserForAuth(mock, uint(userID), false, false, false)

		// Handler: not found on SELECT
		mock.ExpectQuery(`SELECT \* FROM "admins" WHERE id = \$1 AND "admins"\."deleted_at" IS NULL ORDER BY "admins"\."id" LIMIT .*`).
			WithArgs(adminID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}))
		app := setupApp(gdb)
		token := makeJWT(t, []byte(secretString), uint(userID))

		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/admins/%d", adminID), nil)
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
}

// code 500
func TestDeleteAdmin_DBError(t *testing.T) {
	cases := []struct {
		name   string
		status bool
	}{
		{"bypass", true},
		{"unbypass", false},
	}

	for _, c := range cases {
		if c.status {
			t.Setenv("STATUS", "development")
		} else {
			t.Setenv("STATUS", "production")
		}
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)
		userID := 42

		adminID := 5

		preloadUserForAuth(mock, uint(userID), false, false, false)

		// Found, then UPDATE fails -> rollback
		mock.ExpectQuery(`SELECT \* FROM "admins" WHERE id = \$1 AND "admins"\."deleted_at" IS NULL ORDER BY "admins"\."id" LIMIT .*`).
			WithArgs(adminID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(5))

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "admins" SET "deleted_at"=`).
			WillReturnError(fmt.Errorf("update failed"))
		mock.ExpectRollback()

		app := setupApp(gdb)
		token := makeJWT(t, []byte(secretString), uint(userID))

		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/admins/%d", adminID), nil)
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
}

// code 400
func TestDeleteAdmin_BadRequest(t *testing.T) {
	cases := []struct {
		name   string
		status bool
	}{
		{"bypass", true},
		{"unbypass", false},
	}

	for _, c := range cases {
		if c.status {
			t.Setenv("STATUS", "development")
		} else {
			t.Setenv("STATUS", "production")
		}
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)
		userID := 42

		preloadUserForAuth(mock, uint(userID), false, false, false)

		app := setupApp(gdb)
		token := makeJWT(t, []byte(secretString), uint(userID))

		req := httptest.NewRequest(http.MethodDelete, "/admins/not-an-int", nil)
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
}
