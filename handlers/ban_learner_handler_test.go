package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/a2n2k3p4/tutorium-backend/models"
)

/* ------------------ CreateBanLearner ------------------ */

// 201
func TestCreateBanLearner_OK(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)

	userID := uint(42)
	learnerID := uint(10)
	now := time.Now()
	payload := models.BanDetailsLearner{
		LearnerID:      learnerID,
		BanStart:       now,
		BanEnd:         now.Add(2 * time.Hour),
		BanDescription: "spamming",
	}
	byteString, _ := json.Marshal(payload)
	// focus on handler result: assume authorized learner
	preloadUserForAuth(mock, userID, true, false, false)

	// permissive write expectations
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "ban_details_learners".*RETURNING "id"`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	app := setupApp(gdb)
	token := makeJWT(t, []byte(secretString), userID)

	req := httptest.NewRequest(http.MethodPost, "/banlearners/", bytes.NewReader((byteString)))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want %d; body=%s", resp.StatusCode, http.StatusCreated, string(readBody(t, resp.Body)))
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

// 400
func TestCreateBanLearner_BadRequest(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)

	userID := uint(42)
	preloadUserForAuth(mock, userID, true, false, false)

	app := setupApp(gdb)
	token := makeJWT(t, []byte(secretString), userID)

	req := httptest.NewRequest(http.MethodPost, "/banlearners/", bytes.NewBufferString(`{invalid-json}`))
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

// 500
func TestCreateBanLearner_DBError(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)

	userID := uint(42)
	preloadUserForAuth(mock, userID, true, false, false)

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "ban_details_learners".*RETURNING "id"`).
		WillReturnError(fmt.Errorf("db insert failed"))
	mock.ExpectRollback()

	app := setupApp(gdb)
	token := makeJWT(t, []byte(secretString), userID)

	req := httptest.NewRequest(http.MethodPost, "/banlearners/", bytes.NewReader([]byte(`{}`)))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d; body=%s", resp.StatusCode, http.StatusInternalServerError, string(readBody(t, resp.Body)))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

/* ------------------ GetBanLearners ------------------ */

// 200
func TestGetBanLearners_OK(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()

	userID := uint(42)
	preloadUserForAuth(mock, userID, true, false, false)

	mock.ExpectQuery(`SELECT .* FROM "ban_details_learners".*`).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "learner_id"}).
				AddRow(1, 10).
				AddRow(2, 11),
		)

	mock.ExpectQuery(`SELECT .* FROM "learners".*WHERE .*"learners"\."id" (=\s*\$1|IN \(.*\)).*`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(10).AddRow(11))
	app := setupApp(gdb)
	token := makeJWT(t, []byte(secretString), userID)

	req := httptest.NewRequest(http.MethodGet, "/banlearners/", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", resp.StatusCode, http.StatusOK, string(readBody(t, resp.Body)))
	}

	var arr []map[string]any
	_ = json.Unmarshal(readBody(t, resp.Body), &arr)
	if len(arr) != 2 {
		t.Fatalf("expected 2 ban_details_learners, got %d (%v)", len(arr), arr)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

// 500
func TestGetBanLearners_DBError(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)

	userID := uint(42)
	preloadUserForAuth(mock, userID, true, false, false)

	mock.ExpectQuery(`SELECT .* FROM "ban_details_learners".*`).
		WillReturnError(fmt.Errorf("select failed"))

	app := setupApp(gdb)
	token := makeJWT(t, []byte(secretString), userID)

	req := httptest.NewRequest(http.MethodGet, "/banlearners/", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d; body=%s", resp.StatusCode, http.StatusInternalServerError, string(readBody(t, resp.Body)))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

/* ------------------ GetBanLearner ------------------ */

// 200
func TestGetBanLearner_OK(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()

	userID := uint(42)
	banID := uint(3)
	learnerID := uint(7)
	preloadUserForAuth(mock, userID, true, false, false)

	mock.ExpectQuery(`SELECT .* FROM "ban_details_learners".*WHERE .*id .* LIMIT .*`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "learner_id"}).AddRow(banID, learnerID))

	mock.ExpectQuery(`SELECT .* FROM "learners".*WHERE .*"learners"\."id" (=\s*\$1|IN \(.*\)).*`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(learnerID))
	app := setupApp(gdb)
	token := makeJWT(t, []byte(secretString), userID)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/banlearners/%d", learnerID), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", resp.StatusCode, http.StatusOK, string(readBody(t, resp.Body)))
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

// 404
func TestGetBanLearner_NotFound(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)

	userID := uint(42)
	learnerID := uint(999)
	preloadUserForAuth(mock, userID, true, false, false)

	mock.ExpectQuery(`SELECT .* FROM "ban_details_learners".*WHERE .*id = .* LIMIT .*`).
		WillReturnRows(sqlmock.NewRows([]string{})) // empty -> 404

	app := setupApp(gdb)
	token := makeJWT(t, []byte(secretString), userID)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/banlearners/%d", learnerID), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("status = %d, want %d; body=%s", resp.StatusCode, http.StatusNotFound, string(readBody(t, resp.Body)))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

// 500
func TestGetBanLearner_DBError(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)

	userID := uint(42)
	learnerID := uint(7)
	preloadUserForAuth(mock, userID, true, false, false)

	mock.ExpectQuery(`SELECT .* FROM "ban_details_learners".*WHERE .*id = .* LIMIT .*`).
		WillReturnError(fmt.Errorf("select failed"))

	app := setupApp(gdb)
	token := makeJWT(t, []byte(secretString), userID)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/banlearners/%d", learnerID), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d; body=%s", resp.StatusCode, http.StatusInternalServerError, string(readBody(t, resp.Body)))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

// 400
func TestGetBanLearner_BadRequest(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)

	userID := uint(42)
	preloadUserForAuth(mock, userID, true, false, false)

	app := setupApp(gdb)
	token := makeJWT(t, []byte(secretString), userID)

	req := httptest.NewRequest(http.MethodGet, "/banlearners/not-an-int", nil)
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

/* ------------------UpdateBanLearners ------------------ */
func TestUpdateBanLearner_OK(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()

	userID := uint(42)
	banLearnerID := uint(1)
	learnerID := uint(5)
	preloadUserForAuth(mock, userID, true, false, false)

	// Mock the initial findBanLearner query
	now := time.Now()
	payload := models.BanDetailsLearner{
		LearnerID:      learnerID,
		BanStart:       now,
		BanEnd:         now.Add(2 * time.Hour),
		BanDescription: "spamming",
	}
	byteString, _ := json.Marshal(payload)
	banRows := sqlmock.NewRows([]string{"id", "learner_id", "ban_start", "ban_end", "ban_description"}).
		AddRow(banLearnerID, learnerID, time.Now(), time.Now().Add(24*time.Hour), "Original description")

	mock.ExpectQuery(`SELECT .* FROM "ban_details_learners".*WHERE .*id = .* AND .*deleted_at.*IS NULL.*`).
		WillReturnRows(banRows)

	learnerRows := sqlmock.NewRows([]string{"id", "user_id", "created_at", "updated_at", "deleted_at"}).
		AddRow(learnerID, 123, time.Now(), time.Now(), nil)

	mock.ExpectQuery(`SELECT .* FROM "learners".*WHERE .*id.`).
		WillReturnRows(learnerRows)
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "learners".*RETURNING "id"`).
		WillReturnRows(learnerRows)
	mock.ExpectCommit()
	// Mock the update transaction
	mock.ExpectExec(`UPDATE "ban_details_learners" SET .*updated_at.*WHERE .*deleted_at.*IS NULL.*`).
		WillReturnResult(sqlmock.NewResult(1, 1))

	app := setupApp(gdb)
	token := makeJWT(t, []byte(secretString), userID)

	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/banlearners/%d", banLearnerID),
		bytes.NewReader([]byte(byteString)))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", resp.StatusCode, http.StatusOK, string(readBody(t, resp.Body)))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

// 404
func TestUodateBanLearner_NotFound(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)

	userID := uint(42)
	learnerID := uint(12345)
	preloadUserForAuth(mock, userID, true, false, false)

	mock.ExpectQuery(`SELECT .* FROM "ban_details_learners".*WHERE .*id = .* LIMIT .*`).
		WillReturnRows(sqlmock.NewRows([]string{}))

	app := setupApp(gdb)
	token := makeJWT(t, []byte(secretString), userID)

	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/banlearners/%d", learnerID), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("status = %d, want %d; body=%s", resp.StatusCode, http.StatusNotFound, string(readBody(t, resp.Body)))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

// 500
func TestUpdateBanLearner_DBError(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)

	userID := uint(42)
	learnerID := uint(5)
	banLearnerID := uint(1)
	preloadUserForAuth(mock, userID, true, false, false)

	mock.ExpectQuery(`SELECT .* FROM "ban_details_learners".*WHERE .*id = .* LIMIT .*`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(learnerID))
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "ban_details_learners" SET .*"deleted_at".*`).
		WillReturnError(fmt.Errorf("update failed"))
	mock.ExpectRollback()

	app := setupApp(gdb)
	token := makeJWT(t, []byte(secretString), userID)

	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/banlearners/%d", banLearnerID),
		bytes.NewReader([]byte(`{"ban_description":"Updated description"}`)))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d; body=%s", resp.StatusCode, http.StatusInternalServerError, string(readBody(t, resp.Body)))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

// 400
func TestUpdateBanLearner_BadRequest(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)

	userID := uint(42)
	preloadUserForAuth(mock, userID, true, false, false)

	app := setupApp(gdb)
	token := makeJWT(t, []byte(secretString), userID)

	req := httptest.NewRequest(http.MethodDelete, "/banlearners/not-an-int", nil)
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

/* ------------------ DeleteBanLearners ------------------ */

// 200 (soft delete)
func TestDeleteBanLearner_OK_SoftDelete(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)

	userID := uint(42)
	learnerID := uint(5)
	banID := uint(1)
	preloadUserForAuth(mock, userID, true, false, false)

	// find then soft-delete; keep patterns permissive
	mock.ExpectQuery(`SELECT .* FROM "ban_details_learners".*WHERE .*id = .* `).
		WillReturnRows(sqlmock.NewRows([]string{"id", "learner_id"}).AddRow(banID, learnerID))
	mock.ExpectQuery(`SELECT .* FROM "learners".*WHERE .*"learners"\."id" (=\s*\$1|IN \(.*\)).*`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(learnerID))

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "ban_details_learners" SET .*"deleted_at".*`).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	app := setupApp(gdb)
	token := makeJWT(t, []byte(secretString), userID)

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/banlearners/%d", learnerID), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", resp.StatusCode, http.StatusOK, string(readBody(t, resp.Body)))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

// 404
func TestDeleteBanLearner_NotFound(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)

	userID := uint(42)
	learnerID := uint(12345)
	preloadUserForAuth(mock, userID, true, false, false)

	mock.ExpectQuery(`SELECT .* FROM "ban_details_learners".*WHERE .*id = .* LIMIT .*`).
		WillReturnRows(sqlmock.NewRows([]string{}))

	app := setupApp(gdb)
	token := makeJWT(t, []byte(secretString), userID)

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/banlearners/%d", learnerID), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("status = %d, want %d; body=%s", resp.StatusCode, http.StatusNotFound, string(readBody(t, resp.Body)))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

// 500
func TestDeleteBanLearner_DBError(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)

	userID := uint(42)
	learnerID := uint(5)
	preloadUserForAuth(mock, userID, true, false, false)

	mock.ExpectQuery(`SELECT .* FROM "ban_details_learners".*WHERE .*id = .* LIMIT .*`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(learnerID))
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "ban_details_learners" SET .*"deleted_at".*`).
		WillReturnError(fmt.Errorf("update failed"))
	mock.ExpectRollback()

	app := setupApp(gdb)
	token := makeJWT(t, []byte(secretString), userID)

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/banlearners/%d", learnerID), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d; body=%s", resp.StatusCode, http.StatusInternalServerError, string(readBody(t, resp.Body)))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

// 400
func TestDeleteBanLearner_BadRequest(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)

	userID := uint(42)
	preloadUserForAuth(mock, userID, true, false, false)

	app := setupApp(gdb)
	token := makeJWT(t, []byte(secretString), userID)

	req := httptest.NewRequest(http.MethodDelete, "/banlearners/not-an-int", nil)
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
