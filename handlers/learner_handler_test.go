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

/* ------------------ CreateLearner ------------------ */

// code 201
func TestCreateLearner_OK(t *testing.T) {
	mock, cleanup := setupMockGorm(t)
	defer cleanup()

	mock.MatchExpectationsInOrder(false)
	userID := 42

	// Auth: user + preloads
	preloadUserForAuth(mock, uint(userID))

	// Handler: insert learners
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "learners".*RETURNING "id"`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	app := setupApp()
	token := makeJWT(t, uint(userID))

	req := httptest.NewRequest(http.MethodPost, "/learners/", bytes.NewReader([]byte(`{}`)))
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
func TestCreateLearner_BadRequest(t *testing.T) {
	mock, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)
	userID := 42

	preloadUserForAuth(mock, uint(userID))

	app := setupApp()
	token := makeJWT(t, uint(userID))

	req := httptest.NewRequest(http.MethodPost, "/learners/", bytes.NewBufferString(`{invalid-json}`))
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
func TestCreateLearner_DBError(t *testing.T) {
	mock, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)
	userID := 42

	preloadUserForAuth(mock, uint(userID))

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "learners".*RETURNING "id"`).
		WillReturnError(fmt.Errorf("db insert failed"))
	mock.ExpectRollback()

	app := setupApp()
	token := makeJWT(t, uint(userID))

	req := httptest.NewRequest(http.MethodPost, "/learners/", bytes.NewReader([]byte(`{}`)))
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

/* ------------------ GetLearners ------------------ */
//code 200
func TestGetLearners_OK(t *testing.T) {

	mock, cleanup := setupMockGorm(t)
	defer cleanup()

	mock.MatchExpectationsInOrder(false)
	userID := 42

	// Auth: user + preloads
	preloadUserForAuth(mock, uint(userID))

	// Handler: list learners
	mock.ExpectQuery(`SELECT .* FROM "learners".*`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1).AddRow(2))

	app := setupApp()
	token := makeJWT(t, uint(userID))

	req := httptest.NewRequest(http.MethodGet, "/learners/", nil)
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
		t.Fatalf("expected 2 learners, got %d (%v)", len(arr), arr)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

// code 500
func TestGetLearners_DBError(t *testing.T) {

	mock, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)
	userID := 42

	preloadUserForAuth(mock, uint(userID))

	mock.ExpectQuery(`SELECT .* FROM "learners".*`).
		WillReturnError(fmt.Errorf("select failed"))

	app := setupApp()
	token := makeJWT(t, uint(userID))

	req := httptest.NewRequest(http.MethodGet, "/learners/", nil)
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

/* ------------------ GetLearner ------------------ */
//code 200
func TestGetLearner_OK(t *testing.T) {

	mock, cleanup := setupMockGorm(t)
	defer cleanup()

	mock.MatchExpectationsInOrder(false)
	userID := 42
	learnerID := 7

	// Auth: user + preloads
	preloadUserForAuth(mock, uint(userID))

	// Handler: find by id
	mock.ExpectQuery(`SELECT \* FROM "learners" WHERE id = \$1 AND "learners"\."deleted_at" IS NULL ORDER BY "learners"\."id" LIMIT .*`).
		WithArgs(learnerID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(learnerID))

	app := setupApp()
	token := makeJWT(t, uint(userID))

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/learners/%d", learnerID), nil)
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

// code 404
func TestGetLearner_NotFound(t *testing.T) {

	mock, cleanup := setupMockGorm(t)
	defer cleanup()

	mock.MatchExpectationsInOrder(false)
	userID := 42
	learnerID := 999

	// Auth: user + preloads
	preloadUserForAuth(mock, uint(userID))

	// Handler: not found (empty rowset)
	mock.ExpectQuery(`SELECT \* FROM "learners" WHERE id = \$1 AND "learners"\."deleted_at" IS NULL ORDER BY "learners"\."id" LIMIT .*`).
		WithArgs(learnerID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	app := setupApp()
	token := makeJWT(t, uint(userID))

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/learners/%d", learnerID), nil)
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
func TestGetLearner_DBError(t *testing.T) {
	mock, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)
	userID := 42
	learnerID := 7

	preloadUserForAuth(mock, uint(userID))

	mock.ExpectQuery(`SELECT \* FROM "learners" WHERE id = \$1 AND "learners"\."deleted_at" IS NULL ORDER BY "learners"\."id" LIMIT .*`).
		WithArgs(learnerID, 1).
		WillReturnError(fmt.Errorf("select failed"))

	app := setupApp()
	token := makeJWT(t, uint(userID))

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/learners/%d", learnerID), nil)
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
func TestGetLearner_BadRequest(t *testing.T) {
	mock, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)
	userID := 42

	preloadUserForAuth(mock, uint(userID))

	app := setupApp()
	token := makeJWT(t, uint(userID))

	req := httptest.NewRequest(http.MethodGet, "/learners/not-an-int", nil)
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

/* ------------------ DeleteLearners ------------------ */
//code 200
func TestDeleteLearner_OK_SoftDelete(t *testing.T) {
	mock, cleanup := setupMockGorm(t)
	defer cleanup()

	mock.MatchExpectationsInOrder(false)
	userID := 42
	learnerID := 5

	// Auth: user + preloads
	preloadUserForAuth(mock, uint(userID))

	// Handler: find then soft-delete
	mock.ExpectQuery(`SELECT \* FROM "learners" WHERE id = \$1 AND "learners"\."deleted_at" IS NULL ORDER BY "learners"\."id" LIMIT .*`).
		WithArgs(learnerID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(learnerID))

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "learners" SET "deleted_at"=`).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	app := setupApp()
	token := makeJWT(t, uint(userID))

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/learners/%d", learnerID), nil)
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
func TestDeleteLearner_NotFound(t *testing.T) {
	mock, cleanup := setupMockGorm(t)
	defer cleanup()

	mock.MatchExpectationsInOrder(false)
	userID := 42
	learnerID := 12345

	// Auth: user + preloads
	preloadUserForAuth(mock, uint(userID))

	// Handler: not found on SELECT
	mock.ExpectQuery(`SELECT \* FROM "learners" WHERE id = \$1 AND "learners"\."deleted_at" IS NULL ORDER BY "learners"\."id" LIMIT .*`).
		WithArgs(learnerID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))
	app := setupApp()
	token := makeJWT(t, uint(userID))

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/learners/%d", learnerID), nil)
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
func TestDeleteLearner_DBError(t *testing.T) {
	mock, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)
	userID := 42
	learnerID := 5

	preloadUserForAuth(mock, uint(userID))

	// Found, then UPDATE fails -> rollback
	mock.ExpectQuery(`SELECT \* FROM "learners" WHERE id = \$1 AND "learners"\."deleted_at" IS NULL ORDER BY "learners"\."id" LIMIT .*`).
		WithArgs(learnerID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(5))

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "learners" SET "deleted_at"=`).
		WillReturnError(fmt.Errorf("update failed"))
	mock.ExpectRollback()

	app := setupApp()
	token := makeJWT(t, uint(userID))

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/learners/%d", learnerID), nil)
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
func TestDeleteLearner_BadRequest(t *testing.T) {

	mock, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)
	userID := 42

	preloadUserForAuth(mock, uint(userID))

	app := setupApp()
	token := makeJWT(t, uint(userID))

	req := httptest.NewRequest(http.MethodDelete, "/learners/not-an-int", nil)
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
