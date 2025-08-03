package handlers

import (
	"errors"

	"github.com/a2n2k3p4/tutorium-backend/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func UserRoutes(app *fiber.App) {
	app.Post("/user", CreateUser)
	app.Get("/users", GetUsers)
	app.Get("/user/:id", GetUser)
	app.Put("/user/:id", UpdateUser)
	app.Delete("/user/:id", DeleteUser)
}

func CreateUser(c *fiber.Ctx) error {
	var user models.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if err := db.Create(&user).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(200).JSON(user)
}

func GetUsers(c *fiber.Ctx) error {
	users := []models.User{}
	if err := db.Find(&users).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(200).JSON(users)
}

func findUser(id int, user *models.User) error {
	return db.First(user, "id = ?", id).Error
}

func GetUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var user models.User

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	err = findUser(id, &user)
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return c.Status(404).JSON("user not found")
	case err != nil:
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(200).JSON(user)
}

func UpdateUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var user models.User

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	err = findUser(id, &user)
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return c.Status(404).JSON("user not found")
	case err != nil:
		return c.Status(500).JSON(err.Error())
	}

	var user_update models.User
	if err := c.BodyParser(&user_update); err != nil {
		return c.Status(500).JSON(err.Error())
	}

	if err := db.Model(&user).Updates(user_update).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(200).JSON(user)
}

func DeleteUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var user models.User

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	err = findUser(id, &user)
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return c.Status(404).JSON("user not found")
	case err != nil:
		return c.Status(500).JSON(err.Error())
	}

	if err = db.Delete(&user).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}
	return c.Status(200).JSON("Successfully deleted User")
}
