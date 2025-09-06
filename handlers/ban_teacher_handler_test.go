package handlers

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/a2n2k3p4/tutorium-backend/models"
)

/* ------------------ CreateBanTeacher ------------------ */

// 201
func TestCreateBanTeacher_OK(t *testing.T) {
	cases := []struct {
		name       string
		STATUS_env string
	}{
		{"bypass", "development"},
		{"unbypass", "production"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Setenv("STATUS", c.STATUS_env)
			mock, gdb, cleanup := setupMockGorm(t)
			defer cleanup()
			mock.MatchExpectationsInOrder(false)

			table := "ban_details_teachers"
			userID := uint(42)
			teacherID := uint(50)
			now := time.Now()

			ExpAuthUser(userID, true, false, false)(mock)
			ExpInsertReturningID(table, 1)(mock)

			app := setupApp(gdb)

			payload := models.BanDetailsTeacher{
				TeacherID:      teacherID,
				BanStart:       now,
				BanEnd:         now.Add(2 * time.Hour),
				BanDescription: "spamming",
			}

			resp := runHTTP(t, app, httpInput{
				Method:      http.MethodPost,
				Path:        "/banteachers/",
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
}

// 400
func TestCreateBanTeacher_BadRequest(t *testing.T) {
	cases := []struct {
		name       string
		STATUS_env string
	}{
		{"bypass", "development"},
		{"unbypass", "production"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Setenv("STATUS", c.STATUS_env)
			mock, gdb, cleanup := setupMockGorm(t)
			defer cleanup()
			mock.MatchExpectationsInOrder(false)
			userID := uint(42)

			ExpAuthUser(userID, true, false, false)(mock)

			app := setupApp(gdb)
			resp := runHTTP(t, app, httpInput{
				Method: http.MethodPost, Path: "/banteachers/",
				Body: []byte(`{invalid-json}`), ContentType: "application/json", UserID: &userID,
			})
			wantStatus(t, resp, http.StatusBadRequest)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("unmet expectations: %v", err)
			}
		})
	}
}

// 500
func TestCreateBanTeacher_DBError(t *testing.T) {
	cases := []struct {
		name       string
		STATUS_env string
	}{
		{"bypass", "development"},
		{"unbypass", "production"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Setenv("STATUS", c.STATUS_env)
			mock, gdb, cleanup := setupMockGorm(t)
			defer cleanup()
			mock.MatchExpectationsInOrder(false)

			table := "ban_details_teachers"
			userID := uint(42)
			teacherID := uint(50)
			now := time.Now()

			ExpAuthUser(userID, true, false, false)(mock)
			ExpInsertError(table, fmt.Errorf("db insert failed"))(mock)

			payload := models.BanDetailsTeacher{
				TeacherID:      teacherID,
				BanStart:       now,
				BanEnd:         now.Add(2 * time.Hour),
				BanDescription: "spamming",
			}
			app := setupApp(gdb)
			resp := runHTTP(t, app, httpInput{
				Method: http.MethodPost, Path: "/banteachers/",
				Body: jsonBody(payload), ContentType: "application/json", UserID: &userID,
			})
			wantStatus(t, resp, http.StatusInternalServerError)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("unmet expectations: %v", err)
			}
		})
	}
}

/* ------------------ GetBanTeachers ------------------ */
// 200
func TestGetBanTeachers_OK(t *testing.T) {
	cases := []struct {
		name       string
		STATUS_env string
	}{
		{"bypass", "development"},
		{"unbypass", "production"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Setenv("STATUS", c.STATUS_env)
			mock, gdb, cleanup := setupMockGorm(t)
			defer cleanup()
			mock.MatchExpectationsInOrder(false)
			userID := uint(42)

			ExpAuthUser(userID, true, false, false)(mock)
			ExpListRows("ban_details_teachers", []string{"id"}, []any{1}, []any{2})(mock)

			app := setupApp(gdb)
			resp := runHTTP(t, app, httpInput{
				Method: http.MethodGet, Path: "/banteachers/", UserID: &userID,
			})
			wantStatus(t, resp, http.StatusOK)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("unmet expectations: %v", err)
			}
		})
	}
}

// 500
func TestGetBanTeachers_DBError(t *testing.T) {
	cases := []struct {
		name       string
		STATUS_env string
	}{
		{"bypass", "development"},
		{"unbypass", "production"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Setenv("STATUS", c.STATUS_env)
			mock, gdb, cleanup := setupMockGorm(t)
			defer cleanup()
			mock.MatchExpectationsInOrder(false)
			userID := uint(42)

			ExpAuthUser(userID, true, false, false)(mock)
			ExpListError("ban_details_teachers", fmt.Errorf("select failed"))(mock)

			app := setupApp(gdb)
			resp := runHTTP(t, app, httpInput{
				Method: http.MethodGet, Path: "/banteachers/", UserID: &userID,
			})
			wantStatus(t, resp, http.StatusInternalServerError)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("unmet expectations: %v", err)
			}
		})
	}
}

/* ------------------ GetBanTeacher ------------------ */
// 200
func TestGetBanTeacher_OK(t *testing.T) {
	cases := []struct {
		name       string
		STATUS_env string
	}{
		{"bypass", "development"},
		{"unbypass", "production"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Setenv("STATUS", c.STATUS_env)
			mock, gdb, cleanup := setupMockGorm(t)
			defer cleanup()
			mock.MatchExpectationsInOrder(false)

			table := "ban_details_teachers"
			userID := uint(42)
			ban_teacherID := uint(7)

			ExpAuthUser(userID, true, false, false)(mock)
			ExpSelectByIDFound(table, ban_teacherID, []string{"id"}, []any{ban_teacherID})(mock)

			app := setupApp(gdb)
			resp := runHTTP(t, app, httpInput{
				Method: http.MethodGet, Path: fmt.Sprintf("/banteachers/%d", ban_teacherID), UserID: &userID,
			})
			wantStatus(t, resp, http.StatusOK)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("unmet expectations: %v", err)
			}
		})
	}
}

// 404
func TestGetBanTeacher_NotFound(t *testing.T) {
	cases := []struct {
		name       string
		STATUS_env string
	}{
		{"bypass", "development"},
		{"unbypass", "production"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Setenv("STATUS", c.STATUS_env)
			mock, gdb, cleanup := setupMockGorm(t)
			defer cleanup()
			mock.MatchExpectationsInOrder(false)

			table := "ban_details_teachers"
			userID := uint(42)
			ban_teacherID := uint(999)

			ExpAuthUser(userID, true, false, false)(mock)
			ExpSelectByIDEmpty(table, ban_teacherID)(mock)

			app := setupApp(gdb)
			resp := runHTTP(t, app, httpInput{
				Method: http.MethodGet, Path: fmt.Sprintf("/banteachers/%d", ban_teacherID), UserID: &userID,
			})
			wantStatus(t, resp, http.StatusNotFound)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("unmet expectations: %v", err)
			}
		})
	}
}

// 500
func TestGetBanTeacher_DBError(t *testing.T) {
	cases := []struct {
		name       string
		STATUS_env string
	}{
		{"bypass", "development"},
		{"unbypass", "production"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Setenv("STATUS", c.STATUS_env)
			mock, gdb, cleanup := setupMockGorm(t)
			defer cleanup()
			mock.MatchExpectationsInOrder(false)

			table := "ban_details_teachers"
			userID := uint(42)
			ban_teacherID := uint(7)

			ExpAuthUser(userID, true, false, false)(mock)
			ExpSelectByIDError(table, ban_teacherID, fmt.Errorf("select failed"))(mock)

			app := setupApp(gdb)
			resp := runHTTP(t, app, httpInput{
				Method: http.MethodGet, Path: fmt.Sprintf("/banteachers/%d", ban_teacherID), UserID: &userID,
			})
			wantStatus(t, resp, http.StatusInternalServerError)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("unmet expectations: %v", err)
			}
		})
	}
}

// 400
func TestGetBanTeacher_BadRequest(t *testing.T) {
	cases := []struct {
		name       string
		STATUS_env string
	}{
		{"bypass", "development"},
		{"unbypass", "production"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Setenv("STATUS", c.STATUS_env)
			mock, gdb, cleanup := setupMockGorm(t)
			defer cleanup()
			mock.MatchExpectationsInOrder(false)
			userID := uint(42)

			ExpAuthUser(userID, true, false, false)(mock)

			app := setupApp(gdb)
			resp := runHTTP(t, app, httpInput{
				Method: http.MethodGet, Path: "/banteachers/not-an-int", UserID: &userID,
			})
			wantStatus(t, resp, http.StatusBadRequest)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("unmet expectations: %v", err)
			}
		})
	}
}

/* ------------------ UpdateBanTeacher ------------------ */

// 200
func TestUpdateBanTeacher_OK(t *testing.T) {
	cases := []struct {
		name       string
		STATUS_env string
	}{
		{"bypass", "development"},
		{"unbypass", "production"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Setenv("STATUS", c.STATUS_env)
			mock, gdb, cleanup := setupMockGorm(t)
			defer cleanup()
			mock.MatchExpectationsInOrder(false)

			table := "ban_details_teachers"
			const preloadTable = "teachers"

			userID := uint(42)
			teacherID := uint(50)
			ban_teacherID := uint(1)
			now := time.Now()

			ExpAuthUser(userID, true, false, false)(mock)
			ExpSelectByIDFound(table, ban_teacherID,
				[]string{"id", "teacher_id", "ban_start", "ban_end", "ban_description"},
				[]any{ban_teacherID, teacherID, now, now.Add(2 * time.Hour), "flooding"},
			)(mock)

			ExpPreloadField(preloadTable, []string{"id"}, []any{teacherID})(mock)

			ExpUpdateOK(table)(mock)

			app := setupApp(gdb)
			payload := models.BanDetailsTeacher{
				BanDescription: "spamming",
			}
			resp := runHTTP(t, app, httpInput{
				Method: http.MethodPut, Path: fmt.Sprintf("/banteachers/%d", ban_teacherID),
				Body: jsonBody(payload), ContentType: "application/json", UserID: &userID,
			})
			wantStatus(t, resp, http.StatusOK)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("unmet expectations: %v", err)
			}
		})
	}
}

// 404
func TestUpdateBanTeacher_NotFound(t *testing.T) {
	cases := []struct {
		name       string
		STATUS_env string
	}{
		{"bypass", "development"},
		{"unbypass", "production"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Setenv("STATUS", c.STATUS_env)
			mock, gdb, cleanup := setupMockGorm(t)
			defer cleanup()
			mock.MatchExpectationsInOrder(false)

			table := "ban_details_teachers"
			userID := uint(42)
			ban_teacherID := uint(12345)

			ExpAuthUser(userID, true, false, false)(mock)
			ExpSelectByIDEmpty(table, ban_teacherID)(mock)

			app := setupApp(gdb)
			resp := runHTTP(t, app, httpInput{
				Method: http.MethodPut, Path: fmt.Sprintf("/banteachers/%d", ban_teacherID), UserID: &userID,
			})
			wantStatus(t, resp, http.StatusNotFound)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("unmet expectations: %v", err)
			}
		})
	}
}

// 500
func TestUpdateBanTeacher_DBError(t *testing.T) {
	cases := []struct {
		name       string
		STATUS_env string
	}{
		{"bypass", "development"},
		{"unbypass", "production"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Setenv("STATUS", c.STATUS_env)
			mock, gdb, cleanup := setupMockGorm(t)
			defer cleanup()
			mock.MatchExpectationsInOrder(false)

			table := "ban_details_teachers"

			userID := uint(42)

			ban_teacherID := uint(1)

			ExpAuthUser(userID, true, false, false)(mock)
			ExpSelectByIDFound(table, ban_teacherID, []string{"id"}, []any{ban_teacherID})(mock)
			ExpUpdateError(table, fmt.Errorf("update failed"))(mock)

			app := setupApp(gdb)
			payload := models.BanDetailsTeacher{

				BanDescription: "test err",
			}
			resp := runHTTP(t, app, httpInput{
				Method: http.MethodPut, Path: fmt.Sprintf("/banteachers/%d", ban_teacherID),
				Body: jsonBody(payload), ContentType: "application/json", UserID: &userID,
			})
			wantStatus(t, resp, http.StatusInternalServerError)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("unmet expectations: %v", err)
			}
		})
	}
}

// 400
func TestUpdateBanTeacher_BadRequest(t *testing.T) {
	cases := []struct {
		name       string
		STATUS_env string
	}{
		{"bypass", "development"},
		{"unbypass", "production"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Setenv("STATUS", c.STATUS_env)
			mock, gdb, cleanup := setupMockGorm(t)
			defer cleanup()
			mock.MatchExpectationsInOrder(false)
			userID := uint(42)

			ExpAuthUser(userID, true, false, false)(mock)

			app := setupApp(gdb)
			resp := runHTTP(t, app, httpInput{
				Method: http.MethodPut, Path: "/banteachers/not-an-int", UserID: &userID,
			})
			wantStatus(t, resp, http.StatusBadRequest)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("unmet expectations: %v", err)
			}
		})
	}
}

/* ------------------ DeleteBanTeacher ------------------ */

// 200
func TestDeleteBanTeacher_OK_SoftDelete(t *testing.T) {
	cases := []struct {
		name       string
		STATUS_env string
	}{
		{"bypass", "development"},
		{"unbypass", "production"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Setenv("STATUS", c.STATUS_env)
			mock, gdb, cleanup := setupMockGorm(t)
			defer cleanup()
			mock.MatchExpectationsInOrder(false)

			table := "ban_details_teachers"
			userID := uint(42)
			ban_teacherID := uint(5)

			ExpAuthUser(userID, true, false, false)(mock)
			ExpSelectByIDFound(table, ban_teacherID, []string{"id"}, []any{ban_teacherID})(mock)
			ExpSoftDeleteOK(table)(mock)

			app := setupApp(gdb)
			resp := runHTTP(t, app, httpInput{
				Method: http.MethodDelete, Path: fmt.Sprintf("/banteachers/%d", ban_teacherID), UserID: &userID,
			})
			wantStatus(t, resp, http.StatusOK)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("unmet expectations: %v", err)
			}
		})
	}
}

// 404
func TestDeleteBanTeacher_NotFound(t *testing.T) {
	cases := []struct {
		name       string
		STATUS_env string
	}{
		{"bypass", "development"},
		{"unbypass", "production"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Setenv("STATUS", c.STATUS_env)
			mock, gdb, cleanup := setupMockGorm(t)
			defer cleanup()
			mock.MatchExpectationsInOrder(false)

			table := "ban_details_teachers"
			userID := uint(42)
			ban_teacherID := uint(12345)

			ExpAuthUser(userID, true, false, false)(mock)
			ExpSelectByIDEmpty(table, ban_teacherID)(mock)

			app := setupApp(gdb)
			resp := runHTTP(t, app, httpInput{
				Method: http.MethodDelete, Path: fmt.Sprintf("/banteachers/%d", ban_teacherID), UserID: &userID,
			})
			wantStatus(t, resp, http.StatusNotFound)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("unmet expectations: %v", err)
			}
		})
	}
}

// 500
func TestDeleteBanTeacher_DBError(t *testing.T) {
	cases := []struct {
		name       string
		STATUS_env string
	}{
		{"bypass", "development"},
		{"unbypass", "production"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Setenv("STATUS", c.STATUS_env)
			mock, gdb, cleanup := setupMockGorm(t)
			defer cleanup()
			mock.MatchExpectationsInOrder(false)

			table := "ban_details_teachers"
			userID := uint(42)
			ban_teacherID := uint(5)

			ExpAuthUser(userID, true, false, false)(mock)
			ExpSelectByIDFound(table, ban_teacherID, []string{"id"}, []any{ban_teacherID})(mock)
			ExpSoftDeleteError(table, fmt.Errorf("update failed"))(mock)

			app := setupApp(gdb)
			resp := runHTTP(t, app, httpInput{
				Method: http.MethodDelete, Path: fmt.Sprintf("/banteachers/%d", ban_teacherID), UserID: &userID,
			})
			wantStatus(t, resp, http.StatusInternalServerError)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("unmet expectations: %v", err)
			}
		})
	}
}

// 400
func TestDeleteBanTeacher_BadRequest(t *testing.T) {
	cases := []struct {
		name       string
		STATUS_env string
	}{
		{"bypass", "development"},
		{"unbypass", "production"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Setenv("STATUS", c.STATUS_env)
			mock, gdb, cleanup := setupMockGorm(t)
			defer cleanup()
			mock.MatchExpectationsInOrder(false)
			userID := uint(42)

			ExpAuthUser(userID, true, false, false)(mock)

			app := setupApp(gdb)
			resp := runHTTP(t, app, httpInput{
				Method: http.MethodDelete, Path: "/banteachers/not-an-int", UserID: &userID,
			})
			wantStatus(t, resp, http.StatusBadRequest)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("unmet expectations: %v", err)
			}
		})
	}
}
