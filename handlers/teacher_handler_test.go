package handlers

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/a2n2k3p4/tutorium-backend/models"
)

/* ------------------ CreateTeacher ------------------ */

// 201
func TestCreateTeacher_OK(t *testing.T) {
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

			table := "teachers"
			userID := uint(42)

			ExpAuthUser(userID, false, false, false)(mock)
			ExpInsertReturningID(table, 1)(mock)

			app := setupApp(gdb)

			payload := models.Teacher{
				UserID: userID,
			}

			resp := runHTTP(t, app, httpInput{
				Method:      http.MethodPost,
				Path:        "/teachers/",
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
func TestCreateTeacher_BadRequest(t *testing.T) {
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

			ExpAuthUser(userID, false, false, false)(mock)

			app := setupApp(gdb)
			resp := runHTTP(t, app, httpInput{
				Method: http.MethodPost, Path: "/teachers/",
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
func TestCreateTeacher_DBError(t *testing.T) {
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

			table := "teachers"
			userID := uint(42)

			ExpAuthUser(userID, false, false, false)(mock)
			ExpInsertError(table, fmt.Errorf("db insert failed"))(mock)

			payload := models.Teacher{
				UserID: userID,
			}
			app := setupApp(gdb)
			resp := runHTTP(t, app, httpInput{
				Method: http.MethodPost, Path: "/teachers/",
				Body: jsonBody(payload), ContentType: "application/json", UserID: &userID,
			})
			wantStatus(t, resp, http.StatusInternalServerError)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("unmet expectations: %v", err)
			}
		})
	}
}

/* ------------------ GetTeachers ------------------ */
// 200
func TestGetTeachers_OK(t *testing.T) {
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

			ExpAuthUser(userID, false, false, false)(mock)
			ExpListRows("teachers", []string{"id"}, []any{1}, []any{2})(mock)

			app := setupApp(gdb)
			resp := runHTTP(t, app, httpInput{
				Method: http.MethodGet, Path: "/teachers/", UserID: &userID,
			})
			wantStatus(t, resp, http.StatusOK)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("unmet expectations: %v", err)
			}
		})
	}
}

// 500
func TestGetTeachers_DBError(t *testing.T) {
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

			ExpAuthUser(userID, false, false, false)(mock)
			ExpListError("teachers", fmt.Errorf("select failed"))(mock)

			app := setupApp(gdb)
			resp := runHTTP(t, app, httpInput{
				Method: http.MethodGet, Path: "/teachers/", UserID: &userID,
			})
			wantStatus(t, resp, http.StatusInternalServerError)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("unmet expectations: %v", err)
			}
		})
	}
}

/* ------------------ GetTeacher ------------------ */
// 200
func TestGetTeacher_OK(t *testing.T) {
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

			table := "teachers"
			userID := uint(42)
			teacherID := uint(7)

			ExpAuthUser(userID, false, false, false)(mock)
			ExpSelectByIDFound(table, teacherID, []string{"id"}, []any{teacherID})(mock)

			app := setupApp(gdb)
			resp := runHTTP(t, app, httpInput{
				Method: http.MethodGet, Path: fmt.Sprintf("/teachers/%d", teacherID), UserID: &userID,
			})
			wantStatus(t, resp, http.StatusOK)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("unmet expectations: %v", err)
			}
		})
	}
}

// 404
func TestGetTeacher_NotFound(t *testing.T) {
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

			table := "teachers"
			userID := uint(42)
			teacherID := uint(999)

			ExpAuthUser(userID, false, false, false)(mock)
			ExpSelectByIDEmpty(table, teacherID)(mock)

			app := setupApp(gdb)
			resp := runHTTP(t, app, httpInput{
				Method: http.MethodGet, Path: fmt.Sprintf("/teachers/%d", teacherID), UserID: &userID,
			})
			wantStatus(t, resp, http.StatusNotFound)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("unmet expectations: %v", err)
			}
		})
	}
}

// 500
func TestGetTeacher_DBError(t *testing.T) {
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

			table := "teachers"
			userID := uint(42)
			teacherID := uint(7)

			ExpAuthUser(userID, false, false, false)(mock)
			ExpSelectByIDError(table, teacherID, fmt.Errorf("select failed"))(mock)

			app := setupApp(gdb)
			resp := runHTTP(t, app, httpInput{
				Method: http.MethodGet, Path: fmt.Sprintf("/teachers/%d", teacherID), UserID: &userID,
			})
			wantStatus(t, resp, http.StatusInternalServerError)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("unmet expectations: %v", err)
			}
		})
	}
}

// 400
func TestGetTeacher_BadRequest(t *testing.T) {
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

			ExpAuthUser(userID, false, false, false)(mock)

			app := setupApp(gdb)
			resp := runHTTP(t, app, httpInput{
				Method: http.MethodGet, Path: "/teachers/not-an-int", UserID: &userID,
			})
			wantStatus(t, resp, http.StatusBadRequest)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("unmet expectations: %v", err)
			}
		})
	}
}

/* ------------------ DeleteTeacher ------------------ */

// 200
func TestDeleteTeacher_OK_SoftDelete(t *testing.T) {
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

			table := "teachers"
			userID := uint(42)
			teacherID := uint(5)

			ExpAuthUser(userID, false, false, false)(mock)
			ExpSelectByIDFound(table, teacherID, []string{"id"}, []any{teacherID})(mock)
			ExpSoftDeleteOK(table)(mock)

			app := setupApp(gdb)
			resp := runHTTP(t, app, httpInput{
				Method: http.MethodDelete, Path: fmt.Sprintf("/teachers/%d", teacherID), UserID: &userID,
			})
			wantStatus(t, resp, http.StatusOK)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("unmet expectations: %v", err)
			}
		})
	}
}

// 404
func TestDeleteTeacher_NotFound(t *testing.T) {
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

			table := "teachers"
			userID := uint(42)
			teacherID := uint(12345)

			ExpAuthUser(userID, false, false, false)(mock)
			ExpSelectByIDEmpty(table, teacherID)(mock)

			app := setupApp(gdb)
			resp := runHTTP(t, app, httpInput{
				Method: http.MethodDelete, Path: fmt.Sprintf("/teachers/%d", teacherID), UserID: &userID,
			})
			wantStatus(t, resp, http.StatusNotFound)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("unmet expectations: %v", err)
			}
		})
	}
}

// 500
func TestDeleteTeacher_DBError(t *testing.T) {
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

			table := "teachers"
			userID := uint(42)
			teacherID := uint(5)

			ExpAuthUser(userID, false, false, false)(mock)
			ExpSelectByIDFound(table, teacherID, []string{"id"}, []any{teacherID})(mock)
			ExpSoftDeleteError(table, fmt.Errorf("update failed"))(mock)

			app := setupApp(gdb)
			resp := runHTTP(t, app, httpInput{
				Method: http.MethodDelete, Path: fmt.Sprintf("/teachers/%d", teacherID), UserID: &userID,
			})
			wantStatus(t, resp, http.StatusInternalServerError)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("unmet expectations: %v", err)
			}
		})
	}
}

// 400
func TestDeleteTeacher_BadRequest(t *testing.T) {
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

			ExpAuthUser(userID, false, false, false)(mock)

			app := setupApp(gdb)
			resp := runHTTP(t, app, httpInput{
				Method: http.MethodDelete, Path: "/teachers/not-an-int", UserID: &userID,
			})
			wantStatus(t, resp, http.StatusBadRequest)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("unmet expectations: %v", err)
			}
		})
	}
}
