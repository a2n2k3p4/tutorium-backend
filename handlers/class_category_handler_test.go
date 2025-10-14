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

/* ------------------ CreateClassCategory ------------------ */

// 201
func TestCreateClassCategory_OK(t *testing.T) {
	table := "class_categories"
	userID := uint(42)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, true, false, false)(mock)
			ExpInsertReturningID(table, 1)(mock)

			req := jsonBody(models.ClassCategory{
				ClassCategory: "test",
			})
			*payload = req
			*uID = userID
		},
		http.StatusCreated,
		http.MethodPost,
		"/class_categories/",
	)
}

// 400
func TestCreateClassCategory_BadRequest(t *testing.T) {
	userID := uint(42)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, true, false, false)(mock)
			*uID = userID
		},
		http.StatusBadRequest,
		http.MethodPost,
		"/class_categories/",
	)
}

// 500
func TestCreateClassCategory_DBError(t *testing.T) {
	table := "class_categories"
	userID := uint(42)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, true, false, false)(mock)
			ExpInsertError(table, fmt.Errorf("db insert failed"))(mock)
			req := jsonBody(models.ClassCategory{
				ClassCategory: "test",
			})
			*payload = req
			*uID = userID
		},
		http.StatusInternalServerError,
		http.MethodPost,
		"/class_categories/",
	)
}

/* ------------------ GetClassCategories ------------------ */

// 200
func TestGetClassCategories_OK(t *testing.T) {
	userID := uint(42)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, true, false, false)(mock)
			ExpListRows("class_categories", []string{"id"}, []any{1}, []any{2})(mock)
			*uID = userID
		},
		http.StatusOK,
		http.MethodGet,
		"/class_categories/",
	)
}

// 500
func TestGetClassCategories_DBError(t *testing.T) {
	userID := uint(42)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, true, false, false)(mock)
			ExpListError("class_categories", fmt.Errorf("select failed"))(mock)
			*uID = userID
		},
		http.StatusInternalServerError,
		http.MethodGet,
		"/class_categories/",
	)
}

/* ------------------ GetClassCategory ------------------ */

// 200
func TestGetClassCategory_OK(t *testing.T) {
	table := "class_categories"
	userID := uint(42)
	classCategoryID := uint(7)
	classID := uint(5)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, true, false, false)(mock)
			ExpSelectByIDFound(table, classCategoryID, []string{"id"}, []any{classCategoryID})(mock)
			ExpPreloadM2M("class_class_categories", "classes", "class_category_id", "class_id", [][2]any{{classCategoryID, classID}})(mock)
			*uID = userID
		},
		http.StatusOK,
		http.MethodGet,
		fmt.Sprintf("/class_categories/%d", classCategoryID),
	)
}

// 404
func TestGetClassCategory_NotFound(t *testing.T) {
	table := "class_categories"
	userID := uint(42)
	classCategoryID := uint(999)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, true, false, false)(mock)
			ExpSelectByIDEmpty(table, classCategoryID)(mock)
			*uID = userID
		},
		http.StatusNotFound,
		http.MethodGet,
		fmt.Sprintf("/class_categories/%d", classCategoryID),
	)
}

// 500
func TestGetClassCategory_DBError(t *testing.T) {
	table := "class_categories"
	userID := uint(42)
	classCategoryID := uint(7)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, true, false, false)(mock)
			ExpSelectByIDError(table, classCategoryID, fmt.Errorf("select failed"))(mock)
			*uID = userID
		},
		http.StatusInternalServerError,
		http.MethodGet,
		fmt.Sprintf("/class_categories/%d", classCategoryID),
	)
}

// 400
func TestGetClassCategory_BadRequest(t *testing.T) {
	userID := uint(42)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, true, false, false)(mock)
			*uID = userID
		},
		http.StatusBadRequest,
		http.MethodGet,
		"/class_categories/not-an-int",
	)
}

/* ------------------ UpdateClassCategory ------------------ */

// 200
func TestUpdateClassCategory_OK(t *testing.T) {
	table := "class_categories"
	userID := uint(42)
	classCategoryID := uint(1)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, true, false, false)(mock)
			ExpSelectByIDFound(table, classCategoryID, []string{"id", "class_category"}, []any{classCategoryID, "test"})(mock)
			ExpUpdateOK(table)(mock)

			req := jsonBody(models.ClassCategory{
				ClassCategory: "edit test",
			})
			*payload = req
			*uID = userID
		},
		http.StatusOK,
		http.MethodPut,
		fmt.Sprintf("/class_categories/%d", classCategoryID),
	)
}

// 404
func TestUpdateClassCategory_NotFound(t *testing.T) {
	table := "class_categories"
	userID := uint(42)
	classCategoryID := uint(12345)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, true, false, false)(mock)
			ExpSelectByIDEmpty(table, classCategoryID)(mock)
			*uID = userID
		},
		http.StatusNotFound,
		http.MethodPut,
		fmt.Sprintf("/class_categories/%d", classCategoryID),
	)
}

// 500
func TestUpdateClassCategory_DBError(t *testing.T) {
	table := "class_categories"
	userID := uint(42)
	classCategoryID := uint(1)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, true, false, false)(mock)
			ExpSelectByIDFound(table, classCategoryID, []string{"id", "class_category"}, []any{classCategoryID, "test"})(mock)
			ExpUpdateError(table, fmt.Errorf("update failed"))(mock)

			req := jsonBody(models.ClassCategory{
				ClassCategory: "edit test",
			})
			*payload = req
			*uID = userID
		},
		http.StatusInternalServerError,
		http.MethodPut,
		fmt.Sprintf("/class_categories/%d", classCategoryID),
	)
}

// 400
func TestUpdateClassCategory_BadRequest(t *testing.T) {
	userID := uint(42)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, true, false, false)(mock)
			*uID = userID
		},
		http.StatusBadRequest,
		http.MethodPut,
		"/class_categories/not-an-int",
	)
}

/* ------------------ DeleteClassCategory ------------------ */

// 200 (soft delete)
func TestDeleteClassCategory_OK_SoftDelete(t *testing.T) {
	table := "class_categories"
	userID := uint(42)
	classCategoryID := uint(5)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, true, false, false)(mock)
			ExpSelectByIDFound(table, classCategoryID, []string{"id"}, []any{classCategoryID})(mock)
			ExpClearClassesForCategory(classCategoryID)(mock)
			ExpClearLearnersForCategory(classCategoryID)(mock)
			ExpSoftDeleteOK(table)(mock)
			*uID = userID
		},
		http.StatusOK,
		http.MethodDelete,
		fmt.Sprintf("/class_categories/%d", classCategoryID),
	)
}

// 404
func TestDeleteClassCategory_NotFound(t *testing.T) {
	table := "class_categories"
	userID := uint(42)
	classCategoryID := uint(12345)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, true, false, false)(mock)
			ExpSelectByIDEmpty(table, classCategoryID)(mock)
			*uID = userID
		},
		http.StatusNotFound,
		http.MethodDelete,
		fmt.Sprintf("/class_categories/%d", classCategoryID),
	)
}

// 500
func TestDeleteClassCategory_DBError(t *testing.T) {
	table := "class_categories"
	userID := uint(42)
	classCategoryID := uint(5)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, true, false, false)(mock)
			ExpSelectByIDFound(table, classCategoryID, []string{"id"}, []any{classCategoryID})(mock)
			ExpClearClassesForCategory(classCategoryID)(mock)
			ExpClearLearnersForCategory(classCategoryID)(mock)
			ExpSoftDeleteError(table, fmt.Errorf("update failed"))(mock)
			*uID = userID
		},
		http.StatusInternalServerError,
		http.MethodDelete,
		fmt.Sprintf("/class_categories/%d", classCategoryID),
	)
}

// 400
func TestDeleteClassCategory_BadRequest(t *testing.T) {
	userID := uint(42)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, true, false, false)(mock)
			*uID = userID
		},
		http.StatusBadRequest,
		http.MethodDelete,
		"/class_categories/not-an-int",
	)
}
