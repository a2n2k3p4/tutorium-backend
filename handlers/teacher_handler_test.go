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

/* ------------------ CreateTeacher ------------------ */

// 201
func TestCreateTeacher_OK(t *testing.T) {
	table := "teachers"
	userID := uint(42)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			ExpInsertReturningID(table, 1)(mock)

			req := jsonBody(models.Teacher{
				UserID: userID,
			})
			*payload = req
			*uID = userID
		},
		http.StatusCreated,
		http.MethodPost,
		"/teachers/",
	)
}

// 400
func TestCreateTeacher_BadRequest(t *testing.T) {
	userID := uint(42)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			*uID = userID
		},
		http.StatusBadRequest,
		http.MethodPost,
		"/teachers/",
	)
}

// 500
func TestCreateTeacher_DBError(t *testing.T) {
	table := "teachers"
	userID := uint(42)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			ExpInsertError(table, fmt.Errorf("db insert failed"))(mock)

			req := jsonBody(models.Teacher{
				UserID: userID,
			})
			*payload = req
			*uID = userID
		},
		http.StatusInternalServerError,
		http.MethodPost,
		"/teachers/",
	)
}

/* ------------------ GetTeachers ------------------ */

// 200
func TestGetTeachers_OK(t *testing.T) {
	userID := uint(42)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			ExpListRows("teachers", []string{"id"}, []any{1}, []any{2})(mock)
			*uID = userID
		},
		http.StatusOK,
		http.MethodGet,
		"/teachers/",
	)
}

// 500
func TestGetTeachers_DBError(t *testing.T) {
	userID := uint(42)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			ExpListError("teachers", fmt.Errorf("select failed"))(mock)
			*uID = userID
		},
		http.StatusInternalServerError,
		http.MethodGet,
		"/teachers/",
	)
}

/* ------------------ GetTeacher ------------------ */

// 200
func TestGetTeacher_OK(t *testing.T) {
	table := "teachers"
	userID := uint(42)
	teacherID := uint(7)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			ExpSelectByIDFound(table, teacherID, []string{"id"}, []any{teacherID})(mock)
			*uID = userID
		},
		http.StatusOK,
		http.MethodGet,
		fmt.Sprintf("/teachers/%d", teacherID),
	)
}

// 404
func TestGetTeacher_NotFound(t *testing.T) {
	table := "teachers"
	userID := uint(42)
	teacherID := uint(999)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			ExpSelectByIDEmpty(table, teacherID)(mock)
			*uID = userID
		},
		http.StatusNotFound,
		http.MethodGet,
		fmt.Sprintf("/teachers/%d", teacherID),
	)
}

// 500
func TestGetTeacher_DBError(t *testing.T) {
	table := "teachers"
	userID := uint(42)
	teacherID := uint(7)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			ExpSelectByIDError(table, teacherID, fmt.Errorf("select failed"))(mock)
			*uID = userID
		},
		http.StatusInternalServerError,
		http.MethodGet,
		fmt.Sprintf("/teachers/%d", teacherID),
	)
}

// 400
func TestGetTeacher_BadRequest(t *testing.T) {
	userID := uint(42)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			*uID = userID
		},
		http.StatusBadRequest,
		http.MethodGet,
		"/teachers/not-an-int",
	)
}

/* ------------------ DeleteTeacher ------------------ */

// 200
func TestDeleteTeacher_OK_SoftDelete(t *testing.T) {
	table := "teachers"
	userID := uint(42)
	teacherID := uint(5)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			ExpSelectByIDFound(table, teacherID, []string{"id"}, []any{teacherID})(mock)
			ExpSoftDeleteOK(table)(mock)
			*uID = userID
		},
		http.StatusOK,
		http.MethodDelete,
		fmt.Sprintf("/teachers/%d", teacherID),
	)
}

// 404
func TestDeleteTeacher_NotFound(t *testing.T) {
	table := "teachers"
	userID := uint(42)
	teacherID := uint(12345)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			ExpSelectByIDEmpty(table, teacherID)(mock)
			*uID = userID
		},
		http.StatusNotFound,
		http.MethodDelete,
		fmt.Sprintf("/teachers/%d", teacherID),
	)
}

// 500
func TestDeleteTeacher_DBError(t *testing.T) {
	table := "teachers"
	userID := uint(42)
	teacherID := uint(5)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			ExpSelectByIDFound(table, teacherID, []string{"id"}, []any{teacherID})(mock)
			ExpSoftDeleteError(table, fmt.Errorf("update failed"))(mock)
			*uID = userID
		},
		http.StatusInternalServerError,
		http.MethodDelete,
		fmt.Sprintf("/teachers/%d", teacherID),
	)
}

// 400
func TestDeleteTeacher_BadRequest(t *testing.T) {
	userID := uint(42)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			*uID = userID
		},
		http.StatusBadRequest,
		http.MethodDelete,
		"/teachers/not-an-int",
	)
}
