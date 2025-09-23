package handlers

import (
	"fmt"
	"time"

	"github.com/a2n2k3p4/tutorium-backend/middlewares"
	"github.com/a2n2k3p4/tutorium-backend/models"
	"github.com/gofiber/fiber/v2"
)

func JitsiRoutes(app *fiber.App) {
	jitsi_meeting := NewJitsiHandler("https://meet.jit.si")
	jitsi := app.Group("/jitsi", middlewares.ProtectedMiddleware(), middlewares.TeacherRequired())
	jitsi.Get("/start", jitsi_meeting.StartMeeting)
}

// Create baseURL variable type
type JitsiHandler struct {
	BaseURL string `json:"base_url"`
}

func NewJitsiHandler(baseURL string) *JitsiHandler {
	return &JitsiHandler{BaseURL: baseURL}
}

func (h *JitsiHandler) StartMeeting(c *fiber.Ctx) error {
	user, ok := c.Locals("currentUser").(*models.User)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized teacher not found"})
	}

	// Generate a unique room name using teacher ID and current timestamp
	room := fmt.Sprintf("KUtutorium_%d_%d", user.Teacher.UserID, time.Now().Unix())
	link := fmt.Sprintf("%s/%s", h.BaseURL, room)

	return c.JSON(fiber.Map{
		"meeting_link": link,
	})
}
