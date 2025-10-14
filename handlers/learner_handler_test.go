package handlers

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/a2n2k3p4/tutorium-backend/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

/* ------------------ CreateLearner ------------------ */

// 201
func TestCreateLearner_OK(t *testing.T) {
	table := "learners"
	userID := uint(42)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			ExpInsertReturningID(table, 1)(mock)

			req := jsonBody(models.Learner{
				UserID: userID,
			})
			*payload = req
			*uID = userID
		},
		http.StatusCreated,
		http.MethodPost,
		"/learners/",
	)
}

// 400
func TestCreateLearner_BadRequest(t *testing.T) {
	userID := uint(42)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			*uID = userID
		},
		http.StatusBadRequest,
		http.MethodPost,
		"/learners/",
	)
}

// 500
func TestCreateLearner_DBError(t *testing.T) {
	table := "learners"
	userID := uint(42)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			ExpInsertError(table, fmt.Errorf("db insert failed"))(mock)

			req := jsonBody(models.Learner{
				UserID: userID,
			})
			*payload = req
			*uID = userID
		},
		http.StatusInternalServerError,
		http.MethodPost,
		"/learners/",
	)
}

/* ------------------ GetLearners ------------------ */

// 200
func TestGetLearners_OK(t *testing.T) {
	userID := uint(42)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			ExpListRows("learners", []string{"id"}, []any{1}, []any{2})(mock)
			ExpPreloadLearnersInterestedEmpty(1, 2)(mock)
			*uID = userID
		},
		http.StatusOK,
		http.MethodGet,
		"/learners/",
	)
}

// 500
func TestGetLearners_DBError(t *testing.T) {
	userID := uint(42)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			ExpListError("learners", fmt.Errorf("select failed"))(mock)
			*uID = userID
		},
		http.StatusInternalServerError,
		http.MethodGet,
		"/learners/",
	)
}

/* ------------------ GetLearner ------------------ */

// 200
func TestGetLearner_OK(t *testing.T) {
	table := "learners"
	userID := uint(42)
	learnerID := uint(7)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			ExpSelectByIDFound(table, learnerID, []string{"id"}, []any{learnerID})(mock)
			ExpPreloadLearnerInterestedEmpty(learnerID)(mock)
			*uID = userID
		},
		http.StatusOK,
		http.MethodGet,
		fmt.Sprintf("/learners/%d", learnerID),
	)
}

// 404
func TestGetLearner_NotFound(t *testing.T) {
	table := "learners"
	userID := uint(42)
	learnerID := uint(999)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			ExpSelectByIDEmpty(table, learnerID)(mock)
			*uID = userID
		},
		http.StatusNotFound,
		http.MethodGet,
		fmt.Sprintf("/learners/%d", learnerID),
	)
}

// 500
func TestGetLearner_DBError(t *testing.T) {
	table := "learners"
	userID := uint(42)
	learnerID := uint(7)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			ExpSelectByIDError(table, learnerID, fmt.Errorf("select failed"))(mock)
			*uID = userID
		},
		http.StatusInternalServerError,
		http.MethodGet,
		fmt.Sprintf("/learners/%d", learnerID),
	)
}

// 400
func TestGetLearner_BadRequest(t *testing.T) {
	userID := uint(42)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			*uID = userID
		},
		http.StatusBadRequest,
		http.MethodGet,
		"/learners/not-an-int",
	)
}

/* ------------------ DeleteLearner ------------------ */

// 200
func TestDeleteLearner_OK_SoftDelete(t *testing.T) {
	table := "learners"
	userID := uint(42)
	learnerID := uint(5)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			ExpSelectByIDFound(table, learnerID, []string{"id"}, []any{learnerID})(mock)
			ExpClearInterestedForLearner(learnerID)(mock)
			ExpSoftDeleteOK(table)(mock)
			*uID = userID
		},
		http.StatusOK,
		http.MethodDelete,
		fmt.Sprintf("/learners/%d", learnerID),
	)
}

// 404
func TestDeleteLearner_NotFound(t *testing.T) {
	table := "learners"
	userID := uint(42)
	learnerID := uint(12345)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			ExpSelectByIDEmpty(table, learnerID)(mock)
			*uID = userID
		},
		http.StatusNotFound,
		http.MethodDelete,
		fmt.Sprintf("/learners/%d", learnerID),
	)
}

// 500
func TestDeleteLearner_DBError(t *testing.T) {
	table := "learners"
	userID := uint(42)
	learnerID := uint(5)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			ExpSelectByIDFound(table, learnerID, []string{"id"}, []any{learnerID})(mock)
			ExpClearInterestedForLearner(learnerID)(mock)
			ExpSoftDeleteError(table, fmt.Errorf("update failed"))(mock)
			*uID = userID
		},
		http.StatusInternalServerError,
		http.MethodDelete,
		fmt.Sprintf("/learners/%d", learnerID),
	)
}

// 400
func TestDeleteLearner_BadRequest(t *testing.T) {
	userID := uint(42)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			*uID = userID
		},
		http.StatusBadRequest,
		http.MethodDelete,
		"/learners/not-an-int",
	)
}

/* ------------------ AddLearnerInterests ------------------ */
func TestAddLearnerInterests_OK(t *testing.T) {
	userID := uint(42)
	learnerID := uint(7)
	newCats := []uint{2, 4, 6}

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uid *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			mock.ExpectQuery(regexp.QuoteMeta(
				`SELECT * FROM "learners" WHERE "learners"."id" = $1 AND "learners"."deleted_at" IS NULL ORDER BY "learners"."id" LIMIT $2`,
			)).
				WithArgs(driver.Value(int64(learnerID)), driver.Value(1)).
				WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(learnerID)))
			ExpPreloadLearnerInterestedEmpty(learnerID)(mock)
			ExpSelectCategoriesByIDs(newCats...)(mock)
			ExpAppendLearnerInterests(learnerID, newCats...)(mock)
			ExpPreloadLearnerInterestedOrdered(learnerID, newCats...)(mock)
			type body struct {
				CategoryIDs []int `json:"class_category_ids"`
			}
			b, _ := json.Marshal(body{CategoryIDs: []int{2, 4, 6}})
			*payload = b
			*uid = userID
		},
		http.StatusOK,
		http.MethodPost,
		fmt.Sprintf("/learners/%d/interests", learnerID),
	)
}

/* ------------------ DeleteLearnerInterests ------------------ */
func TestDeleteLearnerInterests_OK(t *testing.T) {
	userID := uint(42)
	learnerID := uint(7)
	delCats := []uint{2, 4}
	remaining := []uint{6}

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uid *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			mock.ExpectQuery(regexp.QuoteMeta(
				`SELECT * FROM "learners" WHERE "learners"."id" = $1 AND "learners"."deleted_at" IS NULL ORDER BY "learners"."id" LIMIT $2`,
			)).
				WithArgs(driver.Value(int64(learnerID)), driver.Value(1)).
				WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(learnerID)))
			ExpSelectCategoriesByIDs(delCats...)(mock)
			ExpDeleteLearnerInterests(learnerID, delCats...)(mock)
			ExpPreloadLearnerInterestedOrdered(learnerID, remaining...)(mock)
			type body struct {
				CategoryIDs []int `json:"class_category_ids"`
			}
			b, _ := json.Marshal(body{CategoryIDs: []int{2, 4}})
			*payload = b
			*uid = userID
		},
		http.StatusOK,
		http.MethodDelete,
		fmt.Sprintf("/learners/%d/interests", learnerID),
	)
}

/* ------------------ GetClassInterests ------------------ */
func TestGetClassInterestsByLearnerID_OK(t *testing.T) {
	userID := uint(42)
	learnerID := uint(7)
	cats := []uint{3, 5}

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uid *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			mock.ExpectQuery(regexp.QuoteMeta(
				`SELECT * FROM "learners" WHERE "learners"."id" = $1 AND "learners"."deleted_at" IS NULL ORDER BY "learners"."id" LIMIT $2`,
			)).
				WithArgs(driver.Value(int64(learnerID)), driver.Value(1)).
				WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(learnerID)))
			ExpPreloadLearnerInterestedOrdered(learnerID, cats...)(mock)

			*uid = userID
		},
		http.StatusOK,
		http.MethodGet,
		fmt.Sprintf("/learners/%d/interests", learnerID),
	)
}
