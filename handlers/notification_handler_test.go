package handlers

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/a2n2k3p4/tutorium-backend/models"
)

/* ------------------ CreateNotification ------------------ */

// 201
func TestCreateNotification_OK(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)

		table := "notifications"
		userID := uint(42)

		ExpAuthUser(userID, true, false, false)(mock)
		ExpInsertReturningID(table, 1)(mock)

		app := setupApp(gdb)

		payload := models.Notification{
			UserID:                  userID,
			NotificationType:        "test",
			NotificationDescription: "Lorem ipsum dolor sit amet, consectetur adipiscing elit",
			NotificationDate:        time.Now(),
			ReadFlag:                false,
		}

		resp := runHTTP(t, app, httpInput{
			Method:      http.MethodPost,
			Path:        "/notifications/",
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
func TestCreateNotification_BadRequest(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)
		userID := uint(42)

		ExpAuthUser(userID, true, false, false)(mock)

		app := setupApp(gdb)
		resp := runHTTP(t, app, httpInput{
			Method: http.MethodPost, Path: "/notifications/",
			Body: []byte(`{invalid-json}`), ContentType: "application/json", UserID: &userID,
		})
		wantStatus(t, resp, http.StatusBadRequest)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

// 500
func TestCreateNotification_DBError(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)

		table := "notifications"
		userID := uint(42)

		ExpAuthUser(userID, true, false, false)(mock)
		ExpInsertError(table, fmt.Errorf("db insert failed"))(mock)

		app := setupApp(gdb)
		resp := runHTTP(t, app, httpInput{
			Method: http.MethodPost, Path: "/notifications/",
			Body: []byte(`{}`), ContentType: "application/json", UserID: &userID,
		})
		wantStatus(t, resp, http.StatusInternalServerError)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

/* ------------------ GetNotifications ------------------ */
// 200
func TestGetNotifications_OK(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)
		userID := uint(42)

		ExpAuthUser(userID, false, false, false)(mock)
		ExpListRows("notifications", []string{"id"}, []any{1}, []any{2})(mock)

		app := setupApp(gdb)
		resp := runHTTP(t, app, httpInput{
			Method: http.MethodGet, Path: "/notifications/", UserID: &userID,
		})
		wantStatus(t, resp, http.StatusOK)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

// 500
func TestGetNotifications_DBError(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)
		userID := uint(42)

		ExpAuthUser(userID, false, false, false)(mock)
		ExpListError("notifications", fmt.Errorf("select failed"))(mock)

		app := setupApp(gdb)
		resp := runHTTP(t, app, httpInput{
			Method: http.MethodGet, Path: "/notifications/", UserID: &userID,
		})
		wantStatus(t, resp, http.StatusInternalServerError)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

/* ------------------ GetNotification ------------------ */
// 200
func TestGetNotification_OK(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)

		table := "notifications"
		userID := uint(42)
		notificationID := uint(7)

		ExpAuthUser(userID, false, false, false)(mock)
		ExpSelectByIDFound(table, notificationID, []string{"id"}, []any{notificationID})(mock)

		app := setupApp(gdb)
		resp := runHTTP(t, app, httpInput{
			Method: http.MethodGet, Path: fmt.Sprintf("/notifications/%d", notificationID), UserID: &userID,
		})
		wantStatus(t, resp, http.StatusOK)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

// 404
func TestGetNotification_NotFound(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)

		table := "notifications"
		userID := uint(42)
		notificationID := uint(999)

		ExpAuthUser(userID, false, false, false)(mock)
		ExpSelectByIDEmpty(table, notificationID)(mock)

		app := setupApp(gdb)
		resp := runHTTP(t, app, httpInput{
			Method: http.MethodGet, Path: fmt.Sprintf("/notifications/%d", notificationID), UserID: &userID,
		})
		wantStatus(t, resp, http.StatusNotFound)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

// 500
func TestGetNotification_DBError(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)

		table := "notifications"
		userID := uint(42)
		notificationID := uint(7)

		ExpAuthUser(userID, false, false, false)(mock)
		ExpSelectByIDError(table, notificationID, fmt.Errorf("select failed"))(mock)

		app := setupApp(gdb)
		resp := runHTTP(t, app, httpInput{
			Method: http.MethodGet, Path: fmt.Sprintf("/notifications/%d", notificationID), UserID: &userID,
		})
		wantStatus(t, resp, http.StatusInternalServerError)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

// 400
func TestGetNotification_BadRequest(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)
		userID := uint(42)

		ExpAuthUser(userID, false, false, false)(mock)

		app := setupApp(gdb)
		resp := runHTTP(t, app, httpInput{
			Method: http.MethodGet, Path: "/notifications/not-an-int", UserID: &userID,
		})
		wantStatus(t, resp, http.StatusBadRequest)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

/* ------------------ UpdateNotification ------------------ */

// 200
func TestUpdateNotification_OK(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)

		table := "notifications"
		preloadTable := "users"
		userID := uint(42)
		notificationID := uint(1)

		ExpAuthUser(userID, false, false, false)(mock)
		ExpSelectByIDFound(table, notificationID,
			[]string{"id", "user_id", "notification_type", "notification_description", "notification_date", "read_flag"},
			[]any{notificationID, userID, "original type", "Lorem", time.Now(), false},
		)(mock)

		ExpPreloadField(preloadTable, []string{"id"}, []any{userID})(mock)
		ExpUpdateOK(table)(mock)

		app := setupApp(gdb)
		payload := models.Notification{
			UserID: userID, NotificationType: "edit test",
			NotificationDescription: "Lorem ipsum dolor sit amet, consectetur adipiscing elit",
			NotificationDate:        time.Now(), ReadFlag: false,
		}
		resp := runHTTP(t, app, httpInput{
			Method: http.MethodPut, Path: fmt.Sprintf("/notifications/%d", notificationID),
			Body: jsonBody(payload), ContentType: "application/json", UserID: &userID,
		})
		wantStatus(t, resp, http.StatusOK)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

// 404
func TestUpdateNotification_NotFound(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)

		table := "notifications"
		userID := uint(42)
		notificationID := uint(12345)

		ExpAuthUser(userID, false, false, false)(mock)
		ExpSelectByIDEmpty(table, notificationID)(mock)

		app := setupApp(gdb)
		resp := runHTTP(t, app, httpInput{
			Method: http.MethodPut, Path: fmt.Sprintf("/notifications/%d", notificationID), UserID: &userID,
		})
		wantStatus(t, resp, http.StatusNotFound)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

// 500
func TestUpdateNotification_DBError(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)

		table := "notifications"
		userID := uint(42)
		notificationID := uint(1)

		ExpAuthUser(userID, false, false, false)(mock)
		ExpSelectByIDFound(table, notificationID, []string{"id"}, []any{notificationID})(mock)
		ExpUpdateError(table, fmt.Errorf("update failed"))(mock)

		app := setupApp(gdb)
		payload := models.Notification{
			UserID: userID, NotificationType: "edit test",
			NotificationDescription: "Lorem ipsum dolor sit amet, consectetur adipiscing elit",
			NotificationDate:        time.Now(), ReadFlag: false,
		}
		resp := runHTTP(t, app, httpInput{
			Method: http.MethodPut, Path: fmt.Sprintf("/notifications/%d", notificationID),
			Body: jsonBody(payload), ContentType: "application/json", UserID: &userID,
		})
		wantStatus(t, resp, http.StatusInternalServerError)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

// 400
func TestUpdateNotification_BadRequest(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)
		userID := uint(42)

		ExpAuthUser(userID, false, false, false)(mock)

		app := setupApp(gdb)
		resp := runHTTP(t, app, httpInput{
			Method: http.MethodPut, Path: "/notifications/not-an-int", UserID: &userID,
		})
		wantStatus(t, resp, http.StatusBadRequest)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

/* ------------------ DeleteNotification ------------------ */

// 200 (soft delete)
func TestDeleteNotification_OK_SoftDelete(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)

		table := "notifications"
		userID := uint(42)
		notificationID := uint(5)

		ExpAuthUser(userID, false, false, false)(mock)
		ExpSelectByIDFound(table, notificationID, []string{"id"}, []any{notificationID})(mock)
		ExpSoftDeleteOK(table)(mock)

		app := setupApp(gdb)
		resp := runHTTP(t, app, httpInput{
			Method: http.MethodDelete, Path: fmt.Sprintf("/notifications/%d", notificationID), UserID: &userID,
		})
		wantStatus(t, resp, http.StatusOK)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

// 404
func TestDeleteNotification_NotFound(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)

		table := "notifications"
		userID := uint(42)
		notificationID := uint(12345)

		ExpAuthUser(userID, false, false, false)(mock)
		ExpSelectByIDEmpty(table, notificationID)(mock)

		app := setupApp(gdb)
		resp := runHTTP(t, app, httpInput{
			Method: http.MethodDelete, Path: fmt.Sprintf("/notifications/%d", notificationID), UserID: &userID,
		})
		wantStatus(t, resp, http.StatusNotFound)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

// 500
func TestDeleteNotification_DBError(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)

		table := "notifications"
		userID := uint(42)
		notificationID := uint(5)

		ExpAuthUser(userID, false, false, false)(mock)
		ExpSelectByIDFound(table, notificationID, []string{"id"}, []any{notificationID})(mock)
		ExpSoftDeleteError(table, fmt.Errorf("update failed"))(mock)

		app := setupApp(gdb)
		resp := runHTTP(t, app, httpInput{
			Method: http.MethodDelete, Path: fmt.Sprintf("/notifications/%d", notificationID), UserID: &userID,
		})
		wantStatus(t, resp, http.StatusInternalServerError)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

// 400
func TestDeleteNotification_BadRequest(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)
		userID := uint(42)

		ExpAuthUser(userID, false, false, false)(mock)

		app := setupApp(gdb)
		resp := runHTTP(t, app, httpInput{
			Method: http.MethodDelete, Path: "/notifications/not-an-int", UserID: &userID,
		})
		wantStatus(t, resp, http.StatusBadRequest)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}
