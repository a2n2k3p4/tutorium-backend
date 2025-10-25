package handlers

import (
	"errors"
	"time"

	"github.com/a2n2k3p4/tutorium-backend/middlewares"
	"github.com/a2n2k3p4/tutorium-backend/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func LearnerRoutes(app *fiber.App) {
	learner := app.Group("/learners", middlewares.ProtectedMiddleware(), middlewares.BanMiddleware())

	learner.Post("/", CreateLearner)
	learner.Get("/", GetLearners)
	learner.Get("/:id", GetLearner)
	// learner.Put("/:id", UpdateLearner) No application logic for updating learner
	learner.Delete("/:id", DeleteLearner)
	learner.Get("/:id/recommended", RecommendClasses)

	learner.Get("/:id/interests", GetClassInterestsByLearnerID)
	learner.Post("/:id/interests", AddLearnerInterests)
	learner.Delete("/:id/interests", DeleteLearnerInterests)

}

// CreateLearner godoc
//
//	@Summary		Create a new learner
//	@Description	CreateLearner creates a new Learner record
//	@Tags			Learners
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			learner	body		models.LearnerDoc	true	"Learner payload"
//	@Success		201		{object}	models.LearnerDoc
//	@Failure		400		{string}	string	"Invalid input"
//	@Failure		500		{string}	string	"Server error"
//	@Router			/learners [post]
func CreateLearner(c *fiber.Ctx) error {
	var learner models.Learner

	if err := c.BodyParser(&learner); err != nil {
		return c.Status(400).JSON(err.Error())
	}
	db, err := middlewares.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	if err := db.Create(&learner).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(201).JSON(learner)
}

// GetLearners godoc
//
//	@Summary		List all learners
//	@Description	GetLearners retrieves all Learner records
//	@Tags			Learners
//	@Security		BearerAuth
//	@Produce		json
//	@Success		200	{array}		models.LearnerDoc
//	@Failure		500	{string}	string	"Server error"
//	@Router			/learners [get]
func GetLearners(c *fiber.Ctx) error {
	//input : learner ID from url
	//output : learner data, including interested categories
	learners := []models.Learner{}
	db, err := middlewares.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}
	// Preload Interested categories via the `many2many:interested_class_categories` link
	if err := db.Preload("Interested", func(tx *gorm.DB) *gorm.DB {
		return tx.
			Select("id", "class_category").
			Order("class_categories.class_category")
	}).Find(&learners).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(200).JSON(learners)
}

func findLearner(db *gorm.DB, id int, learner *models.Learner) error {
	//helper for GetLearner and DeleteLearner
	//Preload Interested class categories
	return db.
		Preload("Interested", func(tx *gorm.DB) *gorm.DB {
			return tx.
				Select("id", "class_category").
				Order("class_categories.class_category")
		}).First(learner, "id = ?", id).Error
}

// GetLearner godoc
//
//	@Summary		Get learner by ID
//	@Description	GetLearner retrieves a single Learner by its ID
//	@Tags			Learners
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id	path		int	true	"Learner ID"
//	@Success		200	{object}	models.LearnerDoc
//	@Failure		400	{string}	string	"Invalid ID"
//	@Failure		404	{string}	string	"Learner not found"
//	@Failure		500	{string}	string	"Server error"
//	@Router			/learners/{id} [get]
func GetLearner(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var learner models.Learner

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}
	db, err := middlewares.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	err = findLearner(db, id, &learner)
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return c.Status(404).JSON("Learner not found")
	case err != nil:
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(200).JSON(learner)
}

// DeleteLearner godoc
//
//	@Summary		Delete a learner by ID
//	@Description	DeleteLearner removes a Learner record by its ID
//	@Tags			Learners
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id	path		int		true	"Learner ID"
//	@Success		200	{string}	string	"Successfully deleted Learner"
//	@Failure		400	{string}	string	"Invalid ID"
//	@Failure		404	{string}	string	"Learner not found"
//	@Failure		500	{string}	string	"Server error"
//	@Router			/learners/{id} [delete]
func DeleteLearner(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var learner models.Learner

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	db, err := middlewares.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	err = findLearner(db, id, &learner)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return c.Status(404).JSON("Learner not found")
		default:
			return c.Status(500).JSON(err.Error())
		}
	}
	// Clear many2many association with class categories
	if err := db.Model(&learner).Association("Interested").Clear(); err != nil {
		return c.Status(500).JSON(err.Error())
	}

	if err = db.Delete(&learner).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}
	return c.Status(200).JSON("Successfully deleted Learner")
}

// RecommendClasses godoc
//
//	@Summary		Recommend enrollable classes for a learner
//	@Description	Returns two buckets: **recommended_classes** (enrollable classes matching the learner’s interested categories) and **remaining_classes** (other enrollable classes). If no interests or no matches, `recommended_found` is false and `remaining_classes` contains all enrollable classes.
//	@Tags			Learners
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id	path		int		true	"Learner ID"
//	@Success		200	{object}	models.RecommendClassesDoc
//	@Failure		400	{string}	string	"Please ensure that :id is an integer"
//	@Failure		404	{string}	string	"learner not found"
//	@Failure		500	{string}	string	"Server error"
//	@Router			/learners/{id}/recommended [get]
func RecommendClasses(c *fiber.Ctx) error {
	/*
		input : learner ID from url
		output :
			{
			    "recommended_found": true|false,
			    "recommended_classes": [...],  // enrollable & match interests
			    "remaining_classes":   [...]   // enrollable & outside interests
			  }
	*/
	// Parse :id from url
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	db, err := middlewares.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	// Load learner and interests
	var lr models.Learner
	if err := db.Preload("Interested").First(&lr, "id = ?", id).Error; err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return c.Status(404).JSON("learner not found")
		default:
			return c.Status(500).JSON(err.Error())
		}
	}

	now := time.Now()

	// Return all classes that have at least one enrollable session
	fetchAllActive := func() ([]models.Class, error) {
		var allActive []models.Class
		err := db.
			Model(&models.Class{}).
			Joins("JOIN class_sessions cs ON cs.class_id = classes.id AND cs.enrollment_deadline > ?", now).
			Select("classes.*").
			Group("classes.id").
			Preload("Teacher").
			Preload("Categories").
			Find(&allActive).Error
		return allActive, err
	}

	type response struct {
		RecommendedFound   bool           `json:"recommended_found"`
		RecommendedClasses []models.Class `json:"recommended_classes"`
		RemainingClasses   []models.Class `json:"remaining_classes"`
	}

	// No interests → recommended empty, remaining = all active
	if len(lr.Interested) == 0 {
		allActive, err := fetchAllActive()
		if err != nil {
			return c.Status(500).JSON(err.Error())
		}
		return c.Status(200).JSON(response{
			RecommendedFound:   false,
			RecommendedClasses: []models.Class{},
			RemainingClasses:   allActive,
		})
	}

	// Collect interested category IDs
	catIDs := make([]uint, 0, len(lr.Interested))
	for _, cat := range lr.Interested {
		catIDs = append(catIDs, cat.ID)
	}

	// 1) Recommended (active + match any interest)
	var recommended []models.Class
	recQ := db.Model(&models.Class{}).
		Joins("JOIN class_sessions cs ON cs.class_id = classes.id AND cs.enrollment_deadline > ?", now).
		Joins("JOIN class_class_categories ccc ON ccc.class_id = classes.id").
		Where("ccc.class_category_id IN ?", catIDs).
		Select("classes.*").
		Group("classes.id").
		Order("COUNT(ccc.class_category_id) DESC").
		Preload("Teacher").
		Preload("Categories")
	if err := recQ.Find(&recommended).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	// If no recommended → remaining = all active
	if len(recommended) == 0 {
		allActive, err := fetchAllActive()
		if err != nil {
			return c.Status(500).JSON(err.Error())
		}
		return c.Status(200).JSON(response{
			RecommendedFound:   false,
			RecommendedClasses: []models.Class{},
			RemainingClasses:   allActive,
		})
	}

	// 2) Remaining (active BUT not in recommended)
	recIDs := make([]uint, 0, len(recommended))
	for _, cls := range recommended {
		recIDs = append(recIDs, cls.ID)
	}

	var remaining []models.Class
	remQ := db.Model(&models.Class{}).
		Joins("JOIN class_sessions cs ON cs.class_id = classes.id AND cs.enrollment_deadline > ?", now).
		Where("classes.id NOT IN ?", recIDs).
		Select("classes.*").
		Group("classes.id").
		Preload("Teacher").
		Preload("Categories")
	if err := remQ.Find(&remaining).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(200).JSON(response{
		RecommendedFound:   true,
		RecommendedClasses: recommended,
		RemainingClasses:   remaining,
	})
}

// AddLearnerInterests godoc
//
//	@Summary		Add categories to a learner's interests
//	@Description	Appends rows to `interested_class_categories` for the given learner. Duplicate IDs are ignored. Returns the updated Learner with `Interested` preloaded.
//	@Tags			Learners
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int							true	"Learner ID"
//	@Param			body	body		models.ClassCategoryIDsDoc	true	"IDs of class categories to add"
//	@Success		200		{object}	models.LearnerDoc
//	@Failure		400		{string}	string	"invalid learner ID"
//	@Failure		400		{string}	string	"invalid body"
//	@Failure		400		{string}	string	"no class category IDs provided"
//	@Failure		400		{string}	string	"no valid class category IDs"
//	@Failure		404		{string}	string	"learner not found"
//	@Failure		500		{string}	string	"Server error"
//	@Router			/learners/{id}/interests [post]
func AddLearnerInterests(c *fiber.Ctx) error {
	/*
		input : learner ID from url, body with class_category_ids array
		output : updated Learner with Interested categories
		example body :{ "class_category_ids": [1,2,3] }
	*/
	type payload struct {
		CategoryIDs []int `json:"class_category_ids"`
	}

	// Get :id from route
	learnerID, err := c.ParamsInt("id")
	if err != nil || learnerID <= 0 {
		return c.Status(400).JSON("invalid learner ID")
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

	// Validate learner exists
	var learner models.Learner
	err = db.First(&learner, learnerID).Error
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return c.Status(404).JSON("learner not found")
	case err != nil:
		return c.Status(500).JSON(err.Error())
	}

	// Deduplicate and validate IDs
	seen := make(map[int]bool)
	validIDs := make([]int, 0, len(p.CategoryIDs))
	for _, id := range p.CategoryIDs {
		if id > 0 && !seen[id] {
			seen[id] = true
			validIDs = append(validIDs, id)
		}
	}
	if len(validIDs) == 0 {
		return c.Status(400).JSON("no valid class category IDs")
	}

	// Get existing interests to avoid duplicates
	var existing []models.ClassCategory
	if err := db.Model(&learner).Association("Interested").Find(&existing); err != nil {
		return c.Status(500).JSON(err.Error())
	}
	existingSet := make(map[uint]bool, len(existing))
	for _, cat := range existing {
		existingSet[cat.ID] = true
	}

	// Load categories by IDs
	var categories []models.ClassCategory
	if err := db.Where("id IN ?", validIDs).Find(&categories).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	// Filter only new categories
	newCategories := make([]models.ClassCategory, 0, len(categories))
	for _, cat := range categories {
		if !existingSet[cat.ID] {
			newCategories = append(newCategories, cat)
		}
	}

	// Append new interests (if any)
	if len(newCategories) > 0 {
		if err := db.Model(&learner).Association("Interested").Append(newCategories); err != nil {
			return c.Status(500).JSON(err.Error())
		}
	}

	// Return updated learner with Interested preloaded
	var result models.Learner
	if err := db.Preload("Interested").First(&result, learnerID).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}
	return c.Status(200).JSON(result)
}

// DeleteLearnerInterests godoc
//
//	@Summary		Remove categories from a learner's interests
//	@Description	Deletes the association rows in `interested_class_categories` for the given learner and category IDs. Returns the updated Learner with `Interested` preloaded and ordered.
//	@Tags			Learners
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int						true	"Learner ID"
//	@Param			body	body		models.ClassCategoryIDsDoc	true	"IDs of class categories to remove"
//	@Success		200		{object}	models.LearnerDoc
//	@Failure		400		{string}	string	"invalid learner ID"
//	@Failure		400		{string}	string	"invalid body"
//	@Failure		400		{string}	string	"no class category IDs provided"
//	@Failure		400		{string}	string	"no valid class category IDs"
//	@Failure		400		{string}	string	"no matching class categories found"
//	@Failure		404		{string}	string	"learner not found"
//	@Failure		500		{string}	string	"Server error"
//	@Router			/learners/{id}/interests [delete]
func DeleteLearnerInterests(c *fiber.Ctx) error {
	/*
		input : learner ID from url, body with class_category_ids array
		output : updated Learner with deleted Interested categories
		example body :{ "class_category_ids": [1,2,3] }
	*/
	type payload struct {
		CategoryIDs []int `json:"class_category_ids"`
	}

	learnerID, err := c.ParamsInt("id")
	if err != nil || learnerID <= 0 {
		return c.Status(400).JSON("invalid learner ID")
	}

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

	// Validate learner exists
	var learner models.Learner
	if err := db.First(&learner, learnerID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(404).JSON("learner not found")
		}
		return c.Status(500).JSON(err.Error())
	}

	// Deduplicate & validate IDs
	seen := make(map[int]bool, len(p.CategoryIDs))
	validIDs := make([]int, 0, len(p.CategoryIDs))
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
	if err := db.Model(&learner).Association("Interested").Delete(categories); err != nil {
		return c.Status(500).JSON(err.Error())
	}

	// Return updated learner with Interested preloaded (ordered)
	var result models.Learner
	if err := db.
		Preload("Interested", func(tx *gorm.DB) *gorm.DB {
			return tx.Select("id", "class_category").Order("class_categories.class_category")
		}).
		First(&result, learnerID).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(200).JSON(result)
}

// GetClassInterestsByLearnerID godoc
//
//	@Summary		Get a learner's interested class categories
//	@Description	Returns the list of class category names the learner is interested in (alphabetically ordered).
//	@Tags			Learners
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id	path		int		true	"Learner ID"
//	@Success		200	{object}	map[string][]string	"keys: categories"
//	@Failure		400	{string}	string	"invalid :id"
//	@Failure		404	{string}	string	"learner not found"
//	@Failure		500	{string}	string	"Server error"
//	@Router			/learners/{id}/interests [get]
func GetClassInterestsByLearnerID(c *fiber.Ctx) error {
	/*
		input : learner ID from url
		output : array of class category names the learner is interested in
	*/
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).SendString("invalid :id")
	}

	db, err := middlewares.GetDB(c)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	var learner models.Learner
	if err := db.
		Preload("Interested", func(tx *gorm.DB) *gorm.DB {
			// select only needed fields & order by name
			return tx.Select("id", "class_category").Order("class_categories.class_category")
		}).
		First(&learner, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(404).SendString("learner not found")
		}
		return c.Status(500).SendString(err.Error())
	}

	names := make([]string, 0, len(learner.Interested))
	for _, cat := range learner.Interested {
		names = append(names, cat.ClassCategory)
	}
	// Keep the same response shape as your class example
	return c.Status(200).JSON(fiber.Map{"categories": names})
}
