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

/* ------------------ CreateReview ------------------ */

// 201
func TestCreateReview_OK(t *testing.T) {
	table := "reviews"
	userID := uint(42)
	learnerID := uint(10)
	classID := uint(5)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, true)(mock)
			ExpInsertReturningID(table, 1)(mock)

			req := jsonBody(models.Review{
				LearnerID: learnerID,
				ClassID:   classID,
				Rating:    5,
				Comment:   "Lorem Ipsum",
			})
			*payload = req
			*uID = userID
		},
		http.StatusCreated,
		http.MethodPost,
		"/reviews/",
	)
}

// 400
func TestCreateReview_BadRequest(t *testing.T) {
	userID := uint(42)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, true)(mock)
			*uID = userID
		},
		http.StatusBadRequest,
		http.MethodPost,
		"/reviews/",
	)
}

// 500
func TestCreateReview_DBError(t *testing.T) {
	table := "reviews"
	userID := uint(42)
	learnerID := uint(10)
	classID := uint(5)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, true)(mock)
			ExpInsertError(table, fmt.Errorf("db insert failed"))(mock)

			req := jsonBody(models.Review{
				LearnerID: learnerID,
				ClassID:   classID,
				Rating:    5,
				Comment:   "Lorem Ipsum",
			})
			*payload = req
			*uID = userID
		},
		http.StatusInternalServerError,
		http.MethodPost,
		"/reviews/",
	)
}

/* ------------------ GetReviews ------------------ */

// 200
func TestGetReviews_OK(t *testing.T) {
	userID := uint(42)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			ExpListRows("reviews", []string{"id"}, []any{1}, []any{2})(mock)
			*uID = userID
		},
		http.StatusOK,
		http.MethodGet,
		"/reviews/",
	)
}

// 500
func TestGetReviews_DBError(t *testing.T) {
	userID := uint(42)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			ExpListError("reviews", fmt.Errorf("select failed"))(mock)
			*uID = userID
		},
		http.StatusInternalServerError,
		http.MethodGet,
		"/reviews/",
	)
}

/* ------------------ GetReview ------------------ */

// 200
func TestGetReview_OK(t *testing.T) {
	table := "reviews"
	userID := uint(42)
	reviewID := uint(7)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			ExpSelectByIDFound(table, reviewID, []string{"id"}, []any{reviewID})(mock)
			*uID = userID
		},
		http.StatusOK,
		http.MethodGet,
		fmt.Sprintf("/reviews/%d", reviewID),
	)
}

// 404
func TestGetReview_NotFound(t *testing.T) {
	table := "reviews"
	userID := uint(42)
	reviewID := uint(999)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			ExpSelectByIDEmpty(table, reviewID)(mock)
			*uID = userID
		},
		http.StatusNotFound,
		http.MethodGet,
		fmt.Sprintf("/reviews/%d", reviewID),
	)
}

// 500
func TestGetReview_DBError(t *testing.T) {
	table := "reviews"
	userID := uint(42)
	reviewID := uint(7)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			ExpSelectByIDError(table, reviewID, fmt.Errorf("select failed"))(mock)
			*uID = userID
		},
		http.StatusInternalServerError,
		http.MethodGet,
		fmt.Sprintf("/reviews/%d", reviewID),
	)
}

// 400
func TestGetReview_BadRequest(t *testing.T) {
	userID := uint(42)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			*uID = userID
		},
		http.StatusBadRequest,
		http.MethodGet,
		"/reviews/not-an-int",
	)
}

/* ------------------ UpdateReview ------------------ */

// 200
func TestUpdateReview_OK(t *testing.T) {
	table := "reviews"
	preloadTable := "classes"
	preloadTable2 := "learners"
	userID := uint(42)
	reviewID := uint(1)
	learnerID := uint(10)
	classID := uint(5)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, true)(mock)
			ExpSelectByIDFound(table, reviewID,
				[]string{"id", "learner_id", "class_id", "rating", "comment"},
				[]any{reviewID, learnerID, classID, 4, "Lorem"},
			)(mock)

			ExpPreloadField(preloadTable, []string{"id"}, []any{classID})(mock)
			ExpPreloadField(preloadTable2, []string{"id"}, []any{learnerID})(mock)
			ExpUpdateOK(table)(mock)

			req := jsonBody(models.Review{
				LearnerID: learnerID,
				ClassID:   classID,
				Rating:    5,
				Comment:   "Lorem Ipsum",
			})
			*payload = req
			*uID = userID
		},
		http.StatusOK,
		http.MethodPut,
		fmt.Sprintf("/reviews/%d", reviewID),
	)
}

// 404
func TestUpdateReview_NotFound(t *testing.T) {
	table := "reviews"
	userID := uint(42)
	reviewID := uint(12345)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, true)(mock)
			ExpSelectByIDEmpty(table, reviewID)(mock)
			*uID = userID
		},
		http.StatusNotFound,
		http.MethodPut,
		fmt.Sprintf("/reviews/%d", reviewID),
	)
}

// 500
func TestUpdateReview_DBError(t *testing.T) {
	table := "reviews"
	userID := uint(42)
	reviewID := uint(1)
	learnerID := uint(10)
	classID := uint(5)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, true)(mock)
			ExpSelectByIDFound(table, reviewID, []string{"id"}, []any{reviewID})(mock)
			ExpUpdateError(table, fmt.Errorf("update failed"))(mock)

			req := jsonBody(models.Review{
				LearnerID: learnerID,
				ClassID:   classID,
				Rating:    5,
				Comment:   "Lorem Ipsum",
			})
			*payload = req
			*uID = userID
		},
		http.StatusInternalServerError,
		http.MethodPut,
		fmt.Sprintf("/reviews/%d", reviewID),
	)
}

// 400
func TestUpdateReview_BadRequest(t *testing.T) {
	userID := uint(42)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, true)(mock)
			*uID = userID
		},
		http.StatusBadRequest,
		http.MethodPut,
		"/reviews/not-an-int",
	)
}

/* ------------------ DeleteReview ------------------ */

// 200
func TestDeleteReview_OK_SoftDelete(t *testing.T) {
	table := "reviews"
	userID := uint(42)
	reviewID := uint(5)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, true)(mock)
			ExpSelectByIDFound(table, reviewID, []string{"id"}, []any{reviewID})(mock)
			ExpSoftDeleteOK(table)(mock)
			*uID = userID
		},
		http.StatusOK,
		http.MethodDelete,
		fmt.Sprintf("/reviews/%d", reviewID),
	)
}

// 404
func TestDeleteReview_NotFound(t *testing.T) {
	table := "reviews"
	userID := uint(42)
	reviewID := uint(12345)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, true)(mock)
			ExpSelectByIDEmpty(table, reviewID)(mock)
			*uID = userID
		},
		http.StatusNotFound,
		http.MethodDelete,
		fmt.Sprintf("/reviews/%d", reviewID),
	)
}

// 500
func TestDeleteReview_DBError(t *testing.T) {
	table := "reviews"
	userID := uint(42)
	reviewID := uint(5)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, true)(mock)
			ExpSelectByIDFound(table, reviewID, []string{"id"}, []any{reviewID})(mock)
			ExpSoftDeleteError(table, fmt.Errorf("update failed"))(mock)
			*uID = userID
		},
		http.StatusInternalServerError,
		http.MethodDelete,
		fmt.Sprintf("/reviews/%d", reviewID),
	)
}

// 400
func TestDeleteReview_BadRequest(t *testing.T) {
	userID := uint(42)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, true)(mock)
			*uID = userID
		},
		http.StatusBadRequest,
		http.MethodDelete,
		"/reviews/not-an-int",
	)
}
