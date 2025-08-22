package handlers

import (
	"errors"

	"github.com/a2n2k3p4/tutorium-backend/middleware"
	"github.com/a2n2k3p4/tutorium-backend/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func NotificationRoutes(app *fiber.App) {
	notification := app.Group("/notification", middleware.ProtectedMiddleware())
	notification.Get("/", GetNotifications)
	notification.Get("/:id", GetNotification)
	notification.Put("/:id", UpdateNotification)
	notification.Delete("/:id", DeleteNotification)

	notificationAdmin := notification.Group("/", middleware.AdminRequired())
	notificationAdmin.Post("/", CreateNotification)
}

// CreateNotification godoc
//
//	@Summary		Create a new notification
//	@Description	CreateNotification creates a new Notification record
//	@Tags			Notifications
//	@Accept			json
//	@Produce		json
//	@Param			notification	body		models.NotificationDoc	true	"Notification payload"
//	@Success		201				{object}	models.NotificationDoc
//	@Failure		400				{object}	map[string]string	"Invalid input"
//	@Failure		500				{object}	map[string]string	"Server error"
//	@Router			/notification [post]
func CreateNotification(c *fiber.Ctx) error {
	var notification models.Notification

	if err := c.BodyParser(&notification); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if err := db.Create(&notification).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(201).JSON(notification)
}

// GetNotifications godoc
//
//	@Summary		List all notifications
//	@Description	GetNotifications retrieves all Notification records with associated User
//	@Tags			Notifications
//	@Produce		json
//	@Success		200	{array}		models.NotificationDoc
//	@Failure		500	{object}	map[string]string	"Server error"
//	@Router			/notifications [get]
func GetNotifications(c *fiber.Ctx) error {
	var notifications []models.Notification
	if err := db.Preload("User").Find(&notifications).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}
	return c.Status(200).JSON(notifications)
}

func findNotification(id int, notification *models.Notification) error {
	return db.Preload("User").First(notification, "id = ?", id).Error
}

// GetNotification godoc
//
//	@Summary		Get notification by ID
//	@Description	GetNotification retrieves a single Notification by its ID, including the User
//	@Tags			Notifications
//	@Produce		json
//	@Param			id	path		int	true	"Notification ID"
//	@Success		200	{object}	models.NotificationDoc
//	@Failure		400	{object}	map[string]string	"Invalid ID"
//	@Failure		404	{object}	map[string]string	"Notification not found"
//	@Failure		500	{object}	map[string]string	"Server error"
//	@Router			/notification/{id} [get]
func GetNotification(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var notification models.Notification

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	err = findNotification(id, &notification)
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return c.Status(404).JSON("notification not found")
	case err != nil:
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(200).JSON(notification)
}

// UpdateNotification godoc
//
//	@Summary		Update an existing notification
//	@Description	UpdateNotification updates a Notification record by its ID
//	@Tags			Notifications
//	@Accept			json
//	@Produce		json
//	@Param			id				path		int						true	"Notification ID"
//	@Param			notification	body		models.NotificationDoc	true	"Updated notification payload"
//	@Success		200				{object}	models.NotificationDoc
//	@Failure		400				{object}	map[string]string	"Invalid input"
//	@Failure		404				{object}	map[string]string	"Notification not found"
//	@Failure		500				{object}	map[string]string	"Server error"
//	@Router			/notification/{id} [put]
func UpdateNotification(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var notification models.Notification

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	err = findNotification(id, &notification)
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return c.Status(404).JSON("notification not found")
	case err != nil:
		return c.Status(500).JSON(err.Error())
	}

	var notification_updated models.Notification
	if err := c.BodyParser(&notification_updated); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if err := db.Model(&notification).Updates(notification_updated).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(200).JSON(notification)
}

// DeleteNotification godoc
//
//	@Summary		Delete a notification by ID
//	@Description	DeleteNotification removes a Notification record by its ID
//	@Tags			Notifications
//	@Produce		json
//	@Param			id	path		int					true	"Notification ID"
//	@Success		200	{string}	string				"Successfully deleted notification"
//	@Failure		400	{object}	map[string]string	"Invalid ID"
//	@Failure		404	{object}	map[string]string	"Notification not found"
//	@Failure		500	{object}	map[string]string	"Server error"
//	@Router			/notification/{id} [delete]
func DeleteNotification(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var notification models.Notification

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	err = findNotification(id, &notification)
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return c.Status(404).JSON("notification not found")
	case err != nil:
		return c.Status(500).JSON(err.Error())
	}

	if err := db.Delete(&notification).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(200).JSON("Successfully deleted notification")
}
