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

/* ------------------ CreateBanTeacher ------------------ */

// 201
func TestCreateBanTeacher_OK(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)

	userID := uint(42)
	teacherID := uint(10)
	now := time.Now()
	payload := models.BanDetailsTeacher{
		TeacherID:      teacherID,
		BanStart:       now,
		BanEnd:         now.Add(2 * time.Hour),
		BanDescription: "spamming",
	}
	byteString, _ := json.Marshal(payload)
	// focus on handler result: assume authorized teacher
	preloadUserForAuth(mock, userID, true, false, false)

	// permissive write expectations
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "ban_details_teachers".*RETURNING "id"`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	app := setupApp(gdb)
	token := makeJWT(t, []byte(secretString), userID)

	req := httptest.NewRequest(http.MethodPost, "/banteachers/", bytes.NewReader((byteString)))
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
func TestCreateBanTeacher_BadRequest(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)

	userID := uint(42)
	preloadUserForAuth(mock, userID, true, false, false)

	app := setupApp(gdb)
	token := makeJWT(t, []byte(secretString), userID)

	req := httptest.NewRequest(http.MethodPost, "/banteachers/", bytes.NewBufferString(`{invalid-json}`))
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
func TestCreateBanTeacher_DBError(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)

	userID := uint(42)
	preloadUserForAuth(mock, userID, true, false, false)

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "ban_details_teachers".*RETURNING "id"`).
		WillReturnError(fmt.Errorf("db insert failed"))
	mock.ExpectRollback()

	app := setupApp(gdb)
	token := makeJWT(t, []byte(secretString), userID)

	req := httptest.NewRequest(http.MethodPost, "/banteachers/", bytes.NewReader([]byte(`{}`)))
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

/* ------------------ GetBanTeachers ------------------ */

// 200
func TestGetBanTeachers_OK(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()

	userID := uint(42)
	preloadUserForAuth(mock, userID, true, false, false)

	mock.ExpectQuery(`SELECT .* FROM "ban_details_teachers".*`).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "teacher_id"}).
				AddRow(1, 10).
				AddRow(2, 11),
		)

	mock.ExpectQuery(`SELECT .* FROM "teachers".*WHERE .*"teachers"\."id" (=\s*\$1|IN \(.*\)).*`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(10).AddRow(11))
	app := setupApp(gdb)
	token := makeJWT(t, []byte(secretString), userID)

	req := httptest.NewRequest(http.MethodGet, "/banteachers/", nil)
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
		t.Fatalf("expected 2 ban_details_teachers, got %d (%v)", len(arr), arr)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

// 500
func TestGetBanTeachers_DBError(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)

	userID := uint(42)
	preloadUserForAuth(mock, userID, true, false, false)

	mock.ExpectQuery(`SELECT .* FROM "ban_details_teachers".*`).
		WillReturnError(fmt.Errorf("select failed"))

	app := setupApp(gdb)
	token := makeJWT(t, []byte(secretString), userID)

	req := httptest.NewRequest(http.MethodGet, "/banteachers/", nil)
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

/* ------------------ GetBanTeacher ------------------ */

// 200
func TestGetBanTeacher_OK(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()

	userID := uint(42)
	banID := uint(3)
	teacherID := uint(7)
	preloadUserForAuth(mock, userID, true, false, false)

	mock.ExpectQuery(`SELECT .* FROM "ban_details_teachers".*WHERE .*id .* LIMIT .*`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "teacher_id"}).AddRow(banID, teacherID))

	mock.ExpectQuery(`SELECT .* FROM "teachers".*WHERE .*"teachers"\."id" (=\s*\$1|IN \(.*\)).*`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(teacherID))
	app := setupApp(gdb)
	token := makeJWT(t, []byte(secretString), userID)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/banteachers/%d", teacherID), nil)
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
func TestGetBanTeacher_NotFound(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)

	userID := uint(42)
	teacherID := uint(999)
	preloadUserForAuth(mock, userID, true, false, false)

	mock.ExpectQuery(`SELECT .* FROM "ban_details_teachers".*WHERE .*id = .* LIMIT .*`).
		WillReturnRows(sqlmock.NewRows([]string{})) // empty -> 404

	app := setupApp(gdb)
	token := makeJWT(t, []byte(secretString), userID)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/banteachers/%d", teacherID), nil)
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
func TestGetBanTeacher_DBError(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)

	userID := uint(42)
	teacherID := uint(7)
	preloadUserForAuth(mock, userID, true, false, false)

	mock.ExpectQuery(`SELECT .* FROM "ban_details_teachers".*WHERE .*id = .* LIMIT .*`).
		WillReturnError(fmt.Errorf("select failed"))

	app := setupApp(gdb)
	token := makeJWT(t, []byte(secretString), userID)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/banteachers/%d", teacherID), nil)
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
func TestGetBanTeacher_BadRequest(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)

	userID := uint(42)
	preloadUserForAuth(mock, userID, true, false, false)

	app := setupApp(gdb)
	token := makeJWT(t, []byte(secretString), userID)

	req := httptest.NewRequest(http.MethodGet, "/banteachers/not-an-int", nil)
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

/* ------------------UpdateBanTeachers ------------------ */
func TestUpdateBanTeacher_OK(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()

	userID := uint(42)
	banTeacherID := uint(1)
	teacherID := uint(5)
	preloadUserForAuth(mock, userID, true, false, false)

	// Mock the initial findBanTeacher query
	now := time.Now()
	payload := models.BanDetailsTeacher{
		TeacherID:      teacherID,
		BanStart:       now,
		BanEnd:         now.Add(2 * time.Hour),
		BanDescription: "spamming",
	}
	byteString, _ := json.Marshal(payload)
	banRows := sqlmock.NewRows([]string{"id", "teacher_id", "ban_start", "ban_end", "ban_description"}).
		AddRow(banTeacherID, teacherID, time.Now(), time.Now().Add(24*time.Hour), "Original description")

	mock.ExpectQuery(`SELECT .* FROM "ban_details_teachers".*WHERE .*id = .* AND .*deleted_at.*IS NULL.*`).
		WillReturnRows(banRows)

	teacherRows := sqlmock.NewRows([]string{"id", "user_id", "created_at", "updated_at", "deleted_at"}).
		AddRow(teacherID, 123, time.Now(), time.Now(), nil)

	mock.ExpectQuery(`SELECT .* FROM "teachers".*WHERE .*id.`).
		WillReturnRows(teacherRows)
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "teachers".*RETURNING "id"`).
		WillReturnRows(teacherRows)
	mock.ExpectCommit()
	// Mock the update transaction
	mock.ExpectExec(`UPDATE "ban_details_teachers" SET .*updated_at.*WHERE .*deleted_at.*IS NULL.*`).
		WillReturnResult(sqlmock.NewResult(1, 1))

	app := setupApp(gdb)
	token := makeJWT(t, []byte(secretString), userID)

	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/banteachers/%d", banTeacherID),
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
func TestUodateBanTeacher_NotFound(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)

	userID := uint(42)
	teacherID := uint(12345)
	preloadUserForAuth(mock, userID, true, false, false)

	mock.ExpectQuery(`SELECT .* FROM "ban_details_teachers".*WHERE .*id = .* LIMIT .*`).
		WillReturnRows(sqlmock.NewRows([]string{}))

	app := setupApp(gdb)
	token := makeJWT(t, []byte(secretString), userID)

	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/banteachers/%d", teacherID), nil)
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
func TestUpdateBanTeacher_DBError(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)

	userID := uint(42)
	teacherID := uint(5)
	banTeacherID := uint(1)
	preloadUserForAuth(mock, userID, true, false, false)

	mock.ExpectQuery(`SELECT .* FROM "ban_details_teachers".*WHERE .*id = .* LIMIT .*`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(teacherID))
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "ban_details_teachers" SET .*"deleted_at".*`).
		WillReturnError(fmt.Errorf("update failed"))
	mock.ExpectRollback()

	app := setupApp(gdb)
	token := makeJWT(t, []byte(secretString), userID)

	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/banteachers/%d", banTeacherID),
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
func TestUpdateBanTeacher_BadRequest(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)

	userID := uint(42)
	preloadUserForAuth(mock, userID, true, false, false)

	app := setupApp(gdb)
	token := makeJWT(t, []byte(secretString), userID)

	req := httptest.NewRequest(http.MethodDelete, "/banteachers/not-an-int", nil)
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

/* ------------------ DeleteBanTeachers ------------------ */

// 200 (soft delete)
func TestDeleteBanTeacher_OK_SoftDelete(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)

	userID := uint(42)
	teacherID := uint(5)
	banID := uint(1)
	preloadUserForAuth(mock, userID, true, false, false)

	// find then soft-delete; keep patterns permissive
	mock.ExpectQuery(`SELECT .* FROM "ban_details_teachers".*WHERE .*id = .* `).
		WillReturnRows(sqlmock.NewRows([]string{"id", "teacher_id"}).AddRow(banID, teacherID))
	mock.ExpectQuery(`SELECT .* FROM "teachers".*WHERE .*"teachers"\."id" (=\s*\$1|IN \(.*\)).*`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(teacherID))

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "ban_details_teachers" SET .*"deleted_at".*`).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	app := setupApp(gdb)
	token := makeJWT(t, []byte(secretString), userID)

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/banteachers/%d", teacherID), nil)
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
func TestDeleteBanTeacher_NotFound(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)

	userID := uint(42)
	teacherID := uint(12345)
	preloadUserForAuth(mock, userID, true, false, false)

	mock.ExpectQuery(`SELECT .* FROM "ban_details_teachers".*WHERE .*id = .* LIMIT .*`).
		WillReturnRows(sqlmock.NewRows([]string{}))

	app := setupApp(gdb)
	token := makeJWT(t, []byte(secretString), userID)

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/banteachers/%d", teacherID), nil)
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
func TestDeleteBanTeacher_DBError(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)

	userID := uint(42)
	teacherID := uint(5)
	preloadUserForAuth(mock, userID, true, false, false)

	mock.ExpectQuery(`SELECT .* FROM "ban_details_teachers".*WHERE .*id = .* LIMIT .*`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(teacherID))
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "ban_details_teachers" SET .*"deleted_at".*`).
		WillReturnError(fmt.Errorf("update failed"))
	mock.ExpectRollback()

	app := setupApp(gdb)
	token := makeJWT(t, []byte(secretString), userID)

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/banteachers/%d", teacherID), nil)
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
func TestDeleteBanTeacher_BadRequest(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)

	userID := uint(42)
	preloadUserForAuth(mock, userID, true, false, false)

	app := setupApp(gdb)
	token := makeJWT(t, []byte(secretString), userID)

	req := httptest.NewRequest(http.MethodDelete, "/banteachers/not-an-int", nil)
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
