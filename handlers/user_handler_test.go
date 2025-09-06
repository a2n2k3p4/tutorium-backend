package handlers

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/a2n2k3p4/tutorium-backend/models"
	"github.com/stretchr/testify/require"
)

/* ------------------ CreateUser ------------------ */

// 201
func TestCreateUser_OK(t *testing.T) {
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

			table := "users"
			userID := uint(42)
			studentID := "6600000000"

			ExpAuthUser(userID, false, false, false)(mock)
			ExpDoubleInsertReturningID(table, "learners", uint64(userID), 2)(mock)

			app := setupApp(gdb)

			payload := models.User{
				StudentID:         studentID,
				ProfilePictureURL: "",
				FirstName:         "Jane",
				LastName:          "Doe",
				Gender:            "Female",
				PhoneNumber:       "",
				Balance:           0,
			}

			resp := runHTTP(t, app, httpInput{
				Method:      http.MethodPost,
				Path:        "/users/",
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
func TestCreateUser_BadRequest(t *testing.T) {
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
				Method: http.MethodPost, Path: "/users/",
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
func TestCreateUser_DBError(t *testing.T) {
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

			table := "users"
			userID := uint(42)
			studentID := "6600000000"

			ExpAuthUser(userID, false, false, false)(mock)
			ExpInsertError(table, fmt.Errorf("db insert failed"))(mock)

			payload := models.User{
				StudentID:         studentID,
				ProfilePictureURL: "",
				FirstName:         "Jane",
				LastName:          "Doe",
				Gender:            "Female",
				PhoneNumber:       "",
				Balance:           0,
			}
			app := setupApp(gdb)
			resp := runHTTP(t, app, httpInput{
				Method: http.MethodPost, Path: "/users/",
				Body: jsonBody(payload), ContentType: "application/json", UserID: &userID,
			})
			wantStatus(t, resp, http.StatusInternalServerError)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("unmet expectations: %v", err)
			}
		})
	}
}

/* ------------------ GetUsers ------------------ */
// 200
func TestGetUsers_OK(t *testing.T) {
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
			ExpListRows("users", []string{"id"}, []any{1}, []any{2})(mock)
			ExpPreloadCanEmpty("learners", []string{"id", "user_id"})(mock)
			ExpPreloadCanEmpty("teachers", []string{"id", "user_id"})(mock)
			ExpPreloadCanEmpty("admins", []string{"id", "user_id"})(mock)

			app := setupApp(gdb)
			resp := runHTTP(t, app, httpInput{
				Method: http.MethodGet, Path: "/users/", UserID: &userID,
			})
			wantStatus(t, resp, http.StatusOK)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("unmet expectations: %v", err)
			}
		})
	}
}

// 500
func TestGetUsers_DBError(t *testing.T) {
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
			ExpListError("users", fmt.Errorf("select failed"))(mock)

			app := setupApp(gdb)
			resp := runHTTP(t, app, httpInput{
				Method: http.MethodGet, Path: "/users/", UserID: &userID,
			})
			wantStatus(t, resp, http.StatusInternalServerError)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("unmet expectations: %v", err)
			}
		})
	}
}

/* ------------------ GetUser ------------------ */
// 200
func TestGetUser_OK(t *testing.T) {
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

			table := "users"
			userID := uint(42)

			ExpAuthUser(userID, false, false, false)(mock)
			ExpSelectByIDFound(table, userID, []string{"id"}, []any{userID})(mock)
			ExpPreloadCanEmpty("learners", []string{"id", "user_id"})(mock)
			ExpPreloadCanEmpty("teachers", []string{"id", "user_id"})(mock)
			ExpPreloadCanEmpty("admins", []string{"id", "user_id"})(mock)

			app := setupApp(gdb)
			resp := runHTTP(t, app, httpInput{
				Method: http.MethodGet, Path: fmt.Sprintf("/users/%d", userID), UserID: &userID,
			})
			wantStatus(t, resp, http.StatusOK)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("unmet expectations: %v", err)
			}
		})
	}
}

// 404
func TestGetUser_NotFound(t *testing.T) {
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

			table := "users"
			userID := uint(42)

			ExpAuthUser(userID, false, false, false)(mock)
			ExpSelectByIDEmpty(table, userID)(mock)

			app := setupApp(gdb)
			resp := runHTTP(t, app, httpInput{
				Method: http.MethodGet, Path: fmt.Sprintf("/users/%d", userID), UserID: &userID,
			})
			wantStatus(t, resp, http.StatusNotFound)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("unmet expectations: %v", err)
			}
		})
	}
}

// 500
func TestGetUser_DBError(t *testing.T) {
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

			table := "users"
			userID := uint(42)

			ExpAuthUser(userID, false, false, false)(mock)
			ExpSelectByIDError(table, userID, fmt.Errorf("select failed"))(mock)

			app := setupApp(gdb)
			resp := runHTTP(t, app, httpInput{
				Method: http.MethodGet, Path: fmt.Sprintf("/users/%d", userID), UserID: &userID,
			})
			wantStatus(t, resp, http.StatusInternalServerError)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("unmet expectations: %v", err)
			}
		})
	}
}

// 400
func TestGetUser_BadRequest(t *testing.T) {
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
				Method: http.MethodGet, Path: "/users/not-an-int", UserID: &userID,
			})
			wantStatus(t, resp, http.StatusBadRequest)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("unmet expectations: %v", err)
			}
		})
	}
}

/* ------------------ UpdateUser ------------------ */

// 200
func TestUpdateUser_OK(t *testing.T) {
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

			table := "users"
			userID := uint(42)
			studentID := "6600000000"

			ExpAuthUser(userID, false, false, false)(mock)
			ExpSelectByIDFound(table, userID,
				[]string{"id", "student_id", "profilePicture_url", "first_name", "last_name", "gender", "phone_number", "balance"},
				[]any{userID, studentID, "", "Janet", "Doe", "Female", "", 50},
			)(mock)

			ExpPreloadCanEmpty("learners", []string{"id", "user_id"})(mock)
			ExpPreloadCanEmpty("teachers", []string{"id", "user_id"})(mock)
			ExpPreloadCanEmpty("admins", []string{"id", "user_id"})(mock)

			ExpUpdateOK(table)(mock)

			app := setupApp(gdb)
			payload := models.User{
				StudentID:         studentID,
				ProfilePictureURL: "",
				FirstName:         "Jane",
				LastName:          "Doe",
				Gender:            "Female",
				PhoneNumber:       "",
				Balance:           0,
			}
			resp := runHTTP(t, app, httpInput{
				Method: http.MethodPut, Path: fmt.Sprintf("/users/%d", userID),
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
func TestUpdateUser_NotFound(t *testing.T) {
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

			table := "users"
			userID := uint(42)

			ExpAuthUser(userID, false, false, false)(mock)
			ExpSelectByIDEmpty(table, userID)(mock)

			app := setupApp(gdb)
			resp := runHTTP(t, app, httpInput{
				Method: http.MethodPut, Path: fmt.Sprintf("/users/%d", userID), UserID: &userID,
			})
			wantStatus(t, resp, http.StatusNotFound)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("unmet expectations: %v", err)
			}
		})
	}
}

// 500
func TestUpdateUser_DBError(t *testing.T) {
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

			table := "users"
			userID := uint(42)
			studentID := "6600000000"

			ExpAuthUser(userID, false, false, false)(mock)
			ExpSelectByIDFound(table, userID, []string{"id"}, []any{userID})(mock)
			ExpPreloadCanEmpty("learners", []string{"id", "user_id"})(mock)
			ExpPreloadCanEmpty("teachers", []string{"id", "user_id"})(mock)
			ExpPreloadCanEmpty("admins", []string{"id", "user_id"})(mock)
			ExpUpdateError(table, fmt.Errorf("update failed"))(mock)

			app := setupApp(gdb)
			payload := models.User{
				StudentID:         studentID,
				ProfilePictureURL: "",
				FirstName:         "Jane",
				LastName:          "Doe",
				Gender:            "Female",
				PhoneNumber:       "",
				Balance:           0,
			}
			resp := runHTTP(t, app, httpInput{
				Method: http.MethodPut, Path: fmt.Sprintf("/users/%d", userID),
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
func TestUpdateUser_BadRequest(t *testing.T) {
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
				Method: http.MethodPut, Path: "/users/not-an-int", UserID: &userID,
			})
			wantStatus(t, resp, http.StatusBadRequest)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("unmet expectations: %v", err)
			}
		})
	}
}

/* ------------------ DeleteUser ------------------ */

// 200
func TestDeleteUser_OK_SoftDelete(t *testing.T) {
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

			table := "users"
			userID := uint(42)

			ExpAuthUser(userID, false, false, false)(mock)
			ExpSelectByIDFound(table, userID, []string{"id"}, []any{userID})(mock)
			ExpPreloadCanEmpty("learners", []string{"id", "user_id"})(mock)
			ExpPreloadCanEmpty("teachers", []string{"id", "user_id"})(mock)
			ExpPreloadCanEmpty("admins", []string{"id", "user_id"})(mock)
			ExpSoftDeleteOK(table)(mock)
			ExpSoftDeleteOKWithAllowNoTransaction("learners")(mock)
			ExpSoftDeleteOKWithAllowNoTransaction("teachers")(mock)
			ExpSoftDeleteOKWithAllowNoTransaction("admins")(mock)

			app := setupApp(gdb)
			resp := runHTTP(t, app, httpInput{
				Method: http.MethodDelete, Path: fmt.Sprintf("/users/%d", userID), UserID: &userID,
			})
			wantStatus(t, resp, http.StatusOK)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("unmet expectations: %v", err)
			}
		})
	}
}

// 404
func TestDeleteUser_NotFound(t *testing.T) {
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

			table := "users"
			userID := uint(42)

			ExpAuthUser(userID, false, false, false)(mock)
			ExpSelectByIDEmpty(table, userID)(mock)

			app := setupApp(gdb)
			resp := runHTTP(t, app, httpInput{
				Method: http.MethodDelete, Path: fmt.Sprintf("/users/%d", userID), UserID: &userID,
			})
			wantStatus(t, resp, http.StatusNotFound)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("unmet expectations: %v", err)
			}
		})
	}
}

// 500
func TestDeleteUser_DBError(t *testing.T) {
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

			table := "users"
			userID := uint(42)

			ExpAuthUser(userID, false, false, false)(mock)
			ExpSelectByIDFound(table, userID, []string{"id"}, []any{userID})(mock)
			ExpPreloadCanEmpty("learners", []string{"id", "user_id"})(mock)
			ExpPreloadCanEmpty("teachers", []string{"id", "user_id"})(mock)
			ExpPreloadCanEmpty("admins", []string{"id", "user_id"})(mock)
			ExpSoftDeleteError(table, fmt.Errorf("update failed"))(mock)
			ExpSoftDeleteOKWithAllowNoTransaction("learners")(mock)
			ExpSoftDeleteOKWithAllowNoTransaction("teachers")(mock)
			ExpSoftDeleteOKWithAllowNoTransaction("admins")(mock)

			app := setupApp(gdb)
			resp := runHTTP(t, app, httpInput{
				Method: http.MethodDelete, Path: fmt.Sprintf("/users/%d", userID), UserID: &userID,
			})
			wantStatus(t, resp, http.StatusInternalServerError)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("unmet expectations: %v", err)
			}
		})
	}
}

// 400
func TestDeleteUser_BadRequest(t *testing.T) {
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
				Method: http.MethodDelete, Path: "/users/not-an-int", UserID: &userID,
			})
			wantStatus(t, resp, http.StatusBadRequest)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("unmet expectations: %v", err)
			}
		})
	}
}

/* ------------------ processProfilePicture ------------------ */

func TestProcessProfilePicture_Success(t *testing.T) {
	fu := &fakeUploader{}
	app, ctx := newCtxWithUploader(fu)
	defer app.ReleaseCtx(ctx)

	u := &models.User{ProfilePictureURL: tinyPNGRawBase64}

	err := processProfilePicture(ctx, u)
	require.NoError(t, err)

	require.Equal(t, "users", fu.lastBucket)
	require.NotEmpty(t, fu.lastFilename)
	require.NotEmpty(t, fu.lastData)

	require.Regexp(t, regexp.MustCompile(`^users/\d+\.(png|jpg|jpeg|gif|webp|bin)$`), u.ProfilePictureURL)
}

func TestProcessProfilePicture_SkipEmptyAndHTTP(t *testing.T) {
	fu := &fakeUploader{}
	app, ctx := newCtxWithUploader(fu)
	defer app.ReleaseCtx(ctx)

	u1 := &models.User{ProfilePictureURL: ""}
	require.NoError(t, processProfilePicture(ctx, u1))
	require.Equal(t, "", u1.ProfilePictureURL)
	require.Equal(t, "", fu.lastFilename)

	u2 := &models.User{ProfilePictureURL: "http://example.com/pic.png"}
	require.NoError(t, processProfilePicture(ctx, u2))
	require.Equal(t, "http://example.com/pic.png", u2.ProfilePictureURL)
	require.Equal(t, "", fu.lastFilename)
}

func TestProcessProfilePicture_InvalidBase64(t *testing.T) {
	fu := &fakeUploader{}
	app, ctx := newCtxWithUploader(fu)
	defer app.ReleaseCtx(ctx)

	u := &models.User{ProfilePictureURL: "!!!not-b64!!!"}
	err := processProfilePicture(ctx, u)
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid base64 image")
}

func TestProcessProfilePicture_InvalidImage(t *testing.T) {
	fu := &fakeUploader{}
	app, ctx := newCtxWithUploader(fu)
	defer app.ReleaseCtx(ctx)

	s := base64.StdEncoding.EncodeToString([]byte("hello"))
	u := &models.User{ProfilePictureURL: s}
	err := processProfilePicture(ctx, u)
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid image")
}
