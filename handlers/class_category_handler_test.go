package handlers

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/a2n2k3p4/tutorium-backend/models"
)

/* ------------------ CreateClassCategory ------------------ */

// 201
func TestCreateClassCategory_OK(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)

		table := "class_categories"
		userID := uint(42)

		ExpAuthUser(userID, true, false, false)(mock)
		ExpInsertReturningID(table, 1)(mock)

		app := setupApp(gdb)

		payload := models.ClassCategory{
			ClassCategory: "test",
		}

		resp := runHTTP(t, app, httpInput{
			Method:      http.MethodPost,
			Path:        "/class_categories/",
			Body:        jsonBody(payload),
			ContentType: "application/json",
			UserID:      &userID,
		})
		wantStatus(t, resp, http.StatusCreated)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

// 400
func TestCreateClassCategory_BadRequest(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)
		userID := uint(42)

		ExpAuthUser(userID, true, false, false)(mock)

		app := setupApp(gdb)
		resp := runHTTP(t, app, httpInput{
			Method: http.MethodPost, Path: "/class_categories/",
			Body: []byte(`{invalid-json}`), ContentType: "application/json", UserID: &userID,
		})
		wantStatus(t, resp, http.StatusBadRequest)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

// 500
func TestCreateClassCategory_DBError(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)

		table := "class_categories"
		userID := uint(42)

		ExpAuthUser(userID, true, false, false)(mock)
		ExpInsertError(table, fmt.Errorf("db insert failed"))(mock)

		app := setupApp(gdb)
		resp := runHTTP(t, app, httpInput{
			Method: http.MethodPost, Path: "/class_categories/",
			Body: []byte(`{}`), ContentType: "application/json", UserID: &userID,
		})
		wantStatus(t, resp, http.StatusInternalServerError)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

/* ------------------ GetClassCategories ------------------ */
// 200
func TestGetClassCategories_OK(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)
		userID := uint(42)

		ExpAuthUser(userID, false, false, false)(mock)
		ExpListRows("class_categories", []string{"id"}, []any{1}, []any{2})(mock)

		app := setupApp(gdb)
		resp := runHTTP(t, app, httpInput{
			Method: http.MethodGet, Path: "/class_categories/", UserID: &userID,
		})
		wantStatus(t, resp, http.StatusOK)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

// 500
func TestGetClassCategories_DBError(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)
		userID := uint(42)

		ExpAuthUser(userID, false, false, false)(mock)
		ExpListError("class_categories", fmt.Errorf("select failed"))(mock)

		app := setupApp(gdb)
		resp := runHTTP(t, app, httpInput{
			Method: http.MethodGet, Path: "/class_categories/", UserID: &userID,
		})
		wantStatus(t, resp, http.StatusInternalServerError)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

/* ------------------ GetClassCategory ------------------ */
// 200
func TestGetClassCategory_OK(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)

		table := "class_categories"
		userID := uint(42)
		classCategoryID := uint(7)
		classID := uint(5)

		ExpAuthUser(userID, false, false, false)(mock)
		ExpSelectByIDFound(table, classCategoryID, []string{"id"}, []any{classCategoryID})(mock)
		ExpPreloadM2M("class_class_categories", "classes", "class_category_id", "class_id", [][2]any{{classCategoryID, classID}})(mock)

		app := setupApp(gdb)
		resp := runHTTP(t, app, httpInput{
			Method: http.MethodGet, Path: fmt.Sprintf("/class_categories/%d", classCategoryID), UserID: &userID,
		})
		wantStatus(t, resp, http.StatusOK)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

// 404
func TestGetClassCategory_NotFound(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)

		table := "class_categories"
		userID := uint(42)
		classCategoryID := uint(999)

		ExpAuthUser(userID, false, false, false)(mock)
		ExpSelectByIDEmpty(table, classCategoryID)(mock)

		app := setupApp(gdb)
		resp := runHTTP(t, app, httpInput{
			Method: http.MethodGet, Path: fmt.Sprintf("/class_categories/%d", classCategoryID), UserID: &userID,
		})
		wantStatus(t, resp, http.StatusNotFound)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

// 500
func TestGetClassCategory_DBError(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)

		table := "class_categories"
		userID := uint(42)
		classCategoryID := uint(7)

		ExpAuthUser(userID, false, false, false)(mock)
		ExpSelectByIDError(table, classCategoryID, fmt.Errorf("select failed"))(mock)

		app := setupApp(gdb)
		resp := runHTTP(t, app, httpInput{
			Method: http.MethodGet, Path: fmt.Sprintf("/class_categories/%d", classCategoryID), UserID: &userID,
		})
		wantStatus(t, resp, http.StatusInternalServerError)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

// 400
func TestGetClassCategory_BadRequest(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)
		userID := uint(42)

		ExpAuthUser(userID, false, false, false)(mock)

		app := setupApp(gdb)
		resp := runHTTP(t, app, httpInput{
			Method: http.MethodGet, Path: "/class_categories/not-an-int", UserID: &userID,
		})
		wantStatus(t, resp, http.StatusBadRequest)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

/* ------------------ UpdateClassCategory ------------------ */

// 200
func TestUpdateClassCategory_OK(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)

		table := "class_categories"

		userID := uint(42)
		classCategoryID := uint(1)

		ExpAuthUser(userID, true, false, false)(mock)
		ExpSelectByIDFound(table, classCategoryID, []string{"id", "class_category"}, []any{classCategoryID, "test"})(mock)
		ExpUpdateOK(table)(mock)

		app := setupApp(gdb)
		payload := models.ClassCategory{
			ClassCategory: "edit test",
		}
		resp := runHTTP(t, app, httpInput{
			Method: http.MethodPut, Path: fmt.Sprintf("/class_categories/%d", classCategoryID),
			Body: jsonBody(payload), ContentType: "application/json", UserID: &userID,
		})
		wantStatus(t, resp, http.StatusOK)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

// 404
func TestUpdateClassCategory_NotFound(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)

		table := "class_categories"
		userID := uint(42)
		classCategoryID := uint(12345)

		ExpAuthUser(userID, true, false, false)(mock)
		ExpSelectByIDEmpty(table, classCategoryID)(mock)

		app := setupApp(gdb)
		resp := runHTTP(t, app, httpInput{
			Method: http.MethodPut, Path: fmt.Sprintf("/class_categories/%d", classCategoryID), UserID: &userID,
		})
		wantStatus(t, resp, http.StatusNotFound)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

// 500
func TestUpdateClassCategory_DBError(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)

		table := "class_categories"
		userID := uint(42)
		classCategoryID := uint(1)

		ExpAuthUser(userID, true, false, false)(mock)
		ExpSelectByIDFound(table, classCategoryID, []string{"id", "class_category"}, []any{classCategoryID, "test"})(mock)
		ExpUpdateError(table, fmt.Errorf("update failed"))(mock)

		app := setupApp(gdb)
		payload := models.ClassCategory{
			ClassCategory: "edit test",
		}
		resp := runHTTP(t, app, httpInput{
			Method: http.MethodPut, Path: fmt.Sprintf("/class_categories/%d", classCategoryID),
			Body: jsonBody(payload), ContentType: "application/json", UserID: &userID,
		})
		wantStatus(t, resp, http.StatusInternalServerError)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

// 400
func TestUpdateClassCategory_BadRequest(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)
		userID := uint(42)

		ExpAuthUser(userID, true, false, false)(mock)

		app := setupApp(gdb)
		resp := runHTTP(t, app, httpInput{
			Method: http.MethodPut, Path: "/class_categories/not-an-int", UserID: &userID,
		})
		wantStatus(t, resp, http.StatusBadRequest)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

/* ------------------ DeleteClassCategory ------------------ */

// 200 (soft delete)
func TestDeleteClassCategory_OK_SoftDelete(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)

		table := "class_categories"
		userID := uint(42)
		classCategoryID := uint(5)

		ExpAuthUser(userID, true, false, false)(mock)
		ExpSelectByIDFound(table, classCategoryID, []string{"id"}, []any{classCategoryID})(mock)
		ExpSoftDeleteOK(table)(mock)

		app := setupApp(gdb)
		resp := runHTTP(t, app, httpInput{
			Method: http.MethodDelete, Path: fmt.Sprintf("/class_categories/%d", classCategoryID), UserID: &userID,
		})
		wantStatus(t, resp, http.StatusOK)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

// 404
func TestDeleteClassCategory_NotFound(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)

		table := "class_categories"
		userID := uint(42)
		classCategoryID := uint(12345)

		ExpAuthUser(userID, true, false, false)(mock)
		ExpSelectByIDEmpty(table, classCategoryID)(mock)

		app := setupApp(gdb)
		resp := runHTTP(t, app, httpInput{
			Method: http.MethodDelete, Path: fmt.Sprintf("/class_categories/%d", classCategoryID), UserID: &userID,
		})
		wantStatus(t, resp, http.StatusNotFound)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

// 500
func TestDeleteClassCategory_DBError(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)

		table := "class_categories"
		userID := uint(42)
		classCategoryID := uint(5)

		ExpAuthUser(userID, true, false, false)(mock)
		ExpSelectByIDFound(table, classCategoryID, []string{"id"}, []any{classCategoryID})(mock)
		ExpSoftDeleteError(table, fmt.Errorf("update failed"))(mock)

		app := setupApp(gdb)
		resp := runHTTP(t, app, httpInput{
			Method: http.MethodDelete, Path: fmt.Sprintf("/class_categories/%d", classCategoryID), UserID: &userID,
		})
		wantStatus(t, resp, http.StatusInternalServerError)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

// 400
func TestDeleteClassCategory_BadRequest(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)
		userID := uint(42)

		ExpAuthUser(userID, true, false, false)(mock)

		app := setupApp(gdb)
		resp := runHTTP(t, app, httpInput{
			Method: http.MethodDelete, Path: "/class_categories/not-an-int", UserID: &userID,
		})
		wantStatus(t, resp, http.StatusBadRequest)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}
