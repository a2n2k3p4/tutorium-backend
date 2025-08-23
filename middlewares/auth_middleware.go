package middlewares

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/a2n2k3p4/tutorium-backend/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

/* ---------- secret provider (prod uses env, tests can override) ---------- */

var (
	secretOnce     sync.Once
	cachedSecret   []byte
	secretProvider = defaultSecretProvider
)

// default: load from env (.env for local dev)
func defaultSecretProvider() []byte {
	secretOnce.Do(func() {
		// Try to load ../.env once (no-op if file missing)
		_ = godotenv.Load("../.env")
		s := os.Getenv("JWT_SECRET")
		if s == "" {
			log.Fatal("JWT_SECRET environment variable is not set")
		}
		cachedSecret = []byte(s)
	})
	return cachedSecret
}

// SetSecretProvider allows tests (or main) to inject a secret getter.
// Pass nil to restore default env-based provider.
func SetSecretProvider(p func() []byte) {
	if p == nil {
		secretProvider = defaultSecretProvider
		return
	}
	// reset any previously cached secret
	secretOnce = sync.Once{}
	cachedSecret = nil
	secretProvider = p
}

/* ------------------------------- middleware ------------------------------ */

func ProtectedMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(401).JSON(fiber.Map{"error": "missing or invalid token"})
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &Claims{}

		secret := secretProvider()

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return secret, nil
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
