package handlers

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/a2n2k3p4/tutorium-backend/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

/* ------------------ CreateBanTeacher ------------------ */

// 201
func TestCreateBanTeacher_OK(t *testing.T) {
	table := "ban_details_teachers"
	userID := uint(42)
	teacherID := uint(50)
	now := time.Now()

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, true, false, false)(mock)
			ExpInsertReturningID(table, 1)(mock)

			req := jsonBody(models.BanDetailsTeacher{
				TeacherID:      teacherID,
				BanStart:       now,
				BanEnd:         now.Add(2 * time.Hour),
				BanDescription: "spamming",
			})

			*payload = req
			*uID = userID
		},
		http.StatusCreated,
		http.MethodPost,
		"/banteachers/",
	)
}

// 400
func TestCreateBanTeacher_BadRequest(t *testing.T) {
	userID := uint(42)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, true, false, false)(mock)
			*uID = userID
		},
		http.StatusBadRequest,
		http.MethodPost,
		"/banteachers/",
	)
}

// 500
func TestCreateBanTeacher_DBError(t *testing.T) {
	table := "ban_details_teachers"
	userID := uint(42)
	teacherID := uint(50)
	now := time.Now()

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, true, false, false)(mock)
			ExpInsertError(table, fmt.Errorf("db insert failed"))(mock)

			req := jsonBody(models.BanDetailsTeacher{
				TeacherID:      teacherID,
				BanStart:       now,
				BanEnd:         now.Add(2 * time.Hour),
				BanDescription: "spamming",
			})

			*payload = req
			*uID = userID
		},
		http.StatusInternalServerError,
		http.MethodPost,
		"/banteachers/",
	)
}

/* ------------------ GetBanTeachers ------------------ */

// 200
func TestGetBanTeachers_OK(t *testing.T) {
	userID := uint(42)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, true, false, false)(mock)
			ExpListRows("ban_details_teachers", []string{"id"}, []any{1}, []any{2})(mock)
			*uID = userID
		},
		http.StatusOK,
		http.MethodGet,
		"/banteachers/",
	)
}

// 500
func TestGetBanTeachers_DBError(t *testing.T) {
	userID := uint(42)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, true, false, false)(mock)
			ExpListError("ban_details_teachers", fmt.Errorf("select failed"))(mock)
			*uID = userID
		},
		http.StatusInternalServerError,
		http.MethodGet,
		"/banteachers/",
	)
}

/* ------------------ GetBanTeacher ------------------ */

// 200
func TestGetBanTeacher_OK(t *testing.T) {
	table := "ban_details_teachers"
	userID := uint(42)
	ban_teacherID := uint(7)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, true, false, false)(mock)
			ExpSelectByIDFound(table, ban_teacherID, []string{"id"}, []any{ban_teacherID})(mock)
			*uID = userID
		},
		http.StatusOK,
		http.MethodGet,
		fmt.Sprintf("/banteachers/%d", ban_teacherID),
	)
}

// 404
func TestGetBanTeacher_NotFound(t *testing.T) {
	table := "ban_details_teachers"
	userID := uint(42)
	ban_teacherID := uint(999)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, true, false, false)(mock)
			ExpSelectByIDEmpty(table, ban_teacherID)(mock)
			*uID = userID
		},
		http.StatusNotFound,
		http.MethodGet,
		fmt.Sprintf("/banteachers/%d", ban_teacherID),
	)
}

// 500
func TestGetBanTeacher_DBError(t *testing.T) {
	table := "ban_details_teachers"
	userID := uint(42)
	ban_teacherID := uint(7)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, true, false, false)(mock)
			ExpSelectByIDError(table, ban_teacherID, fmt.Errorf("select failed"))(mock)
			*uID = userID
		},
		http.StatusInternalServerError,
		http.MethodGet,
		fmt.Sprintf("/banteachers/%d", ban_teacherID),
	)
}

// 400
func TestGetBanTeacher_BadRequest(t *testing.T) {
	userID := uint(42)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, true, false, false)(mock)
			*uID = userID
		},
		http.StatusBadRequest,
		http.MethodGet,
		"/banteachers/not-an-int",
	)
}

/* ------------------ UpdateBanTeacher ------------------ */

// 200
func TestUpdateBanTeacher_OK(t *testing.T) {
	table := "ban_details_teachers"
	preloadTable := "teachers"
	userID := uint(42)
	teacherID := uint(50)
	ban_teacherID := uint(1)
	now := time.Now()

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, true, false, false)(mock)
			ExpSelectByIDFound(table, ban_teacherID,
				[]string{"id", "teacher_id", "ban_start", "ban_end", "ban_description"},
				[]any{ban_teacherID, teacherID, now, now.Add(2 * time.Hour), "flooding"},
			)(mock)

			ExpPreloadField(preloadTable, []string{"id"}, []any{teacherID})(mock)
			ExpUpdateOK(table)(mock)

			req := jsonBody(models.BanDetailsTeacher{
				BanDescription: "spamming",
			})
			*payload = req
			*uID = userID
		},
		http.StatusOK,
		http.MethodPut,
		fmt.Sprintf("/banteachers/%d", ban_teacherID),
	)
}

// 404
func TestUpdateBanTeacher_NotFound(t *testing.T) {
	table := "ban_details_teachers"
	userID := uint(42)
	ban_teacherID := uint(12345)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, true, false, false)(mock)
			ExpSelectByIDEmpty(table, ban_teacherID)(mock)
			*uID = userID
		},
		http.StatusNotFound,
		http.MethodPut,
		fmt.Sprintf("/banteachers/%d", ban_teacherID),
	)
}

// 500
func TestUpdateBanTeacher_DBError(t *testing.T) {
	table := "ban_details_teachers"
	userID := uint(42)
	ban_teacherID := uint(1)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, true, false, false)(mock)
			ExpSelectByIDFound(table, ban_teacherID, []string{"id"}, []any{ban_teacherID})(mock)
			ExpUpdateError(table, fmt.Errorf("update failed"))(mock)

			req := jsonBody(models.BanDetailsTeacher{
				BanDescription: "spamming",
			})
			*payload = req
			*uID = userID
		},
		http.StatusInternalServerError,
		http.MethodPut,
		fmt.Sprintf("/banteachers/%d", ban_teacherID),
	)
}

// 400
func TestUpdateBanTeacher_BadRequest(t *testing.T) {
	userID := uint(42)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, true, false, false)(mock)
			*uID = userID
		},
		http.StatusBadRequest,
		http.MethodPut,
		"/banteachers/not-an-int",
	)
}

/* ------------------ DeleteBanTeacher ------------------ */

// 200
func TestDeleteBanTeacher_OK_SoftDelete(t *testing.T) {
	table := "ban_details_teachers"
	userID := uint(42)
	ban_teacherID := uint(5)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, true, false, false)(mock)
			ExpSelectByIDFound(table, ban_teacherID, []string{"id"}, []any{ban_teacherID})(mock)
			ExpSoftDeleteOK(table)(mock)
			*uID = userID
		},
		http.StatusOK,
		http.MethodDelete,
		fmt.Sprintf("/banteachers/%d", ban_teacherID),
	)
}

// 404
func TestDeleteBanTeacher_NotFound(t *testing.T) {
	table := "ban_details_teachers"
	userID := uint(42)
	ban_teacherID := uint(12345)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, true, false, false)(mock)
			ExpSelectByIDEmpty(table, ban_teacherID)(mock)
			*uID = userID
		},
		http.StatusNotFound,
		http.MethodDelete,
		fmt.Sprintf("/banteachers/%d", ban_teacherID),
	)
}

// 500
func TestDeleteBanTeacher_DBError(t *testing.T) {
	table := "ban_details_teachers"
	userID := uint(42)
	ban_teacherID := uint(5)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, true, false, false)(mock)
			ExpSelectByIDFound(table, ban_teacherID, []string{"id"}, []any{ban_teacherID})(mock)
			ExpSoftDeleteError(table, fmt.Errorf("update failed"))(mock)
			*uID = userID
		},
		http.StatusInternalServerError,
		http.MethodDelete,
		fmt.Sprintf("/banteachers/%d", ban_teacherID),
	)
}

// 400
func TestDeleteBanTeacher_BadRequest(t *testing.T) {
	userID := uint(42)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, true, false, false)(mock)
			*uID = userID
		},
		http.StatusBadRequest,
		http.MethodDelete,
		"/banteachers/not-an-int",
	)
}
