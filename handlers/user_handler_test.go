package handlers

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/a2n2k3p4/tutorium-backend/models"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

/* ------------------ CreateUser ------------------ */

// 201
func TestCreateUser_OK(t *testing.T) {
	table := "users"
	userID := uint(42)
	studentID := "6600000000"

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			ExpDoubleInsertReturningID(table, "learners", uint64(userID), 2)(mock)

			req := jsonBody(models.User{
				StudentID:         studentID,
				ProfilePictureURL: "",
				FirstName:         "Jane",
				LastName:          "Doe",
				Gender:            "Female",
				PhoneNumber:       "",
				Balance:           0,
			})
			*payload = req
			*uID = userID
		},
		http.StatusCreated,
		http.MethodPost,
		"/users/",
	)
}

// 400
func TestCreateUser_BadRequest(t *testing.T) {
	userID := uint(42)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			*uID = userID
		},
		http.StatusBadRequest,
		http.MethodPost,
		"/users/",
	)
}

// 500
func TestCreateUser_DBError(t *testing.T) {
	table := "users"
	userID := uint(42)
	studentID := "6600000000"

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			ExpInsertError(table, fmt.Errorf("db insert failed"))(mock)

			req := jsonBody(models.User{
				StudentID:         studentID,
				ProfilePictureURL: "",
				FirstName:         "Jane",
				LastName:          "Doe",
				Gender:            "Female",
				PhoneNumber:       "",
				Balance:           0,
			})
			*payload = req
			*uID = userID
		},
		http.StatusInternalServerError,
		http.MethodPost,
		"/users/",
	)
}

/* ------------------ GetUsers ------------------ */

// 200
func TestGetUsers_OK(t *testing.T) {
	userID := uint(42)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, true, false, false)(mock)
			ExpListRows("users", []string{"id"}, []any{1}, []any{2})(mock)
			ExpPreloadCanEmpty("learners", []string{"id", "user_id"})(mock)
			ExpPreloadCanEmpty("teachers", []string{"id", "user_id"})(mock)
			ExpPreloadCanEmpty("admins", []string{"id", "user_id"})(mock)
			*uID = userID
		},
		http.StatusOK,
		http.MethodGet,
		"/users/",
	)
}

// 500
func TestGetUsers_DBError(t *testing.T) {
	userID := uint(42)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, true, false, false)(mock)
			ExpListError("users", fmt.Errorf("select failed"))(mock)
			*uID = userID
		},
		http.StatusInternalServerError,
		http.MethodGet,
		"/users/",
	)
}

/* ------------------ GetUser ------------------ */

// 200
func TestGetUser_OK(t *testing.T) {
	table := "users"
	userID := uint(42)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			ExpSelectByIDFound(table, userID, []string{"id"}, []any{userID})(mock)
			ExpPreloadCanEmpty("learners", []string{"id", "user_id"})(mock)
			ExpPreloadCanEmpty("teachers", []string{"id", "user_id"})(mock)
			ExpPreloadCanEmpty("admins", []string{"id", "user_id"})(mock)
			*uID = userID
		},
		http.StatusOK,
		http.MethodGet,
		fmt.Sprintf("/users/%d", userID),
	)
}

// 404
func TestGetUser_NotFound(t *testing.T) {
	table := "users"
	userID := uint(42)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			ExpSelectByIDEmpty(table, userID)(mock)
			*uID = userID
		},
		http.StatusNotFound,
		http.MethodGet,
		fmt.Sprintf("/users/%d", userID),
	)
}

// 500
func TestGetUser_DBError(t *testing.T) {
	table := "users"
	userID := uint(42)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			ExpSelectByIDError(table, userID, fmt.Errorf("select failed"))(mock)
			*uID = userID
		},
		http.StatusInternalServerError,
		http.MethodGet,
		fmt.Sprintf("/users/%d", userID),
	)
}

// 400
func TestGetUser_BadRequest(t *testing.T) {
	userID := uint(42)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			*uID = userID
		},
		http.StatusBadRequest,
		http.MethodGet,
		"/users/not-an-int",
	)
}

/* ------------------ UpdateUser ------------------ */

// 200
func TestUpdateUser_OK(t *testing.T) {
	table := "users"
	userID := uint(42)
	studentID := "6600000000"

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			ExpSelectByIDFound(table, userID,
				[]string{"id", "student_id", "profilePicture_url", "first_name", "last_name", "gender", "phone_number", "balance"},
				[]any{userID, studentID, "", "Janet", "Doe", "Female", "", 50},
			)(mock)

			ExpPreloadCanEmpty("learners", []string{"id", "user_id"})(mock)
			ExpPreloadCanEmpty("teachers", []string{"id", "user_id"})(mock)
			ExpPreloadCanEmpty("admins", []string{"id", "user_id"})(mock)

			ExpUpdateOK(table)(mock)

			req := jsonBody(models.User{
				StudentID:         studentID,
				ProfilePictureURL: "",
				FirstName:         "Jane",
				LastName:          "Doe",
				Gender:            "Female",
				PhoneNumber:       "",
				Balance:           0,
			})
			*payload = req
			*uID = userID
		},
		http.StatusOK,
		http.MethodPut,
		fmt.Sprintf("/users/%d", userID),
	)
}

// 404
func TestUpdateUser_NotFound(t *testing.T) {
	table := "users"
	userID := uint(42)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			ExpSelectByIDEmpty(table, userID)(mock)
			*uID = userID
		},
		http.StatusNotFound,
		http.MethodPut,
		fmt.Sprintf("/users/%d", userID),
	)
}

// 500
func TestUpdateUser_DBError(t *testing.T) {
	table := "users"
	userID := uint(42)
	studentID := "6600000000"

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			ExpSelectByIDFound(table, userID, []string{"id"}, []any{userID})(mock)
			ExpPreloadCanEmpty("learners", []string{"id", "user_id"})(mock)
			ExpPreloadCanEmpty("teachers", []string{"id", "user_id"})(mock)
			ExpPreloadCanEmpty("admins", []string{"id", "user_id"})(mock)
			ExpUpdateError(table, fmt.Errorf("update failed"))(mock)

			req := jsonBody(models.User{
				StudentID:         studentID,
				ProfilePictureURL: "",
				FirstName:         "Jane",
				LastName:          "Doe",
				Gender:            "Female",
				PhoneNumber:       "",
				Balance:           0,
			})
			*payload = req
			*uID = userID
		},
		http.StatusInternalServerError,
		http.MethodPut,
		fmt.Sprintf("/users/%d", userID),
	)
}

// 400
func TestUpdateUser_BadRequest(t *testing.T) {
	userID := uint(42)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			*uID = userID
		},
		http.StatusBadRequest,
		http.MethodPut,
		"/users/not-an-int",
	)
}

/* ------------------ DeleteUser ------------------ */

// 200
func TestDeleteUser_OK_SoftDelete(t *testing.T) {
	table := "users"
	userID := uint(42)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			ExpSelectByIDFound(table, userID, []string{"id"}, []any{userID})(mock)
			ExpPreloadCanEmpty("learners", []string{"id", "user_id"})(mock)
			ExpPreloadCanEmpty("teachers", []string{"id", "user_id"})(mock)
			ExpPreloadCanEmpty("admins", []string{"id", "user_id"})(mock)
			ExpSoftDeleteOK(table)(mock)
			ExpSoftDeleteOKWithAllowNoTransaction("learners")(mock)
			ExpSoftDeleteOKWithAllowNoTransaction("teachers")(mock)
			ExpSoftDeleteOKWithAllowNoTransaction("admins")(mock)
			*uID = userID
		},
		http.StatusOK,
		http.MethodDelete,
		fmt.Sprintf("/users/%d", userID),
	)
}

// 404
func TestDeleteUser_NotFound(t *testing.T) {
	table := "users"
	userID := uint(42)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			ExpSelectByIDEmpty(table, userID)(mock)
			*uID = userID
		},
		http.StatusNotFound,
		http.MethodDelete,
		fmt.Sprintf("/users/%d", userID),
	)
}

// 500
func TestDeleteUser_DBError(t *testing.T) {
	table := "users"
	userID := uint(42)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			ExpSelectByIDFound(table, userID, []string{"id"}, []any{userID})(mock)
			ExpPreloadCanEmpty("learners", []string{"id", "user_id"})(mock)
			ExpPreloadCanEmpty("teachers", []string{"id", "user_id"})(mock)
			ExpPreloadCanEmpty("admins", []string{"id", "user_id"})(mock)
			ExpSoftDeleteError(table, fmt.Errorf("update failed"))(mock)
			ExpSoftDeleteOKWithAllowNoTransaction("learners")(mock)
			ExpSoftDeleteOKWithAllowNoTransaction("teachers")(mock)
			ExpSoftDeleteOKWithAllowNoTransaction("admins")(mock)
			*uID = userID
		},
		http.StatusInternalServerError,
		http.MethodDelete,
		fmt.Sprintf("/users/%d", userID),
	)
}

// 400
func TestDeleteUser_BadRequest(t *testing.T) {
	userID := uint(42)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			*uID = userID
		},
		http.StatusBadRequest,
		http.MethodDelete,
		"/users/not-an-int",
	)
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
