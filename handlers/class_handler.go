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
	"gorm.io/gorm/clause"
)

func ClassRoutes(app *fiber.App) {
	class := app.Group("/classes", middlewares.ProtectedMiddleware())
	class.Get("/", GetClasses)
	class.Get("/detailed", FindClassesDetailed)
	class.Get("/:id", GetClass)

	classProtected := class.Group("/", middlewares.TeacherRequired())
	classProtected.Post("/", CreateClass)
	classProtected.Put("/:id", UpdateClass)
	classProtected.Delete("/:id", DeleteClass)
	classProtected.Get("/:id/class_categories", GetClassCategoriesByClassID)
	classProtected.Post("/:id/class_categories", AddClassCategories)
	classProtected.Delete("/:id/class_categories", DeleteClassCategories)
}

// CreateClass godoc
//
//	@Summary		Create a new class
//	@Description	CreateClass creates a new Class record
//	@Tags			Classes
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			class	body		models.ClassDoc	true	"Class payload"
//	@Success		201		{object}	models.ClassDoc
//	@Failure		400		{string}	string	"Invalid input"
//	@Failure		500		{string}	string	"Server error"
//	@Router			/classes [post]
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
			return c.Status(400).JSON("some categories not found")
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
//	@Summary		List all classes
//	@Description	Retrieve a list of classes filtered by optional query parameters: categories, and rating range
//	@Tags			Classes
//	@Security		BearerAuth
//	@Produce		json
//	@Param			category		query	[]string	false	"Filter by one or more categories (OR relation)"
//	@Param			min_rating		query	string		false	"Minimum class rating"
//	@Param			max_rating		query	string		false	"Maximum class rating"
//	@Success		200	{array}		models.ClassDoc
//	@Failure		400	{string}	string	"Invalid query parameters"
//	@Failure		500	{string}	string	"Server error"
//	@Router			/classes [get]
func GetClasses(c *fiber.Ctx) error {
	db, err := middlewares.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	var filters struct {
		Categories []string `query:"category"`
		MinRating  string   `query:"min_rating"`
		MaxRating  string   `query:"max_rating"`
	}

	if err := c.QueryParser(&filters); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if len(filters.Categories) == 1 && strings.Contains(filters.Categories[0], ",") {
		filters.Categories = strings.Split(filters.Categories[0], ",")
		for i := range filters.Categories {
			filters.Categories[i] = strings.TrimSpace(filters.Categories[i])
		}
	}

	type ClassResponse struct {
		ID               uint    `json:"id"`
		ClassName        string  `json:"class_name"`
		BannerPictureURL string  `json:"banner_picture_url"`
		Rating           float64 `json:"rating"`
		TeacherName      string  `json:"teacher_name"`
	}
	var results []ClassResponse

	// Get teacher's FirstName and LastName from users table
	query := db.Table("classes").
		Select(`
			classes.id,
			classes.class_name,
			classes.banner_picture_url,
			classes.rating,
			CONCAT(users.first_name, ' ', users.last_name) AS teacher_name
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

	return c.Status(200).JSON(results)
}

func findClass(db *gorm.DB, id int, class *models.Class) error {
	return db.First(class, "id = ?", id).Error
}

// GetClass godoc
//
//	@Summary		Get class by ID
//	@Description	GetClass retrieves a single Class by its ID, including Teacher and Categories
//	@Tags			Classes
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id	path		int	true	"Class ID"
//	@Success		200	{object}	models.ClassDoc
//	@Failure		400	{string}	string	"Invalid ID"
//	@Failure		404	{string}	string	"Class not found"
//	@Failure		500	{string}	string	"Server error"
//	@Router			/classes/{id} [get]
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

// UpdateClass godoc
//
//	@Summary		Update an existing class
//	@Description	UpdateClass updates a Class record by its ID
//	@Tags			Classes
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int				true	"Class ID"
//	@Param			class	body		models.ClassDoc	true	"Updated class payload"
//	@Success		200		{object}	models.ClassDoc
//	@Failure		400		{string}	string	"Invalid input"
//	@Failure		404		{string}	string	"Class not found"
//	@Failure		500		{string}	string	"Server error"
//	@Router			/classes/{id} [put]
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
		Omit(clause.Associations).
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
			return c.Status(400).JSON("some categories not found")
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
//	@Summary		Delete a class by ID
//	@Description	DeleteClass removes a Class record by its ID
//	@Tags			Classes
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id	path		int					true	"Class ID"
//	@Success		200	{string}	string				"Successfully deleted class"
//	@Failure		400	{string}	string	"Invalid ID"
//	@Failure		404	{string}	string	"Class not found"
//	@Failure		500	{string}	string	"Server error"
//	@Router			/classes/{id} [delete]
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

		mc := c.Locals("minio").(storage.Uploader)
		filename := storage.GenerateFilename(http.DetectContentType(b[:min(512, len(b))]))
		uploaded, err := mc.UploadBytes(context.Background(), "classes", filename, b)
		if err != nil {
			return err
		}
		class.BannerPictureURL = uploaded
	}
	return nil
}

func FindClassesDetailed(c *fiber.Ctx) error {
	// Query detailed class info (Teacher, Categories) by ID or category list (or return all if neither)
	// Example 1: /classes/detailed?id=1
	// Example 2: /classes/detailed?categories=Yoga,Math
	// Example 3: /classes/detailed

	db, err := middlewares.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// Check ID parameter
	if idParam := c.Query("id"); idParam != "" {
		id, err := strconv.Atoi(idParam)
		if err != nil || id <= 0 {
			return c.Status(400).JSON(fiber.Map{"error": "invalid id parameter"})
		}

		var class models.Class
		err = db.Preload("Categories").Preload("Teacher").
			First(&class, "id = ?", id).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(404).JSON("class not found")
		} else if err != nil {
			return c.Status(500).JSON(err.Error())
		}

		// Generate presigned URL if available
		if mc, ok := c.Locals("minio").(*storage.Client); ok && class.BannerPictureURL != "" {
			url, err := mc.PresignedGetObject(c.Context(), class.BannerPictureURL, 15*time.Minute)
			if err == nil {
				class.BannerPictureURL = url
			} else {
				class.BannerPictureURL = ""
			}
		}

		return c.Status(200).JSON(class)
	}

	// Else, category-based query
	rawCategories := c.Query("categories")
	parts := strings.Split(rawCategories, ",")
	categories := make([]string, 0, len(parts))
	for _, p := range parts {
		t := strings.TrimSpace(p)
		if t != "" {
			categories = append(categories, t)
		}
	}

	var classes []models.Class
	query := db.Model(&models.Class{}).
		Preload("Categories").
		Preload("Teacher")

	if len(categories) > 0 {
		query = query.
			Joins("JOIN class_class_categories ccc ON ccc.class_id = classes.id").
			Joins("JOIN class_categories cc ON cc.id = ccc.class_category_id").
			Where("cc.class_category IN ?", categories).
			Distinct("classes.*")
	}

	if err := query.Find(&classes).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// Generate presigned URLs for all
	if mc, ok := c.Locals("minio").(*storage.Client); ok {
		for i := range classes {
			if classes[i].BannerPictureURL != "" {
				url, err := mc.PresignedGetObject(c.Context(), classes[i].BannerPictureURL, 15*time.Minute)
				if err == nil {
					classes[i].BannerPictureURL = url
				} else {
					classes[i].BannerPictureURL = ""
				}
			}
		}
	}

	return c.Status(200).JSON(classes)
}

func GetClassCategoriesByClassID(c *fiber.Ctx) error {
	// return class categories names given class ID
	// input : class ID
	// output: list of class categories names
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).SendString("invalid :id")
	}

	db, err := middlewares.GetDB(c)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	var class models.Class
	if err := db.
		Preload("Categories", func(tx *gorm.DB) *gorm.DB {
			return tx.Select("id", "class_category").Order("class_categories.class_category")
		}).
		First(&class, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(404).SendString("class not found")
		}
		return c.Status(500).SendString(err.Error())
	}

	names := make([]string, 0, len(class.Categories))
	for _, cat := range class.Categories {
		names = append(names, cat.ClassCategory)
	}
	return c.Status(200).JSON(fiber.Map{"categories": names})
}

func AddClassCategories(c *fiber.Ctx) error {
	// input: class ID from URL, classCategory IDs from body
	// output: updated class with categories

	type payload struct {
		CategoryIDs []int `json:"class_category_ids"`
	}
	//example JSON body: { "class_category_ids": [1,2,3] }

	// Get :id from route
	classID, err := c.ParamsInt("id")
	if err != nil || classID <= 0 {
		return c.Status(400).JSON("invalid class ID")
	}

	// Parse category IDs from body
	var p payload
	if err := c.BodyParser(&p); err != nil {
		return c.Status(400).JSON("invalid body")
	}
	if len(p.CategoryIDs) == 0 {
		return c.Status(400).JSON("no class category IDs provided")
	}

	db, err := middlewares.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	// Validate class exists
	var class models.Class
	err = db.First(&class, classID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(404).JSON("class not found")
		}
		return c.Status(500).JSON(err.Error())
	}

	// Deduplicate and validate category IDs
	seen := make(map[int]bool)
	var validIDs []int
	for _, id := range p.CategoryIDs {
		if id > 0 && !seen[id] {
			seen[id] = true
			validIDs = append(validIDs, id)
		}
	}
	if len(validIDs) == 0 {
		return c.Status(400).JSON("no valid class category IDs")
	}

	// Get existing categories and avoid duplicates
	var existingCategories []models.ClassCategory
	err = db.Model(&class).Association("Categories").Find(&existingCategories)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	existingSet := make(map[uint]bool)
	for _, cat := range existingCategories {
		existingSet[cat.ID] = true
	}

	// Find and add only new categories
	var categoriesToAdd []models.ClassCategory
	err = db.Where("id IN ?", validIDs).Find(&categoriesToAdd).Error
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	var newCategories []models.ClassCategory
	for _, cat := range categoriesToAdd {
		if !existingSet[cat.ID] {
			newCategories = append(newCategories, cat)
		}
	}

	if len(newCategories) > 0 {
		err = db.Model(&class).Association("Categories").Append(newCategories)
		if err != nil {
			return c.Status(500).JSON(err.Error())
		}
	}

	// Return updated class with categories
	var result models.Class
	err = db.Preload("Categories").First(&result, classID).Error
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(200).JSON(result)
}

func DeleteClassCategories(c *fiber.Ctx) error {
	// input: class ID from URL, classCategory IDs from body
	// output: updated class with deleted categories
	type payload struct {
		CategoryIDs []int `json:"class_category_ids"`
	}
	//example JSON body: { "class_category_ids": [1,2,3] }
	// Get :id from route
	classID, err := c.ParamsInt("id")
	if err != nil || classID <= 0 {
		return c.Status(400).JSON("invalid class ID")
	}

	// Parse category IDs from body
	var p payload
	if err := c.BodyParser(&p); err != nil {
		return c.Status(400).JSON("invalid body")
	}
	if len(p.CategoryIDs) == 0 {
		return c.Status(400).JSON("no class category IDs provided")
	}

	db, err := middlewares.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	// Validate class exists
	var class models.Class
	if err := db.First(&class, classID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(404).JSON("class not found")
		}
		return c.Status(500).JSON(err.Error())
	}

	// Deduplicate IDs
	seen := make(map[int]bool)
	var validIDs []int
	for _, id := range p.CategoryIDs {
		if id > 0 && !seen[id] {
			seen[id] = true
			validIDs = append(validIDs, id)
		}
	}
	if len(validIDs) == 0 {
		return c.Status(400).JSON("no valid class category IDs")
	}

	// Load categories to delete
	var categories []models.ClassCategory
	if err := db.Where("id IN ?", validIDs).Find(&categories).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}
	if len(categories) == 0 {
		return c.Status(400).JSON("no matching class categories found")
	}

	// Delete the association links (not the categories themselves)
	if err := db.Model(&class).Association("Categories").Delete(categories); err != nil {
		return c.Status(500).JSON(err.Error())
	}

	// Return updated class with categories
	var result models.Class
	if err := db.Preload("Categories").First(&result, classID).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(200).JSON(result)
}
