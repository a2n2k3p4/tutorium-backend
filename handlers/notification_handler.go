package handlers

import (
	"errors"

	"github.com/a2n2k3p4/tutorium-backend/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func NotificationRoutes(app *fiber.App) {
	app.Post("/notification", CreateNotification)
	app.Get("/notifications", GetNotifications)
	app.Get("/notification/:id", GetNotification)
	app.Put("/notification/:id", UpdateNotification)
	app.Delete("/notification/:id", DeleteNotification)
}

func CreateNotification(c *fiber.Ctx) error {
	var notification models.Notification

	if err := c.BodyParser(&notification); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if err := db.Create(&notification).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(200).JSON(notification)
}

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
