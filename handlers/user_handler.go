package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/a2n2k3p4/tutorium-backend/middlewares"
	"github.com/a2n2k3p4/tutorium-backend/models"
	"github.com/a2n2k3p4/tutorium-backend/storage"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func UserRoutes(app *fiber.App) {
	user := app.Group("/users", middlewares.ProtectedMiddleware())
	user.Post("/", CreateUser)
	user.Get("/:id", GetUser)
	user.Put("/:id", UpdateUser)
	user.Delete("/:id", DeleteUser)

	userAdmin := user.Group("/", middlewares.AdminRequired())
	userAdmin.Get("/", GetUsers)
}

// CreateUser godoc
//
//	@Summary		Create a new user
//	@Description	CreateUser creates a new user record
//	@Tags			Users
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			user	body		models.UserDoc	true	"User payload"
//	@Success		201		{object}	models.UserDoc
//	@Failure		400		{string}	string	"Invalid input"
//	@Failure		500		{string}	string	"Server error"
//	@Router			/users [post]
func CreateUser(c *fiber.Ctx) error {
	var user models.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if err := processProfilePicture(c, &user); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	db, err := middlewares.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	// begin transaction for create User and Learner
	tx := db.Begin()
	if tx.Error != nil {
		return c.Status(500).JSON(tx.Error.Error())
	}

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(err.Error())
	}

	learner := models.Learner{
		UserID:    user.ID,
		FlagCount: 0,
	}
	if err := tx.Create(&learner).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	// attach to user struct so response contains learner
	user.Learner = &learner
	return c.Status(201).JSON(user)
}

// GetUsers godoc
//
//	@Summary		List all users
//	@Description	GetUsers retrieves all user records
//	@Tags			Users
//	@Security		BearerAuth
//	@Produce		json
//	@Success		200	{array}		models.UserDoc
//	@Failure		500	{string}	string	"Server error"
//	@Router			/users [get]
func GetUsers(c *fiber.Ctx) error {
	users := []models.User{}
	db, err := middlewares.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	if err := db.Preload("Learner").
		Preload("Teacher").
		Preload("Admin").
		Find(&users).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	// Generate presigned URL if ProfilePictureURL exists for each user
	mc, ok := c.Locals("minio").(*storage.Client)
	if ok {
		for i := range users {
			if users[i].ProfilePictureURL != "" {
				presignedURL, err := mc.PresignedGetObject(c.Context(), users[i].ProfilePictureURL, 15*time.Minute)
				if err == nil {
					users[i].ProfilePictureURL = presignedURL
				} else {
					users[i].ProfilePictureURL = ""
				}
			}
		}
	}

	return c.Status(200).JSON(users)
}

func findUser(db *gorm.DB, id int, user *models.User) error {
	return db.Preload("Learner").
		Preload("Teacher").
		Preload("Admin").
		First(user, "id = ?", id).Error
}

// GetUser godoc
//
//	@Summary		Get user by ID
//	@Description	GetUser retrieves a single user by its ID
//	@Tags			Users
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	models.UserDoc
//	@Failure		400	{string}	string	"Invalid ID"
//	@Failure		404	{string}	string	"User not found"
//	@Failure		500	{string}	string	"Server error"
//	@Router			/users/{id} [get]
func GetUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var user models.User

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}
	db, err := middlewares.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	err = findUser(db, id, &user)
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return c.Status(404).JSON("user not found")
	case err != nil:
		return c.Status(500).JSON(err.Error())
	}

	// Generate presigned URL if ProfilePictureURL exists
	mc, ok := c.Locals("minio").(*storage.Client)
	if ok && user.ProfilePictureURL != "" {
		presignedURL, err := mc.PresignedGetObject(c.Context(), user.ProfilePictureURL, 15*time.Minute)
		if err == nil {
			user.ProfilePictureURL = presignedURL
		} else {
			// optional: log the error
			user.ProfilePictureURL = ""
		}
	}

	return c.Status(200).JSON(user)
}

// UpdateUser godoc
//
//	@Summary		Update an existing user
//	@Description	UpdateUser updates a user record by its ID
//	@Tags			Users
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int				true	"User ID"
//	@Param			user	body		models.UserDoc	true	"Updated user payload"
//	@Success		200		{object}	models.UserDoc
//	@Failure		400		{string}	string	"Invalid input"
//	@Failure		404		{string}	string	"User not found"
//	@Failure		500		{string}	string	"Server error"
//	@Router			/users/{id} [put]
func UpdateUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var user models.User

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}
	db, err := middlewares.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	err = findUser(db, id, &user)
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

	if err := processProfilePicture(c, &user_update); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if err := db.Model(&user).
		Omit(clause.Associations).
		Updates(user_update).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(200).JSON(user)
}

// DeleteUser godoc
//
//	@Summary		Delete a user by ID
//	@Description	DeleteUser removes a user record by its ID along with associated Learner, Teacher, and Admin
//	@Tags			Users
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id	path		int		true	"User ID"
//	@Success		200	{string}	string	"Successfully deleted User and associated roles"
//	@Failure		400	{string}	string	"Invalid ID"
//	@Failure		404	{string}	string	"User not found"
//	@Failure		500	{string}	string	"Server error"
//	@Router			/users/{id} [delete]
func DeleteUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var user models.User

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}
	db, err := middlewares.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	err = findUser(db, id, &user)
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return c.Status(404).JSON("user not found")
	case err != nil:
		return c.Status(500).JSON(err.Error())
	}

	if err := db.Select(clause.Associations).Delete(&user).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}
	return c.Status(200).JSON("Successfully deleted User")
}

func processProfilePicture(c *fiber.Ctx, user *models.User) error {
	if user.ProfilePictureURL != "" && !strings.HasPrefix(user.ProfilePictureURL, "http") {
		b, err := storage.DecodeBase64Image(user.ProfilePictureURL)
		if err != nil {
			return fmt.Errorf("invalid base64 image: %w", err)
		}
		if err := validateImageBytes(b); err != nil {
			return fmt.Errorf("invalid image: %w", err)
		}
		up := c.Locals("minio").(storage.Uploader)
		filename := storage.GenerateFilename(http.DetectContentType(b[:min(512, len(b))]))
		objectKey, err := up.UploadBytes(c.Context(), "users", filename, b)
		if err != nil {
			return err
		}
		user.ProfilePictureURL = objectKey
	}
	return nil
}
