package middlewares

import (
	"fmt"
	"strings"

	"github.com/a2n2k3p4/tutorium-backend/config"
	"github.com/a2n2k3p4/tutorium-backend/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

// if Status = development you can bypass all routes
var Status = config.STATUS
var Secret = func() []byte {
	return []byte(config.JWTSecret())
}

// SetSecret lets tests (or a custom main) override how the secret is provided.
// Call this once, *before* the middleware is used. Passing nil restores default.
func SetSecret(f func() []byte) {
	if f == nil {
		Secret = func() []byte {
			return []byte(config.JWTSecret())
		}
		return
	}
	Secret = f
}

/* ------------------------------- middleware ------------------------------ */

func ProtectedMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if Status() == "development" {
			return c.Next()
		}
		authHeader := c.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(401).JSON(fiber.Map{"error": "missing or invalid token"})
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return Secret(), nil
		})

		if err != nil || !token.Valid {
			return c.Status(401).JSON(fiber.Map{"error": "invalid token", "details": err.Error()})
		}

		db, err := GetDB(c)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "db not available"})
		}

		var user models.User
		if err := db.Preload("Learner").Preload("Teacher").Preload("Admin").
			First(&user, claims.UserID).Error; err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "invalid credentials"})
		}

		c.Locals("currentUser", &user)
		return c.Next()
	}
}

func AdminRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if Status() == "development" {
			return c.Next()
		}
		user, ok := c.Locals("currentUser").(*models.User)
		if !ok {
			return c.Status(401).JSON(fiber.Map{"error": "authentication required"})
		}
		if user.Admin == nil {
			return c.Status(403).JSON(fiber.Map{"error": "admin access required"})
		}
		return c.Next()
	}
}

func TeacherRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if Status() == "development" {
			return c.Next()
		}
		user, ok := c.Locals("currentUser").(*models.User)
		if !ok {
			return c.Status(401).JSON(fiber.Map{"error": "authentication required"})
		}
		if user.Teacher == nil {
			return c.Status(403).JSON(fiber.Map{"error": "teacher access required"})
		}
		return c.Next()
	}
}

func LearnerRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if Status() == "development" {
			return c.Next()
		}
		user, ok := c.Locals("currentUser").(*models.User)
		if !ok {
			return c.Status(401).JSON(fiber.Map{"error": "authentication required"})
		}
		if user.Learner == nil {
			return c.Status(403).JSON(fiber.Map{"error": "learner access required"})
		}
		return c.Next()
	}
}
