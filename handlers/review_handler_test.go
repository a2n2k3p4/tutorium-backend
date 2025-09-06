package handlers

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/a2n2k3p4/tutorium-backend/models"
)

/* ------------------ CreateReview ------------------ */

// 201
func TestCreateReview_OK(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)

		table := "reviews"
		userID := uint(42)
		learnerID := uint(10)
		classID := uint(5)

		ExpAuthUser(userID, false, false, true)(mock)
		ExpInsertReturningID(table, 1)(mock)

		app := setupApp(gdb)

		payload := models.Review{
			LearnerID: learnerID,
			ClassID:   classID,
			Rating:    5,
			Comment:   "Lorem Ipsum",
		}

		resp := runHTTP(t, app, httpInput{
			Method:      http.MethodPost,
			Path:        "/reviews/",
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

// 400
func TestCreateReview_BadRequest(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)
		userID := uint(42)

		ExpAuthUser(userID, false, false, true)(mock)

		app := setupApp(gdb)
		resp := runHTTP(t, app, httpInput{
			Method: http.MethodPost, Path: "/reviews/",
			Body: []byte(`{invalid-json}`), ContentType: "application/json", UserID: &userID,
		})
		wantStatus(t, resp, http.StatusBadRequest)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

// 500
func TestCreateReview_DBError(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)

		table := "reviews"
		userID := uint(42)
		learnerID := uint(10)
		classID := uint(5)

		ExpAuthUser(userID, false, false, true)(mock)
		ExpInsertError(table, fmt.Errorf("db insert failed"))(mock)

		payload := models.Review{
			LearnerID: learnerID,
			ClassID:   classID,
			Rating:    5,
			Comment:   "Lorem Ipsum",
		}
		app := setupApp(gdb)
		resp := runHTTP(t, app, httpInput{
			Method: http.MethodPost, Path: "/reviews/",
			Body: jsonBody(payload), ContentType: "application/json", UserID: &userID,
		})
		wantStatus(t, resp, http.StatusInternalServerError)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

/* ------------------ GetReviews ------------------ */
// 200
func TestGetReviews_OK(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)
		userID := uint(42)

		ExpAuthUser(userID, false, false, false)(mock)
		ExpListRows("reviews", []string{"id"}, []any{1}, []any{2})(mock)

		app := setupApp(gdb)
		resp := runHTTP(t, app, httpInput{
			Method: http.MethodGet, Path: "/reviews/", UserID: &userID,
		})
		wantStatus(t, resp, http.StatusOK)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

// 500
func TestGetReviews_DBError(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)
		userID := uint(42)

		ExpAuthUser(userID, false, false, false)(mock)
		ExpListError("reviews", fmt.Errorf("select failed"))(mock)

		app := setupApp(gdb)
		resp := runHTTP(t, app, httpInput{
			Method: http.MethodGet, Path: "/reviews/", UserID: &userID,
		})
		wantStatus(t, resp, http.StatusInternalServerError)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

/* ------------------ GetReview ------------------ */
// 200
func TestGetReview_OK(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)

		table := "reviews"
		userID := uint(42)
		reviewID := uint(7)

		ExpAuthUser(userID, false, false, false)(mock)
		ExpSelectByIDFound(table, reviewID, []string{"id"}, []any{reviewID})(mock)

		app := setupApp(gdb)
		resp := runHTTP(t, app, httpInput{
			Method: http.MethodGet, Path: fmt.Sprintf("/reviews/%d", reviewID), UserID: &userID,
		})
		wantStatus(t, resp, http.StatusOK)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

// 404
func TestGetReview_NotFound(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)

		table := "reviews"
		userID := uint(42)
		reviewID := uint(999)

		ExpAuthUser(userID, false, false, false)(mock)
		ExpSelectByIDEmpty(table, reviewID)(mock)

		app := setupApp(gdb)
		resp := runHTTP(t, app, httpInput{
			Method: http.MethodGet, Path: fmt.Sprintf("/reviews/%d", reviewID), UserID: &userID,
		})
		wantStatus(t, resp, http.StatusNotFound)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

// 500
func TestGetReview_DBError(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)

		table := "reviews"
		userID := uint(42)
		reviewID := uint(7)

		ExpAuthUser(userID, false, false, false)(mock)
		ExpSelectByIDError(table, reviewID, fmt.Errorf("select failed"))(mock)

		app := setupApp(gdb)
		resp := runHTTP(t, app, httpInput{
			Method: http.MethodGet, Path: fmt.Sprintf("/reviews/%d", reviewID), UserID: &userID,
		})
		wantStatus(t, resp, http.StatusInternalServerError)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

// 400
func TestGetReview_BadRequest(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)
		userID := uint(42)

		ExpAuthUser(userID, false, false, false)(mock)

		app := setupApp(gdb)
		resp := runHTTP(t, app, httpInput{
			Method: http.MethodGet, Path: "/reviews/not-an-int", UserID: &userID,
		})
		wantStatus(t, resp, http.StatusBadRequest)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

/* ------------------ UpdateReview ------------------ */

// 200
func TestUpdateReview_OK(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)

		table := "reviews"
		const preloadTable = "classes"
		const preloadTable2 = "learners"
		userID := uint(42)
		reviewID := uint(1)
		learnerID := uint(10)
		classID := uint(5)

		ExpAuthUser(userID, false, false, true)(mock)
		ExpSelectByIDFound(table, reviewID,
			[]string{"id", "learner_id", "class_id", "rating", "comment"},
			[]any{reviewID, learnerID, classID, 4, "Lorem"},
		)(mock)

		ExpPreloadField(preloadTable, []string{"id"}, []any{classID})(mock)
		ExpPreloadField(preloadTable2, []string{"id"}, []any{learnerID})(mock)

		ExpUpdateOK(table)(mock)

		app := setupApp(gdb)
		payload := models.Review{
			LearnerID: learnerID,
			ClassID:   classID,
			Rating:    5,
			Comment:   "Lorem Ipsum",
		}
		resp := runHTTP(t, app, httpInput{
			Method: http.MethodPut, Path: fmt.Sprintf("/reviews/%d", reviewID),
			Body: jsonBody(payload), ContentType: "application/json", UserID: &userID,
		})
		wantStatus(t, resp, http.StatusOK)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

// 404
func TestUpdateReview_NotFound(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)

		table := "reviews"
		userID := uint(42)
		reviewID := uint(12345)

		ExpAuthUser(userID, false, false, true)(mock)
		ExpSelectByIDEmpty(table, reviewID)(mock)

		app := setupApp(gdb)
		resp := runHTTP(t, app, httpInput{
			Method: http.MethodPut, Path: fmt.Sprintf("/reviews/%d", reviewID), UserID: &userID,
		})
		wantStatus(t, resp, http.StatusNotFound)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

// 500
func TestUpdateReview_DBError(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)

		table := "reviews"
		userID := uint(42)
		reviewID := uint(1)
		learnerID := uint(10)
		classID := uint(5)

		ExpAuthUser(userID, false, false, true)(mock)
		ExpSelectByIDFound(table, reviewID, []string{"id"}, []any{reviewID})(mock)
		ExpUpdateError(table, fmt.Errorf("update failed"))(mock)

		app := setupApp(gdb)
		payload := models.Review{
			LearnerID: learnerID,
			ClassID:   classID,
			Rating:    5,
			Comment:   "Lorem Ipsum",
		}
		resp := runHTTP(t, app, httpInput{
			Method: http.MethodPut, Path: fmt.Sprintf("/reviews/%d", reviewID),
			Body: jsonBody(payload), ContentType: "application/json", UserID: &userID,
		})
		wantStatus(t, resp, http.StatusInternalServerError)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

// 400
func TestUpdateReview_BadRequest(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)
		userID := uint(42)

		ExpAuthUser(userID, false, false, true)(mock)

		app := setupApp(gdb)
		resp := runHTTP(t, app, httpInput{
			Method: http.MethodPut, Path: "/reviews/not-an-int", UserID: &userID,
		})
		wantStatus(t, resp, http.StatusBadRequest)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

/* ------------------ DeleteReview ------------------ */

// 200
func TestDeleteReview_OK_SoftDelete(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)

		table := "reviews"
		userID := uint(42)
		reviewID := uint(5)

		ExpAuthUser(userID, false, false, true)(mock)
		ExpSelectByIDFound(table, reviewID, []string{"id"}, []any{reviewID})(mock)
		ExpSoftDeleteOK(table)(mock)

		app := setupApp(gdb)
		resp := runHTTP(t, app, httpInput{
			Method: http.MethodDelete, Path: fmt.Sprintf("/reviews/%d", reviewID), UserID: &userID,
		})
		wantStatus(t, resp, http.StatusOK)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

// 404
func TestDeleteReview_NotFound(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)

		table := "reviews"
		userID := uint(42)
		reviewID := uint(12345)

		ExpAuthUser(userID, false, false, true)(mock)
		ExpSelectByIDEmpty(table, reviewID)(mock)

		app := setupApp(gdb)
		resp := runHTTP(t, app, httpInput{
			Method: http.MethodDelete, Path: fmt.Sprintf("/reviews/%d", reviewID), UserID: &userID,
		})
		wantStatus(t, resp, http.StatusNotFound)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

// 500
func TestDeleteReview_DBError(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)

		table := "reviews"
		userID := uint(42)
		reviewID := uint(5)

		ExpAuthUser(userID, false, false, true)(mock)
		ExpSelectByIDFound(table, reviewID, []string{"id"}, []any{reviewID})(mock)
		ExpSoftDeleteError(table, fmt.Errorf("update failed"))(mock)

		app := setupApp(gdb)
		resp := runHTTP(t, app, httpInput{
			Method: http.MethodDelete, Path: fmt.Sprintf("/reviews/%d", reviewID), UserID: &userID,
		})
		wantStatus(t, resp, http.StatusInternalServerError)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

// 400
func TestDeleteReview_BadRequest(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)
		userID := uint(42)

		ExpAuthUser(userID, false, false, true)(mock)

		app := setupApp(gdb)
		resp := runHTTP(t, app, httpInput{
			Method: http.MethodDelete, Path: "/reviews/not-an-int", UserID: &userID,
		})
		wantStatus(t, resp, http.StatusBadRequest)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}
