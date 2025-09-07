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

/* ------------------ CreateBanLearner ------------------ */

// 201
func TestCreateBanLearner_OK(t *testing.T) {
	table := "ban_details_learners"
	userID := uint(42)
	learnerID := uint(50)
	now := time.Now()
	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, true, false, false)(mock)
			ExpInsertReturningID(table, 1)(mock)

			req := jsonBody(models.BanDetailsLearner{
				LearnerID:      learnerID,
				BanStart:       now,
				BanEnd:         now.Add(2 * time.Hour),
				BanDescription: "spamming",
			},
			)

			*payload = req
			*uID = userID

		},
		http.StatusCreated,
		http.MethodPost,
		"/banlearners/",
	)
}

// 400
func TestCreateBanLearner_BadRequest(t *testing.T) {
	userID := uint(42)
	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, true, false, false)(mock)
			*uID = userID
		},
		http.StatusBadRequest,
		http.MethodPost,
		"/banlearners/",
	)
}

// 500
func TestCreateBanLearner_DBError(t *testing.T) {
	table := "ban_details_learners"
	userID := uint(42)
	learnerID := uint(50)
	now := time.Now()
	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, true, false, false)(mock)
			ExpInsertError(table, fmt.Errorf("db insert failed"))(mock)

			req := jsonBody(models.BanDetailsLearner{
				LearnerID:      learnerID,
				BanStart:       now,
				BanEnd:         now.Add(2 * time.Hour),
				BanDescription: "spamming",
			},
			)

			*payload = req
			*uID = userID
		},
		http.StatusInternalServerError,
		http.MethodPost,
		"/banlearners/",
	)
}

/* ------------------ GetBanLearners ------------------ */
// 200
func TestGetBanLearners_OK(t *testing.T) {
	userID := uint(42)
	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, true, false, false)(mock)
			ExpListRows("ban_details_learners", []string{"id"}, []any{1}, []any{2})(mock)
			*uID = userID
		},
		http.StatusOK,
		http.MethodGet,
		"/banlearners/",
	)
}

// 500
func TestGetBanLearners_DBError(t *testing.T) {
	userID := uint(42)
	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, true, false, false)(mock)
			ExpListError("ban_details_learners", fmt.Errorf("select failed"))(mock)
			*uID = userID
		},
		http.StatusInternalServerError,
		http.MethodGet,
		"/banlearners/",
	)
}

/* ------------------ GetBanLearner ------------------ */
// 200
func TestGetBanLearner_OK(t *testing.T) {
	table := "ban_details_learners"
	userID := uint(42)
	ban_learnerID := uint(7)
	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {

			ExpAuthUser(userID, true, false, false)(mock)
			ExpSelectByIDFound(table, ban_learnerID, []string{"id"}, []any{ban_learnerID})(mock)

		},
		http.StatusOK,
		http.MethodGet,
		fmt.Sprintf("/banlearners/%d", ban_learnerID),
	)
}

// 404
func TestGetBanLearner_NotFound(t *testing.T) {
	table := "ban_details_learners"
	userID := uint(42)
	ban_learnerID := uint(999)
	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {

			ExpAuthUser(userID, true, false, false)(mock)
			ExpSelectByIDEmpty(table, ban_learnerID)(mock)

		},
		http.StatusNotFound,
		http.MethodGet,
		fmt.Sprintf("/banlearners/%d", ban_learnerID),
	)
}

// 500
func TestGetBanLearner_DBError(t *testing.T) {
	table := "ban_details_learners"
	userID := uint(42)
	ban_learnerID := uint(7)
	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, true, false, false)(mock)
			ExpSelectByIDError(table, ban_learnerID, fmt.Errorf("select failed"))(mock)

		},
		http.StatusInternalServerError,
		http.MethodGet,
		fmt.Sprintf("/banlearners/%d", ban_learnerID),
	)
}

// 400
func TestGetBanLearner_BadRequest(t *testing.T) {
	userID := uint(42)
	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {

			ExpAuthUser(userID, true, false, false)(mock)

		},
		http.StatusBadRequest,
		http.MethodGet,
		"/banlearners/not-an-int",
	)
}

/* ------------------ UpdateBanLearner ------------------ */

// 200
func TestUpdateBanLearner_OK(t *testing.T) {
	table := "ban_details_learners"
	preloadTable := "learners"
	userID := uint(42)
	learnerID := uint(50)
	ban_learnerID := uint(1)
	now := time.Now()

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {

			ExpAuthUser(userID, true, false, false)(mock)
			ExpSelectByIDFound(table, ban_learnerID,
				[]string{"id", "learner_id", "ban_start", "ban_end", "ban_description"},
				[]any{ban_learnerID, learnerID, now, now.Add(2 * time.Hour), "flooding"},
			)(mock)

			ExpPreloadField(preloadTable, []string{"id"}, []any{learnerID})(mock)

			ExpUpdateOK(table)(mock)

			req := jsonBody(models.BanDetailsLearner{
				BanDescription: "spamming",
			},
			)

			*payload = req

		},
		http.StatusOK,
		http.MethodPut,
		fmt.Sprintf("/banlearners/%d", ban_learnerID),
	)
}

// 404
func TestUpdateBanLearner_NotFound(t *testing.T) {
	table := "ban_details_learners"
	userID := uint(42)
	ban_learnerID := uint(12345)
	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, true, false, false)(mock)
			ExpSelectByIDEmpty(table, ban_learnerID)(mock)

		},
		http.StatusNotFound,
		http.MethodPut,
		fmt.Sprintf("/banlearners/%d", ban_learnerID),
	)
}

// 500
func TestUpdateBanLearner_DBError(t *testing.T) {
	table := "ban_details_learners"
	userID := uint(42)
	ban_learnerID := uint(1)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, true, false, false)(mock)
			ExpSelectByIDFound(table, ban_learnerID, []string{"id"}, []any{ban_learnerID})(mock)
			ExpUpdateError(table, fmt.Errorf("update failed"))(mock)

			req := jsonBody(models.BanDetailsLearner{
				BanDescription: "spamming",
			},
			)

			*payload = req

		},
		http.StatusInternalServerError,
		http.MethodPut,
		fmt.Sprintf("/banlearners/%d", ban_learnerID),
	)
}

// 400
func TestUpdateBanLearner_BadRequest(t *testing.T) {
	userID := uint(42)
	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, true, false, false)(mock)

		},
		http.StatusBadRequest,
		http.MethodPut,
		"/banlearners/not-an-int",
	)
}

/* ------------------ DeleteBanLearner ------------------ */

// 200
func TestDeleteBanLearner_OK_SoftDelete(t *testing.T) {
	table := "ban_details_learners"
	userID := uint(42)
	ban_learnerID := uint(5)
	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, true, false, false)(mock)
			ExpSelectByIDFound(table, ban_learnerID, []string{"id"}, []any{ban_learnerID})(mock)
			ExpSoftDeleteOK(table)(mock)

		},
		http.StatusOK,
		http.MethodDelete,
		fmt.Sprintf("/banlearners/%d", ban_learnerID),
	)
}

// 404
func TestDeleteBanLearner_NotFound(t *testing.T) {
	table := "ban_details_learners"
	userID := uint(42)
	ban_learnerID := uint(12345)
	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, true, false, false)(mock)
			ExpSelectByIDEmpty(table, ban_learnerID)(mock)

		},
		http.StatusNotFound,
		http.MethodDelete,
		fmt.Sprintf("/banlearners/%d", ban_learnerID),
	)
}

// 500
func TestDeleteBanLearner_DBError(t *testing.T) {
	table := "ban_details_learners"
	userID := uint(42)
	ban_learnerID := uint(5)
	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {

			ExpAuthUser(userID, true, false, false)(mock)
			ExpSelectByIDFound(table, ban_learnerID, []string{"id"}, []any{ban_learnerID})(mock)
			ExpSoftDeleteError(table, fmt.Errorf("update failed"))(mock)

		},
		http.StatusInternalServerError,
		http.MethodDelete,
		fmt.Sprintf("/banlearners/%d", ban_learnerID),
	)
}

// 400
func TestDeleteBanLearner_BadRequest(t *testing.T) {
	userID := uint(42)
	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, true, false, false)(mock)

		},
		http.StatusBadRequest,
		http.MethodDelete,
		"/banlearners/not-an-int",
	)
}
