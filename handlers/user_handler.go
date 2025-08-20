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

// CreateUser godoc
// @Summary      Create a new user
// @Description  CreateUser creates a new user record
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        user  body      models.UserDoc  true  "User payload"
// @Success      201   {object}  models.UserDoc
// @Failure      400   {object}  map[string]string  "Invalid input"
// @Failure      500   {object}  map[string]string  "Server error"
// @Router       /user [post]
func CreateUser(c *fiber.Ctx) error {
	var user models.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if err := db.Create(&user).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(201).JSON(user)
}

// GetUsers godoc
// @Summary      List all users
// @Description  GetUsers retrieves all user records
// @Tags         Users
// @Produce      json
// @Success      200  {array}   models.UserDoc
// @Failure      500  {object}  map[string]string  "Server error"
// @Router       /users [get]
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

// GetUser godoc
// @Summary      Get user by ID
// @Description  GetUser retrieves a single user by its ID
// @Tags         Users
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  models.UserDoc
// @Failure      400  {object}  map[string]string  "Invalid ID"
// @Failure      404  {object}  map[string]string  "User not found"
// @Failure      500  {object}  map[string]string  "Server error"
// @Router       /user/{id} [get]
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

// UpdateUser godoc
// @Summary      Update an existing user
// @Description  UpdateUser updates a user record by its ID
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id    path      int         true  "User ID"
// @Param        user  body      models.UserDoc true  "Updated user payload"
// @Success      200   {object}  models.UserDoc
// @Failure      400   {object}  map[string]string  "Invalid input"
// @Failure      404   {object}  map[string]string  "User not found"
// @Failure      500   {object}  map[string]string  "Server error"
// @Router       /user/{id} [put]
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
		return c.Status(400).JSON(err.Error())
	}

	if err := db.Model(&user).Updates(user_update).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(200).JSON(user)
}

// DeleteUser godoc
// @Summary      Delete a user by ID
// @Description  DeleteUser removes a user record by its ID
// @Tags         Users
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {string}  string  "Successfully deleted User"
// @Failure      400  {object}  map[string]string  "Invalid ID"
// @Failure      404  {object}  map[string]string  "User not found"
// @Failure      500  {object}  map[string]string  "Server error"
// @Router       /user/{id} [delete]
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
