package handlers

import (
	"errors"
	"log"
	"os"

	"github.com/a2n2k3p4/tutorium-backend/middlewares"
	"github.com/a2n2k3p4/tutorium-backend/models"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func MeetingRoutes(app *fiber.App) {
	meetingurl := NewMeetingHandler()
	meeting := app.Group("/meetings", middlewares.ProtectedMiddleware(), middlewares.BanMiddleware(), middlewares.TeacherRequired(), middlewares.LearnerRequired())
	meeting.Get("/:id", meetingurl.GetMeetingLink)
}

var BASE_JITSI_URL string

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Warning: could not load .env file")
	}
	BASE_JITSI_URL = os.Getenv("JITSI_URL")
}

// Create baseURL variable type
type MeetingURL struct {
	BaseURL string `json:"base_url"`
}

func NewMeetingHandler() *MeetingURL {
	return &MeetingURL{BaseURL: BASE_JITSI_URL}
}

// GetMeetingLink godoc
//
//	@Summary		Get meeting link by ClassSession ID
//	@Description	Retrieves the meeting link associated with a given ClassSession ID.
//	@Tags			Meetings
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id	path		int					true	"ClassSession ID"
//	@Success		200	{object}	map[string]string	"meeting_link"
//	@Failure		400	{object}	map[string]string	"Invalid class session ID"
//	@Failure		401	{object}	map[string]string	"Unauthorized"
//	@Failure		404	{object}	map[string]string	"Class session not found or meeting not created"
//	@Failure		500	{string}	string				"Server error"
//	@Router			/meetings/{id} [get]
func (h *MeetingURL) GetMeetingLink(c *fiber.Ctx) error {
	_, ok := c.Locals("currentUser").(*models.User)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	db, err := middlewares.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	// ------- Get class session ID from params -------
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid class session ID"})
	}

	var classSession models.ClassSession

	err = db.First(&classSession, "id = ?", id).Error
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return c.Status(404).JSON(fiber.Map{"error": "class session not found"})
	case err != nil:
		return c.Status(500).JSON(err.Error())
	}

	link := classSession.MeetingUrl
	if link == "" {
		return c.Status(404).JSON(fiber.Map{"error": "no meeting was created for your class session yet"})
	}
	return c.JSON(fiber.Map{
		"meeting_link": link,
	})
}
