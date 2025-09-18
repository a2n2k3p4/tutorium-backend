package middlewares

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	omise "github.com/omise/omise-go"
)

const omiseCtxKey = "tutorium_omise_client"

// OmiseMiddleware injects a shared Omise client into the request context.
func OmiseMiddleware(client *omise.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Locals(omiseCtxKey, client)
		return c.Next()
	}
}

// GetOmise extracts the Omise client from the request context.
func GetOmise(c *fiber.Ctx) (*omise.Client, error) {
	v := c.Locals(omiseCtxKey)
	if v == nil {
		return nil, errors.New("omise client not found in context")
	}
	cli, ok := v.(*omise.Client)
	if !ok {
		return nil, errors.New("invalid omise client in context")
	}
	return cli, nil
}
