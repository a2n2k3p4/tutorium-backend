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

/* ------------------ CreateNotification ------------------ */

// 201
func TestCreateNotification_OK(t *testing.T) {
	table := "notifications"
	userID := uint(42)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, true, false, false)(mock)
			ExpInsertReturningID(table, 1)(mock)

			req := jsonBody(models.Notification{
				UserID:                  userID,
				NotificationType:        "test",
				NotificationDescription: "Lorem ipsum dolor sit amet, consectetur adipiscing elit",
				NotificationDate:        time.Now(),
				ReadFlag:                false,
			})
			*payload = req
			*uID = userID
		},
		http.StatusCreated,
		http.MethodPost,
		"/notifications/",
	)
}

// 400
func TestCreateNotification_BadRequest(t *testing.T) {
	userID := uint(42)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, true, false, false)(mock)
			*uID = userID
		},
		http.StatusBadRequest,
		http.MethodPost,
		"/notifications/",
	)
}

// 500
func TestCreateNotification_DBError(t *testing.T) {
	table := "notifications"
	userID := uint(42)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, true, false, false)(mock)
			ExpInsertError(table, fmt.Errorf("db insert failed"))(mock)
			req := jsonBody(models.Notification{
				UserID:                  userID,
				NotificationType:        "test",
				NotificationDescription: "Lorem ipsum dolor sit amet, consectetur adipiscing elit",
				NotificationDate:        time.Now(),
				ReadFlag:                false,
			})
			*payload = req
			*uID = userID
		},
		http.StatusInternalServerError,
		http.MethodPost,
		"/notifications/",
	)
}

/* ------------------ GetNotifications ------------------ */

// 200
func TestGetNotifications_OK(t *testing.T) {
	userID := uint(42)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			ExpListRows("notifications", []string{"id"}, []any{1}, []any{2})(mock)
			*uID = userID
		},
		http.StatusOK,
		http.MethodGet,
		"/notifications/",
	)
}

// 500
func TestGetNotifications_DBError(t *testing.T) {
	userID := uint(42)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			ExpListError("notifications", fmt.Errorf("select failed"))(mock)
			*uID = userID
		},
		http.StatusInternalServerError,
		http.MethodGet,
		"/notifications/",
	)
}

/* ------------------ GetNotification ------------------ */

// 200
func TestGetNotification_OK(t *testing.T) {
	table := "notifications"
	userID := uint(42)
	notificationID := uint(7)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			ExpSelectByIDFound(table, notificationID, []string{"id"}, []any{notificationID})(mock)
			*uID = userID
		},
		http.StatusOK,
		http.MethodGet,
		fmt.Sprintf("/notifications/%d", notificationID),
	)
}

// 404
func TestGetNotification_NotFound(t *testing.T) {
	table := "notifications"
	userID := uint(42)
	notificationID := uint(999)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			ExpSelectByIDEmpty(table, notificationID)(mock)
			*uID = userID
		},
		http.StatusNotFound,
		http.MethodGet,
		fmt.Sprintf("/notifications/%d", notificationID),
	)
}

// 500
func TestGetNotification_DBError(t *testing.T) {
	table := "notifications"
	userID := uint(42)
	notificationID := uint(7)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			ExpSelectByIDError(table, notificationID, fmt.Errorf("select failed"))(mock)
			*uID = userID
		},
		http.StatusInternalServerError,
		http.MethodGet,
		fmt.Sprintf("/notifications/%d", notificationID),
	)
}

// 400
func TestGetNotification_BadRequest(t *testing.T) {
	userID := uint(42)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			*uID = userID
		},
		http.StatusBadRequest,
		http.MethodGet,
		"/notifications/not-an-int",
	)
}

/* ------------------ UpdateNotification ------------------ */

// 200
func TestUpdateNotification_OK(t *testing.T) {
	table := "notifications"
	preloadTable := "users"
	userID := uint(42)
	notificationID := uint(1)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			ExpSelectByIDFound(table, notificationID,
				[]string{"id", "user_id", "notification_type", "notification_description", "notification_date", "read_flag"},
				[]any{notificationID, userID, "original type", "Lorem", time.Now(), false},
			)(mock)

			ExpPreloadField(preloadTable, []string{"id"}, []any{userID})(mock)
			ExpUpdateOK(table)(mock)

			req := jsonBody(models.Notification{
				UserID:                  userID,
				NotificationType:        "edit test",
				NotificationDescription: "Lorem ipsum dolor sit amet, consectetur adipiscing elit",
				NotificationDate:        time.Now(),
				ReadFlag:                false,
			})
			*payload = req
			*uID = userID
		},
		http.StatusOK,
		http.MethodPut,
		fmt.Sprintf("/notifications/%d", notificationID),
	)
}

// 404
func TestUpdateNotification_NotFound(t *testing.T) {
	table := "notifications"
	userID := uint(42)
	notificationID := uint(12345)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			ExpSelectByIDEmpty(table, notificationID)(mock)
			*uID = userID
		},
		http.StatusNotFound,
		http.MethodPut,
		fmt.Sprintf("/notifications/%d", notificationID),
	)
}

// 500
func TestUpdateNotification_DBError(t *testing.T) {
	table := "notifications"
	userID := uint(42)
	notificationID := uint(1)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			ExpSelectByIDFound(table, notificationID, []string{"id"}, []any{notificationID})(mock)
			ExpUpdateError(table, fmt.Errorf("update failed"))(mock)

			req := jsonBody(models.Notification{
				UserID:                  userID,
				NotificationType:        "edit test",
				NotificationDescription: "Lorem ipsum dolor sit amet, consectetur adipiscing elit",
				NotificationDate:        time.Now(),
				ReadFlag:                false,
			})
			*payload = req
			*uID = userID
		},
		http.StatusInternalServerError,
		http.MethodPut,
		fmt.Sprintf("/notifications/%d", notificationID),
	)
}

// 400
func TestUpdateNotification_BadRequest(t *testing.T) {
	userID := uint(42)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			*uID = userID
		},
		http.StatusBadRequest,
		http.MethodPut,
		"/notifications/not-an-int",
	)
}

/* ------------------ DeleteNotification ------------------ */

// 200
func TestDeleteNotification_OK_SoftDelete(t *testing.T) {
	table := "notifications"
	userID := uint(42)
	notificationID := uint(5)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			ExpSelectByIDFound(table, notificationID, []string{"id"}, []any{notificationID})(mock)
			ExpSoftDeleteOK(table)(mock)
			*uID = userID
		},
		http.StatusOK,
		http.MethodDelete,
		fmt.Sprintf("/notifications/%d", notificationID),
	)
}

// 404
func TestDeleteNotification_NotFound(t *testing.T) {
	table := "notifications"
	userID := uint(42)
	notificationID := uint(12345)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			ExpSelectByIDEmpty(table, notificationID)(mock)
			*uID = userID
		},
		http.StatusNotFound,
		http.MethodDelete,
		fmt.Sprintf("/notifications/%d", notificationID),
	)
}

// 500
func TestDeleteNotification_DBError(t *testing.T) {
	table := "notifications"
	userID := uint(42)
	notificationID := uint(5)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			ExpSelectByIDFound(table, notificationID, []string{"id"}, []any{notificationID})(mock)
			ExpSoftDeleteError(table, fmt.Errorf("update failed"))(mock)
			*uID = userID
		},
		http.StatusInternalServerError,
		http.MethodDelete,
		fmt.Sprintf("/notifications/%d", notificationID),
	)
}

// 400
func TestDeleteNotification_BadRequest(t *testing.T) {
	userID := uint(42)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			*uID = userID
		},
		http.StatusBadRequest,
		http.MethodDelete,
		"/notifications/not-an-int",
	)
}
