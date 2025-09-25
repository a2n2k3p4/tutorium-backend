package handlers

import (
	"errors"

	"github.com/a2n2k3p4/tutorium-backend/middlewares"
	"github.com/a2n2k3p4/tutorium-backend/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func MeetingRoutes(app *fiber.App) {
	meetingurl := NewMeetingHandler()
	meeting := app.Group("/meeting", middlewares.ProtectedMiddleware(), middlewares.TeacherRequired(), middlewares.LearnerRequired())
	meeting.Get("/:id", meetingurl.GetMeetingLink)
}

// Create baseURL variable type
type MeetingURL struct {
	BaseURL string `json:"base_url"`
}

func NewMeetingHandler() *MeetingURL {
	return &MeetingURL{BaseURL: "https://meet.jit.si"}
}

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

	link := classSession.ClassURL
	if link == "" {
		return c.Status(404).JSON(fiber.Map{"error": "no meeting was created for your class session yet"})
	}
	return c.JSON(fiber.Map{
		"meeting_link": link,
	})
}
