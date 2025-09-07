package handlers

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/a2n2k3p4/tutorium-backend/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

/* ------------------ CreateAdmin ------------------ */

// 201
func TestCreateAdmin_OK(t *testing.T) {
	table := "admins"
	userID := uint(42)
	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			ExpInsertReturningID(table, 1)(mock)

			req := jsonBody(models.Admin{
				UserID: userID,
			},
			)

			*payload = req
			*uID = userID

		},
		http.StatusCreated,
		http.MethodPost,
		"/admins/",
	)
}

// 400
func TestCreateAdmin_BadRequest(t *testing.T) {
	userID := uint(42)
	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {

			ExpAuthUser(userID, false, false, false)(mock)

			*uID = userID
		},
		http.StatusBadRequest,
		http.MethodPost,
		"/admins/",
	)
}

// 500
func TestCreateAdmin_DBError(t *testing.T) {
	table := "admins"
	userID := uint(42)
	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			ExpInsertError(table, fmt.Errorf("db insert failed"))(mock)

			req := jsonBody(models.Admin{
				UserID: userID,
			},
			)

			*payload = req
			*uID = userID
		},
		http.StatusInternalServerError,
		http.MethodPost,
		"/admins/",
	)
}

/* ------------------ GetAdmins ------------------ */
// 200
func TestGetAdmins_OK(t *testing.T) {
	table := "admins"
	userID := uint(42)
	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			ExpListRows(table, []string{"id"}, []any{1}, []any{2})(mock)

			*uID = userID
		},
		http.StatusOK,
		http.MethodGet,
		"/admins/",
	)
}

// 500
func TestGetAdmins_DBError(t *testing.T) {
	table := "admins"
	userID := uint(42)
	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			ExpListError(table, fmt.Errorf("select failed"))(mock)

			*uID = userID
		},
		http.StatusInternalServerError,
		http.MethodGet,
		"/admins/",
	)
}

/* ------------------ GetAdmin ------------------ */
// 200
func TestGetAdmin_OK(t *testing.T) {
	table := "admins"
	userID := uint(42)
	adminID := uint(7)
	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			ExpSelectByIDFound(table, adminID, []string{"id"}, []any{adminID})(mock)

			*uID = userID
		},
		http.StatusOK,
		http.MethodGet,
		fmt.Sprintf("/admins/%d", adminID),
	)
}

// 404
func TestGetAdmin_NotFound(t *testing.T) {
	table := "admins"
	userID := uint(42)
	adminID := uint(999)
	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			ExpSelectByIDEmpty(table, adminID)(mock)

			*uID = userID
		},
		http.StatusNotFound,
		http.MethodGet,
		fmt.Sprintf("/admins/%d", adminID),
	)
}

// 500
func TestGetAdmin_DBError(t *testing.T) {
	table := "admins"
	userID := uint(42)
	adminID := uint(7)
	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			ExpSelectByIDError(table, adminID, fmt.Errorf("select failed"))(mock)

			*uID = userID
		},
		http.StatusInternalServerError,
		http.MethodGet,
		fmt.Sprintf("/admins/%d", adminID),
	)
}

// 400
func TestGetAdmin_BadRequest(t *testing.T) {
	userID := uint(42)
	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			*uID = userID
		},
		http.StatusBadRequest,
		http.MethodGet,
		"/admins/not-an-int",
	)
}

/* ------------------ DeleteAdmin ------------------ */

// 200
func TestDeleteAdmin_OK_SoftDelete(t *testing.T) {
	table := "admins"
	userID := uint(42)
	adminID := uint(5)
	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			ExpSelectByIDFound(table, adminID, []string{"id"}, []any{adminID})(mock)
			ExpSoftDeleteOK(table)(mock)
			*uID = userID
		},
		http.StatusOK,
		http.MethodDelete,
		fmt.Sprintf("/admins/%d", adminID),
	)
}

// 404
func TestDeleteAdmin_NotFound(t *testing.T) {
	table := "admins"
	userID := uint(42)
	adminID := uint(12345)
	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			ExpSelectByIDEmpty(table, adminID)(mock)
			*uID = userID
		},
		http.StatusNotFound,
		http.MethodDelete,
		fmt.Sprintf("/admins/%d", adminID),
	)
}

// 500
func TestDeleteAdmin_DBError(t *testing.T) {
	table := "admins"
	userID := uint(42)
	adminID := uint(5)
	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			ExpSelectByIDFound(table, adminID, []string{"id"}, []any{adminID})(mock)
			ExpSoftDeleteError(table, fmt.Errorf("update failed"))(mock)

		},
		http.StatusInternalServerError,
		http.MethodDelete,
		fmt.Sprintf("/admins/%d", adminID),
	)
}

// 400
func TestDeleteAdmin_BadRequest(t *testing.T) {
	userID := uint(42)
	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {

			ExpAuthUser(userID, false, false, false)(mock)

		},
		http.StatusBadRequest,
		http.MethodDelete,
		"/admins/not-an-int",
	)
}
