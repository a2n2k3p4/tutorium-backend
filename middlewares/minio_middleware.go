package middlewares

import (
	"github.com/a2n2k3p4/tutorium-backend/storage"
	"github.com/gofiber/fiber/v2"
)

func MinioMiddleware(client *storage.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Locals("minio", client)
		return c.Next()
	}
}
