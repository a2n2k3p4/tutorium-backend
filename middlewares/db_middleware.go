package middlewares

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

const dbCtxKey = "tutorium_db"

func DBMiddleware(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Locals(dbCtxKey, db)
		return c.Next()
	}
}

func GetDB(c *fiber.Ctx) (*gorm.DB, error) {
	v := c.Locals(dbCtxKey)
	if v == nil {
		return nil, errors.New("database not found in context")
	}
	db, ok := v.(*gorm.DB)
	if !ok {
		return nil, errors.New("invalid database object in context")
	}
	return db, nil
}
