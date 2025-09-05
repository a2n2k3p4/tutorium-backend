package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/a2n2k3p4/tutorium-backend/middlewares"
	"github.com/a2n2k3p4/tutorium-backend/models"
	"github.com/a2n2k3p4/tutorium-backend/storage"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ClassRoutes(app *fiber.App) {
	class := app.Group("/classes", middlewares.ProtectedMiddleware())
	class.Get("/filter", GetFilteredClasses)
	class.Get("/", GetClasses)
	class.Get("/:id", GetClass)

	classProtected := class.Group("/", middlewares.TeacherRequired())
	classProtected.Post("/", CreateClass)
	classProtected.Put("/:id", UpdateClass)
	classProtected.Delete("/:id", DeleteClass)
}

// CreateClass godoc
//
//		@Summary		Create a new class
//		@Description	CreateClass creates a new Class record
//		@Tags			Classes
//	 @Security 		BearerAuth
//		@Accept			json
//		@Produce		json
//		@Param			class	body		models.ClassDoc	true	"Class payload"
//		@Success		201		{object}	models.ClassDoc
//		@Failure		400		{object}	map[string]string	"Invalid input"
//		@Failure		500		{object}	map[string]string	"Server error"
//		@Router			/classes [post]
func CreateClass(c *fiber.Ctx) error {
	var class models.Class

	if err := c.BodyParser(&class); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if err := processBannerPicture(c, &class); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	db, err := middlewares.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	// begin transaction for create Class with ClassCategories
	tx := db.Begin()
	if tx.Error != nil {
		return c.Status(500).JSON(tx.Error.Error())
	}

	if err := tx.Omit("Categories").Create(&class).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(err.Error())
	}

	// Attach categories if provided
	if len(class.Categories) > 0 {
		var cats []models.ClassCategory
		names := make([]string, len(class.Categories))
		for i, cat := range class.Categories {
			names[i] = cat.ClassCategory
		}

		if err := tx.Where("class_category IN ?", names).Find(&cats).Error; err != nil {
			tx.Rollback()
			return c.Status(500).JSON(err.Error())
		}

		if len(cats) != len(names) {
			tx.Rollback()
			return c.Status(400).JSON(fiber.Map{"error": "some categories not found"})
		}

		if err := tx.Model(&class).Association("Categories").Replace(&cats); err != nil {
			tx.Rollback()
			return c.Status(500).JSON(err.Error())
		}
	}

	if err := tx.Commit().Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(201).JSON(class)
}

// GetClasses godoc
//
//		@Summary		List all classes
//		@Description	GetClasses retrieves all Class records with Teacher and Categories relations
//		@Tags			Classes
//	 @Security 		BearerAuth
//		@Produce		json
//		@Success		200	{array}		models.ClassDoc
//		@Failure		500	{object}	map[string]string	"Server error"
//		@Router			/classes [get]
func GetClasses(c *fiber.Ctx) error {
	classes := []models.Class{}
	db, err := middlewares.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	if err := db.Find(&classes).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	// Generate presigned URL if BannerPictureURL exists
	mc, ok := c.Locals("minio").(*storage.Client)
	if ok {
		for i := range classes {
			if classes[i].BannerPictureURL != "" {
				presignedURL, err := mc.PresignedGetObject(c.Context(), classes[i].BannerPictureURL, 15*time.Minute)
				if err == nil {
					classes[i].BannerPictureURL = presignedURL
				} else {
					classes[i].BannerPictureURL = ""
				}
			}
		}
	}

	return c.Status(200).JSON(classes)
}

func findClass(db *gorm.DB, id int, class *models.Class) error {
	return db.First(class, "id = ?", id).Error
}

// GetClass godoc
//
//		@Summary		Get class by ID
//		@Description	GetClass retrieves a single Class by its ID, including Teacher and Categories
//		@Tags			Classes
//	 @Security 		BearerAuth
//		@Produce		json
//		@Param			id	path		int	true	"Class ID"
//		@Success		200	{object}	models.ClassDoc
//		@Failure		400	{object}	map[string]string	"Invalid ID"
//		@Failure		404	{object}	map[string]string	"Class not found"
//		@Failure		500	{object}	map[string]string	"Server error"
//		@Router			/classes/{id} [get]
func GetClass(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var class models.Class

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}
	db, err := middlewares.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	err = db.Preload("Teacher").Preload("Categories").First(&class, "id = ?", id).Error
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return c.Status(404).JSON("class not found")
	case err != nil:
		return c.Status(500).JSON(err.Error())
	}

	// Generate presigned URL if BannerPictureURL exists
	mc, ok := c.Locals("minio").(*storage.Client)
	if ok && class.BannerPictureURL != "" {
		presignedURL, err := mc.PresignedGetObject(c.Context(), class.BannerPictureURL, 15*time.Minute)
		if err == nil {
			class.BannerPictureURL = presignedURL
		} else {
			class.BannerPictureURL = ""
		}
	}

	return c.Status(200).JSON(class)
}

// GetFilteredClasses godoc
//
//		@Summary		Get filtered classes
//		@Description	Retrieve a list of classes filtered by optional query parameters: categories, price range, and rating range
//		@Tags			Classes
//	 @Security 		BearerAuth
//		@Produce		json
//		@Param			category		query	[]string	false	"Filter by one or more categories (OR relation)"
//		@Param			min_price		query	string		false	"Minimum class price"
//		@Param			max_price		query	string		false	"Maximum class price"
//		@Param			min_rating		query	string		false	"Minimum class rating"
//		@Param			max_rating		query	string		false	"Maximum class rating"
//		@Success		200	{array}		models.ClassDoc
//		@Failure		400	{object}	map[string]string	"Invalid query parameters"
//		@Failure		500	{object}	map[string]string	"Server error"
//		@Router			/classes/filter [get]
func GetFilteredClasses(c *fiber.Ctx) error {
	db, err := middlewares.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	var filters struct {
		Categories []string `query:"category"`
		MinPrice   string   `query:"min_price"`
		MaxPrice   string   `query:"max_price"`
		MinRating  string   `query:"min_rating"`
		MaxRating  string   `query:"max_rating"`
	}

	if err := c.QueryParser(&filters); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	type FilteredClassResponse struct {
		ID               uint    `json:"id"`
		ClassName        string  `json:"class_name"`
		BannerPictureURL string  `json:"banner_picture_url"`
		Rating           float64 `json:"rating"`
		Price            float64 `json:"price"`
		TeacherName      string  `json:"teacher_name"`
	}
	var results []FilteredClassResponse

	// Get teacher's FirstName and LastName from users table
	query := db.Table("classes").
		Select(`
			classes.id,
			classes.class_name,
			classes.banner_picture_url,
			classes.rating,
			classes.price,
			users.first_name || ' ' || users.last_name AS teacher_name
		`).
		Joins("JOIN teachers ON teachers.id = classes.teacher_id").
		Joins("JOIN users ON users.id = teachers.user_id")

	// Categories filter
	if len(filters.Categories) > 0 {
		query = query.Where("classes.id IN (?)",
			db.Table("class_class_categories ccc").
				Joins("JOIN class_categories cc ON cc.id = ccc.class_category_id").
				Select("ccc.class_id").
				Where("cc.class_category IN ?", filters.Categories),
		)
	}

	// Price filter
	if filters.MinPrice != "" || filters.MaxPrice != "" {
		min, minErr := strconv.ParseFloat(filters.MinPrice, 64)
		max, maxErr := strconv.ParseFloat(filters.MaxPrice, 64)

		if minErr == nil && maxErr == nil {
			query = query.Where("price BETWEEN ? AND ?", min, max)
		} else if minErr == nil {
			query = query.Where("price >= ?", min)
		} else if maxErr == nil {
			query = query.Where("price <= ?", max)
		}
	}

	// Rating filter
	if filters.MinRating != "" || filters.MaxRating != "" {
		min, minErr := strconv.ParseFloat(filters.MinRating, 64)
		max, maxErr := strconv.ParseFloat(filters.MaxRating, 64)

		if minErr == nil && maxErr == nil {
			query = query.Where("rating BETWEEN ? AND ?", min, max)
		} else if minErr == nil {
			query = query.Where("rating >= ?", min)
		} else if maxErr == nil {
			query = query.Where("rating <= ?", max)
		}
	}

	// Executed query
	if err := query.Scan(&results).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	// Presigned URLs for BannerPicture
	mc, ok := c.Locals("minio").(*storage.Client)
	if ok {
		for i := range results {
			if results[i].BannerPictureURL != "" {
				presignedURL, err := mc.PresignedGetObject(c.Context(), results[i].BannerPictureURL, 15*time.Minute)
				if err == nil {
					results[i].BannerPictureURL = presignedURL
				} else {
					results[i].BannerPictureURL = ""
				}
			}
		}
	}

	return c.JSON(results)
}

// UpdateClass godoc
//
//		@Summary		Update an existing class
//		@Description	UpdateClass updates a Class record by its ID
//		@Tags			Classes
//	 @Security 		BearerAuth
//		@Accept			json
//		@Produce		json
//		@Param			id		path		int				true	"Class ID"
//		@Param			class	body		models.ClassDoc	true	"Updated class payload"
//		@Success		200		{object}	models.ClassDoc
//		@Failure		400		{object}	map[string]string	"Invalid input"
//		@Failure		404		{object}	map[string]string	"Class not found"
//		@Failure		500		{object}	map[string]string	"Server error"
//		@Router			/classes/{id} [put]
func UpdateClass(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var class models.Class

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}
	db, err := middlewares.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	err = findClass(db, id, &class)
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return c.Status(404).JSON("class not found")
	case err != nil:
		return c.Status(500).JSON(err.Error())
	}

	var class_update models.Class
	if err := c.BodyParser(&class_update); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if err := processBannerPicture(c, &class_update); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	// begin transaction for create Class with ClassCategories
	tx := db.Begin()
	if tx.Error != nil {
		return c.Status(500).JSON(tx.Error.Error())
	}

	if err := tx.Model(&class).
		Omit("Teacher").
		Omit("Categories").
		Updates(class_update).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(err.Error())
	}

	if len(class_update.Categories) > 0 {
		var cats []models.ClassCategory
		names := make([]string, len(class_update.Categories))
		for i, cat := range class_update.Categories {
			names[i] = cat.ClassCategory
		}

		if err := tx.Where("class_category IN ?", names).Find(&cats).Error; err != nil {
			tx.Rollback()
			return c.Status(500).JSON(err.Error())
		}

		if len(cats) != len(names) {
			tx.Rollback()
			return c.Status(400).JSON(fiber.Map{"error": "some categories not found"})
		}

		if err := tx.Model(&class).Association("Categories").Replace(&cats); err != nil {
			tx.Rollback()
			return c.Status(500).JSON(err.Error())
		}
	}

	if err := tx.Commit().Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	if err := db.First(&class, id).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(200).JSON(class)
}

// DeleteClass godoc
//
//		@Summary		Delete a class by ID
//		@Description	DeleteClass removes a Class record by its ID
//		@Tags			Classes
//	 @Security 		BearerAuth
//		@Produce		json
//		@Param			id	path		int					true	"Class ID"
//		@Success		200	{string}	string				"Successfully deleted class"
//		@Failure		400	{object}	map[string]string	"Invalid ID"
//		@Failure		404	{object}	map[string]string	"Class not found"
//		@Failure		500	{object}	map[string]string	"Server error"
//		@Router			/classes/{id} [delete]
func DeleteClass(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var class models.Class

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}
	db, err := middlewares.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	err = findClass(db, id, &class)
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return c.Status(404).JSON("class not found")
	case err != nil:
		return c.Status(500).JSON(err.Error())
	}

	if err = db.Delete(&class).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}
	return c.Status(200).JSON("Successfully deleted class")
}

func processBannerPicture(c *fiber.Ctx, class *models.Class) error {
	if class.BannerPictureURL != "" && !strings.HasPrefix(class.BannerPictureURL, "http") {
		b, err := storage.DecodeBase64Image(class.BannerPictureURL)
		if err != nil {
			return fmt.Errorf("invalid base64 image: %w", err)
		}
		if err := validateImageBytes(b); err != nil {
			return fmt.Errorf("invalid image: %w", err)
		}

		mc := c.Locals("minio").(*storage.Client)
		filename := storage.GenerateFilename(http.DetectContentType(b[:min(512, len(b))]))
		uploaded, err := mc.UploadBytes(context.Background(), "classes", filename, b)
		if err != nil {
			return err
		}
		class.BannerPictureURL = uploaded
	}
	return nil
}
