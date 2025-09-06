package handlers

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/a2n2k3p4/tutorium-backend/models"
)

/* ------------------ CreateEnrollment ------------------ */

// 201
func TestCreateEnrollment_OK(t *testing.T) {
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

			table := "enrollments"
			userID := uint(42)
			learnerID := uint(5)
			classSessionID := uint(10)

			ExpAuthUser(userID, false, false, true)(mock)
			ExpInsertReturningID(table, 1)(mock)

			app := setupApp(gdb)

			payload := models.Enrollment{
				LearnerID:        learnerID,
				ClassSessionID:   classSessionID,
				EnrollmentStatus: "Success",
			}

			resp := runHTTP(t, app, httpInput{
				Method:      http.MethodPost,
				Path:        "/enrollments/",
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
func TestCreateEnrollment_BadRequest(t *testing.T) {
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

			ExpAuthUser(userID, false, false, true)(mock)

			app := setupApp(gdb)
			resp := runHTTP(t, app, httpInput{
				Method: http.MethodPost, Path: "/enrollments/",
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
func TestCreateEnrollment_DBError(t *testing.T) {
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

			table := "enrollments"
			userID := uint(42)

			ExpAuthUser(userID, false, false, true)(mock)
			ExpInsertError(table, fmt.Errorf("db insert failed"))(mock)

			app := setupApp(gdb)
			resp := runHTTP(t, app, httpInput{
				Method: http.MethodPost, Path: "/enrollments/",
				Body: []byte(`{}`), ContentType: "application/json", UserID: &userID,
			})
			wantStatus(t, resp, http.StatusInternalServerError)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("unmet expectations: %v", err)
			}
		})
	}
}

/* ------------------ GetEnrollments ------------------ */
// 200
func TestGetEnrollments_OK(t *testing.T) {
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

			ExpAuthUser(userID, false, false, true)(mock)
			ExpListRows("enrollments", []string{"id"}, []any{1}, []any{2})(mock)

			app := setupApp(gdb)
			resp := runHTTP(t, app, httpInput{
				Method: http.MethodGet, Path: "/enrollments/", UserID: &userID,
			})
			wantStatus(t, resp, http.StatusOK)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("unmet expectations: %v", err)
			}
		})
	}
}

// 500
func TestGetEnrollments_DBError(t *testing.T) {
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

			ExpAuthUser(userID, false, false, true)(mock)
			ExpListError("enrollments", fmt.Errorf("select failed"))(mock)

			app := setupApp(gdb)
			resp := runHTTP(t, app, httpInput{
				Method: http.MethodGet, Path: "/enrollments/", UserID: &userID,
			})
			wantStatus(t, resp, http.StatusInternalServerError)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("unmet expectations: %v", err)
			}
		})
	}
}

/* ------------------ GetEnrollment ------------------ */
// 200
func TestGetEnrollment_OK(t *testing.T) {
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

			table := "enrollments"
			userID := uint(42)
			enrollmentID := uint(7)

			ExpAuthUser(userID, false, false, true)(mock)
			ExpSelectByIDFound(table, enrollmentID, []string{"id"}, []any{enrollmentID})(mock)

			app := setupApp(gdb)
			resp := runHTTP(t, app, httpInput{
				Method: http.MethodGet, Path: fmt.Sprintf("/enrollments/%d", enrollmentID), UserID: &userID,
			})
			wantStatus(t, resp, http.StatusOK)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("unmet expectations: %v", err)
			}
		})
	}
}

// 404
func TestGetEnrollment_NotFound(t *testing.T) {
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

			table := "enrollments"
			userID := uint(42)
			enrollmentID := uint(999)

			ExpAuthUser(userID, false, false, true)(mock)
			ExpSelectByIDEmpty(table, enrollmentID)(mock)

			app := setupApp(gdb)
			resp := runHTTP(t, app, httpInput{
				Method: http.MethodGet, Path: fmt.Sprintf("/enrollments/%d", enrollmentID), UserID: &userID,
			})
			wantStatus(t, resp, http.StatusNotFound)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("unmet expectations: %v", err)
			}
		})
	}
}

// 500
func TestGetEnrollment_DBError(t *testing.T) {
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

			table := "enrollments"
			userID := uint(42)
			enrollmentID := uint(7)

			ExpAuthUser(userID, false, false, true)(mock)
			ExpSelectByIDError(table, enrollmentID, fmt.Errorf("select failed"))(mock)

			app := setupApp(gdb)
			resp := runHTTP(t, app, httpInput{
				Method: http.MethodGet, Path: fmt.Sprintf("/enrollments/%d", enrollmentID), UserID: &userID,
			})
			wantStatus(t, resp, http.StatusInternalServerError)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("unmet expectations: %v", err)
			}
		})
	}
}

// 400
func TestGetEnrollment_BadRequest(t *testing.T) {
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

			ExpAuthUser(userID, false, false, true)(mock)

			app := setupApp(gdb)
			resp := runHTTP(t, app, httpInput{
				Method: http.MethodGet, Path: "/enrollments/not-an-int", UserID: &userID,
			})
			wantStatus(t, resp, http.StatusBadRequest)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("unmet expectations: %v", err)
			}
		})
	}
}

/* ------------------ UpdateEnrollment ------------------ */

// 200
func TestUpdateEnrollment_OK(t *testing.T) {
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

			table := "enrollments"
			preloadTable := "learners"
			preloadTable2 := "class_sessions"
			userID := uint(42)
			enrollmentID := uint(1)
			learnerID := uint(5)
			classSessionID := uint(10)

			ExpAuthUser(userID, false, false, true)(mock)
			ExpSelectByIDFound(table, enrollmentID,
				[]string{"id", "learner_id", "class_session_id", "enrollment_status"},
				[]any{enrollmentID, learnerID, classSessionID, "pending"},
			)(mock)

			ExpPreloadField(preloadTable, []string{"id"}, []any{learnerID})(mock)
			ExpPreloadField(preloadTable2, []string{"id"}, []any{classSessionID})(mock)

			ExpUpdateOK(table)(mock)

			app := setupApp(gdb)
			payload := models.Enrollment{
				LearnerID:        learnerID,
				ClassSessionID:   classSessionID,
				EnrollmentStatus: "Success",
			}
			resp := runHTTP(t, app, httpInput{
				Method: http.MethodPut, Path: fmt.Sprintf("/enrollments/%d", enrollmentID),
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
func TestUpdateEnrollment_NotFound(t *testing.T) {
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

			table := "enrollments"
			userID := uint(42)
			enrollmentID := uint(12345)

			ExpAuthUser(userID, false, false, true)(mock)
			ExpSelectByIDEmpty(table, enrollmentID)(mock)

			app := setupApp(gdb)
			resp := runHTTP(t, app, httpInput{
				Method: http.MethodPut, Path: fmt.Sprintf("/enrollments/%d", enrollmentID), UserID: &userID,
			})
			wantStatus(t, resp, http.StatusNotFound)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("unmet expectations: %v", err)
			}
		})
	}
}

// 500
func TestUpdateEnrollment_DBError(t *testing.T) {
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

			table := "enrollments"
			userID := uint(42)
			enrollmentID := uint(1)
			learnerID := uint(5)
			classSessionID := uint(10)

			ExpAuthUser(userID, false, false, true)(mock)
			ExpSelectByIDFound(table, enrollmentID, []string{"id"}, []any{enrollmentID})(mock)
			ExpUpdateError(table, fmt.Errorf("update failed"))(mock)

			app := setupApp(gdb)
			payload := models.Enrollment{
				LearnerID:        learnerID,
				ClassSessionID:   classSessionID,
				EnrollmentStatus: "Success",
			}
			resp := runHTTP(t, app, httpInput{
				Method: http.MethodPut, Path: fmt.Sprintf("/enrollments/%d", enrollmentID),
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
func TestUpdateEnrollment_BadRequest(t *testing.T) {
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

			ExpAuthUser(userID, false, false, true)(mock)

			app := setupApp(gdb)
			resp := runHTTP(t, app, httpInput{
				Method: http.MethodPut, Path: "/enrollments/not-an-int", UserID: &userID,
			})
			wantStatus(t, resp, http.StatusBadRequest)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("unmet expectations: %v", err)
			}
		})
	}
}

/* ------------------ DeleteEnrollment ------------------ */

// 200 (soft delete)
func TestDeleteEnrollment_OK_SoftDelete(t *testing.T) {
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

			table := "enrollments"
			userID := uint(42)
			enrollmentID := uint(5)

			ExpAuthUser(userID, false, false, true)(mock)
			ExpSelectByIDFound(table, enrollmentID, []string{"id"}, []any{enrollmentID})(mock)
			ExpSoftDeleteOK(table)(mock)

			app := setupApp(gdb)
			resp := runHTTP(t, app, httpInput{
				Method: http.MethodDelete, Path: fmt.Sprintf("/enrollments/%d", enrollmentID), UserID: &userID,
			})
			wantStatus(t, resp, http.StatusOK)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("unmet expectations: %v", err)
			}
		})
	}
}

// 404
func TestDeleteEnrollment_NotFound(t *testing.T) {
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

			table := "enrollments"
			userID := uint(42)
			enrollmentID := uint(12345)

			ExpAuthUser(userID, false, false, true)(mock)
			ExpSelectByIDEmpty(table, enrollmentID)(mock)

			app := setupApp(gdb)
			resp := runHTTP(t, app, httpInput{
				Method: http.MethodDelete, Path: fmt.Sprintf("/enrollments/%d", enrollmentID), UserID: &userID,
			})
			wantStatus(t, resp, http.StatusNotFound)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("unmet expectations: %v", err)
			}
		})
	}
}

// 500
func TestDeleteEnrollment_DBError(t *testing.T) {
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

			table := "enrollments"
			userID := uint(42)
			enrollmentID := uint(5)

			ExpAuthUser(userID, false, false, true)(mock)
			ExpSelectByIDFound(table, enrollmentID, []string{"id"}, []any{enrollmentID})(mock)
			ExpSoftDeleteError(table, fmt.Errorf("update failed"))(mock)

			app := setupApp(gdb)
			resp := runHTTP(t, app, httpInput{
				Method: http.MethodDelete, Path: fmt.Sprintf("/enrollments/%d", enrollmentID), UserID: &userID,
			})
			wantStatus(t, resp, http.StatusInternalServerError)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("unmet expectations: %v", err)
			}
		})
	}
}

// 400
func TestDeleteEnrollment_BadRequest(t *testing.T) {
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

			ExpAuthUser(userID, false, false, true)(mock)

			app := setupApp(gdb)
			resp := runHTTP(t, app, httpInput{
				Method: http.MethodDelete, Path: "/enrollments/not-an-int", UserID: &userID,
			})
			wantStatus(t, resp, http.StatusBadRequest)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("unmet expectations: %v", err)
			}
		})
	}
}
