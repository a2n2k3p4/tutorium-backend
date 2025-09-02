package handlers

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/a2n2k3p4/tutorium-backend/models"
)

/* ------------------ CreateBanLearner ------------------ */

// 201
func TestCreateBanLearner_OK(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)

	table := "ban_details_learners"
	userID := uint(42)
	learnerID := uint(50)
	now := time.Now()

	ExpAuthUser(userID, true, false, false)(mock)
	ExpInsertReturningID(table, 1)(mock)

	app := setupApp(gdb)

	payload := models.BanDetailsLearner{
		LearnerID:      learnerID,
		BanStart:       now,
		BanEnd:         now.Add(2 * time.Hour),
		BanDescription: "spamming",
	}

	resp := runHTTP(t, app, httpInput{
		Method:      http.MethodPost,
		Path:        "/banlearners/",
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
func TestCreateBanLearner_BadRequest(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)
	userID := uint(42)

	ExpAuthUser(userID, true, false, false)(mock)

	app := setupApp(gdb)
	resp := runHTTP(t, app, httpInput{
		Method: http.MethodPost, Path: "/banlearners/",
		Body: []byte(`{invalid-json}`), ContentType: "application/json", UserID: &userID,
	})
	wantStatus(t, resp, http.StatusBadRequest)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

// 500
func TestCreateBanLearner_DBError(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)

	table := "ban_details_learners"
	userID := uint(42)
	learnerID := uint(50)
	now := time.Now()

	ExpAuthUser(userID, true, false, false)(mock)
	ExpInsertError(table, fmt.Errorf("db insert failed"))(mock)

	payload := models.BanDetailsLearner{
		LearnerID:      learnerID,
		BanStart:       now,
		BanEnd:         now.Add(2 * time.Hour),
		BanDescription: "spamming",
	}
	app := setupApp(gdb)
	resp := runHTTP(t, app, httpInput{
		Method: http.MethodPost, Path: "/banlearners/",
		Body: jsonBody(payload), ContentType: "application/json", UserID: &userID,
	})
	wantStatus(t, resp, http.StatusInternalServerError)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

/* ------------------ GetBanLearners ------------------ */
// 200
func TestGetBanLearners_OK(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)
	userID := uint(42)

	ExpAuthUser(userID, true, false, false)(mock)
	ExpListRows("ban_details_learners", []string{"id"}, []any{1}, []any{2})(mock)

	app := setupApp(gdb)
	resp := runHTTP(t, app, httpInput{
		Method: http.MethodGet, Path: "/banlearners/", UserID: &userID,
	})
	wantStatus(t, resp, http.StatusOK)
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

	ExpAuthUser(userID, true, false, false)(mock)
	ExpListError("ban_details_learners", fmt.Errorf("select failed"))(mock)

	app := setupApp(gdb)
	resp := runHTTP(t, app, httpInput{
		Method: http.MethodGet, Path: "/banlearners/", UserID: &userID,
	})
	wantStatus(t, resp, http.StatusInternalServerError)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

/* ------------------ GetBanLearner ------------------ */
// 200
func TestGetBanLearner_OK(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)

	table := "ban_details_learners"
	userID := uint(42)
	ban_learnerID := uint(7)

	ExpAuthUser(userID, true, false, false)(mock)
	ExpSelectByIDFound(table, ban_learnerID, []string{"id"}, []any{ban_learnerID})(mock)

	app := setupApp(gdb)
	resp := runHTTP(t, app, httpInput{
		Method: http.MethodGet, Path: fmt.Sprintf("/banlearners/%d", ban_learnerID), UserID: &userID,
	})
	wantStatus(t, resp, http.StatusOK)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

// 404
func TestGetBanLearner_NotFound(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)

	table := "ban_details_learners"
	userID := uint(42)
	ban_learnerID := uint(999)

	ExpAuthUser(userID, true, false, false)(mock)
	ExpSelectByIDEmpty(table, ban_learnerID)(mock)

	app := setupApp(gdb)
	resp := runHTTP(t, app, httpInput{
		Method: http.MethodGet, Path: fmt.Sprintf("/banlearners/%d", ban_learnerID), UserID: &userID,
	})
	wantStatus(t, resp, http.StatusNotFound)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

// 500
func TestGetBanLearner_DBError(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)

	table := "ban_details_learners"
	userID := uint(42)
	ban_learnerID := uint(7)

	ExpAuthUser(userID, true, false, false)(mock)
	ExpSelectByIDError(table, ban_learnerID, fmt.Errorf("select failed"))(mock)

	app := setupApp(gdb)
	resp := runHTTP(t, app, httpInput{
		Method: http.MethodGet, Path: fmt.Sprintf("/banlearners/%d", ban_learnerID), UserID: &userID,
	})
	wantStatus(t, resp, http.StatusInternalServerError)
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

	ExpAuthUser(userID, true, false, false)(mock)

	app := setupApp(gdb)
	resp := runHTTP(t, app, httpInput{
		Method: http.MethodGet, Path: "/banlearners/not-an-int", UserID: &userID,
	})
	wantStatus(t, resp, http.StatusBadRequest)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

/* ------------------ UpdateBanLearner ------------------ */

// 200
func TestUpdateBanLearner_OK(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)

	table := "ban_details_learners"
	const preloadTable = "learners"

	userID := uint(42)
	learnerID := uint(50)
	ban_learnerID := uint(1)
	now := time.Now()

	ExpAuthUser(userID, true, false, false)(mock)
	ExpSelectByIDFound(table, ban_learnerID,
		[]string{"id", "learner_id", "ban_start", "ban_end", "ban_description"},
		[]any{ban_learnerID, learnerID, now, now.Add(2 * time.Hour), "flooding"},
	)(mock)

	ExpPreloadField(preloadTable, []string{"id"}, []any{learnerID})(mock)

	ExpUpdateOK(table)(mock)

	app := setupApp(gdb)
	payload := models.BanDetailsLearner{
		BanDescription: "spamming",
	}
	resp := runHTTP(t, app, httpInput{
		Method: http.MethodPut, Path: fmt.Sprintf("/banlearners/%d", ban_learnerID),
		Body: jsonBody(payload), ContentType: "application/json", UserID: &userID,
	})
	wantStatus(t, resp, http.StatusOK)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

// 404
func TestUpdateBanLearner_NotFound(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)

	table := "ban_details_learners"
	userID := uint(42)
	ban_learnerID := uint(12345)

	ExpAuthUser(userID, true, false, false)(mock)
	ExpSelectByIDEmpty(table, ban_learnerID)(mock)

	app := setupApp(gdb)
	resp := runHTTP(t, app, httpInput{
		Method: http.MethodPut, Path: fmt.Sprintf("/banlearners/%d", ban_learnerID), UserID: &userID,
	})
	wantStatus(t, resp, http.StatusNotFound)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

// 500
func TestUpdateBanLearner_DBError(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)

	table := "ban_details_learners"

	userID := uint(42)

	ban_learnerID := uint(1)

	ExpAuthUser(userID, true, false, false)(mock)
	ExpSelectByIDFound(table, ban_learnerID, []string{"id"}, []any{ban_learnerID})(mock)
	ExpUpdateError(table, fmt.Errorf("update failed"))(mock)

	app := setupApp(gdb)
	payload := models.BanDetailsLearner{

		BanDescription: "test err",
	}
	resp := runHTTP(t, app, httpInput{
		Method: http.MethodPut, Path: fmt.Sprintf("/banlearners/%d", ban_learnerID),
		Body: jsonBody(payload), ContentType: "application/json", UserID: &userID,
	})
	wantStatus(t, resp, http.StatusInternalServerError)
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

	ExpAuthUser(userID, true, false, false)(mock)

	app := setupApp(gdb)
	resp := runHTTP(t, app, httpInput{
		Method: http.MethodPut, Path: "/banlearners/not-an-int", UserID: &userID,
	})
	wantStatus(t, resp, http.StatusBadRequest)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

/* ------------------ DeleteBanLearner ------------------ */

// 200
func TestDeleteBanLearner_OK_SoftDelete(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)

	table := "ban_details_learners"
	userID := uint(42)
	ban_learnerID := uint(5)

	ExpAuthUser(userID, true, false, false)(mock)
	ExpSelectByIDFound(table, ban_learnerID, []string{"id"}, []any{ban_learnerID})(mock)
	ExpSoftDeleteOK(table)(mock)

	app := setupApp(gdb)
	resp := runHTTP(t, app, httpInput{
		Method: http.MethodDelete, Path: fmt.Sprintf("/banlearners/%d", ban_learnerID), UserID: &userID,
	})
	wantStatus(t, resp, http.StatusOK)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

// 404
func TestDeleteBanLearner_NotFound(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)

	table := "ban_details_learners"
	userID := uint(42)
	ban_learnerID := uint(12345)

	ExpAuthUser(userID, true, false, false)(mock)
	ExpSelectByIDEmpty(table, ban_learnerID)(mock)

	app := setupApp(gdb)
	resp := runHTTP(t, app, httpInput{
		Method: http.MethodDelete, Path: fmt.Sprintf("/banlearners/%d", ban_learnerID), UserID: &userID,
	})
	wantStatus(t, resp, http.StatusNotFound)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

// 500
func TestDeleteBanLearner_DBError(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()
	mock.MatchExpectationsInOrder(false)

	table := "ban_details_learners"
	userID := uint(42)
	ban_learnerID := uint(5)

	ExpAuthUser(userID, true, false, false)(mock)
	ExpSelectByIDFound(table, ban_learnerID, []string{"id"}, []any{ban_learnerID})(mock)
	ExpSoftDeleteError(table, fmt.Errorf("update failed"))(mock)

	app := setupApp(gdb)
	resp := runHTTP(t, app, httpInput{
		Method: http.MethodDelete, Path: fmt.Sprintf("/banlearners/%d", ban_learnerID), UserID: &userID,
	})
	wantStatus(t, resp, http.StatusInternalServerError)
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

	ExpAuthUser(userID, true, false, false)(mock)

	app := setupApp(gdb)
	resp := runHTTP(t, app, httpInput{
		Method: http.MethodDelete, Path: "/banlearners/not-an-int", UserID: &userID,
	})
	wantStatus(t, resp, http.StatusBadRequest)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
