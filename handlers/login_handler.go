package handlers

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/a2n2k3p4/tutorium-backend/middlewares"
	"github.com/a2n2k3p4/tutorium-backend/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func LoginRoutes(app *fiber.App) {
	app.Post("/login", LoginHandler)
}

// LoginHandler godoc
//
//	@Summary		Login with KU/Nisit credentials
//	@Description	Authenticate a nisit user via KU API, create the user if not exists, and return a JWT token along with user info
//	@Tags			Login
//	@Accept			json
//	@Produce		json
//	@Param			login	body		models.LoginRequestDoc	true	"Login payload"
//	@Success		200		{object}	models.LoginResponseDoc
//	@Failure		400		{object}	map[string]string	"Invalid input"
//	@Failure		401		{object}	map[string]string	"Unauthorized"
//	@Failure		500		{object}	map[string]string	"Server error"
//	@Router			/login [post]
func LoginHandler(c *fiber.Ctx) error {
	type LoginRequest struct {
		Username       string `json:"username"`
		Password       string `json:"password"`
		ProfilePicture string `json:"profile_picture,omitempty"`
		FirstName      string `json:"first_name"`
		LastName       string `json:"last_name"`
		Gender         string `json:"gender"`
		PhoneNumber    string `json:"phone_number"`
	}

	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if err := godotenv.Load("../.env"); err != nil {
		log.Println(".env file not found, using system environment variables")
	}
	NisitKUBaseURL := os.Getenv("KU_API")
	nisitClient := NewNisitKUClient(NisitKUBaseURL)

	loginResp, err := nisitClient.Login(req.Username, req.Password)
	if err != nil {
		return c.Status(401).JSON(err.Error())
	}

	if loginResp == nil || loginResp.Status != "true" {
		return c.Status(401).JSON(fiber.Map{"error": "invalid credentials"})
	}

	// Decode ProfilePicture from string to []byte and validate it
	var profileBytes []byte
	if req.ProfilePicture != "" {
		profileBytes, err = decodeBase64Image(req.ProfilePicture)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid profile_picture", "detail": err.Error()})
		}
		// validate size and mime type
		if err := validateImageBytes(profileBytes); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid profile_picture", "detail": err.Error()})
		}
	}
	db, err := middlewares.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	var user models.User
	err = db.Where("student_id = ?", loginResp.ID).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		user = models.User{
			StudentID:      loginResp.ID,
			FirstName:      req.FirstName,
			LastName:       req.LastName,
			Gender:         req.Gender,
			PhoneNumber:    req.PhoneNumber,
			ProfilePicture: profileBytes,
			Balance:        0,
		}
		if err := db.Create(&user).Error; err != nil {
			return c.Status(500).JSON(err.Error())
		}
	} else if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	token, err := generateJWT(user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "cannot generate token"})
	}

	return c.JSON(fiber.Map{
		"user":  user,
		"token": token,
	})
}

func generateJWT(user models.User) (string, error) {
	if err := godotenv.Load("../.env"); err != nil {
		log.Println(".env file not found, using system environment variables")
	}
	secretStr := os.Getenv("JWT_SECRET")
	if secretStr == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}
	secret := []byte(secretStr)
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // token expires in 24h
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

// decodeBase64Image decodes data: URI or plain base64 payload.
func decodeBase64Image(s string) ([]byte, error) {
	if s == "" {
		return nil, nil
	}
	if strings.HasPrefix(s, "data:") {
		comma := strings.IndexByte(s, ',')
		if comma < 0 {
			return nil, errors.New("invalid data URI")
		}
		s = s[comma+1:]
	}
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, fmt.Errorf("base64 decode: %w", err)
	}
	return b, nil
}

// validateImageBytes checks size and MIME type of the image bytes.
func validateImageBytes(b []byte) error {
	const MaxProfileImageBytes = 2 * 1024 * 1024
	if len(b) == 0 {
		return nil // nothing to validate
	}
	if len(b) > MaxProfileImageBytes {
		return fmt.Errorf("image too large (max %d bytes)", MaxProfileImageBytes)
	}
	// Detect content type from first 512 bytes (http.DetectContentType)
	sz := 512
	if len(b) < sz {
		sz = len(b)
	}
	mtype := http.DetectContentType(b[:sz])
	switch mtype {
	case "image/jpeg", "image/png", "image/gif", "image/webp":
		return nil
	default:
		return fmt.Errorf("unsupported image type: %s", mtype)
	}
}
