package middlewares

import (
	"github.com/a2n2k3p4/tutorium-backend/storage"
	"github.com/gofiber/fiber/v2"
)

func MinioMiddleware(u storage.Uploader) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Locals("minio", u)
		return c.Next()
	}
}
