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
	jitsi_teacher := app.Group("/jitsi", middlewares.ProtectedMiddleware(), middlewares.TeacherRequired())
	jitsi_teacher.Get("/start", jitsi_meeting.StartMeeting)
	jitsi_learner := app.Group("/jitsi", middlewares.ProtectedMiddleware(), middlewares.LearnerRequired())
	jitsi_learner.Get("/link", jitsi_meeting.GetMeetingLink)
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

func (h *JitsiHandler) GetMeetingLink(c *fiber.Ctx) error {
	user, ok := c.Locals("currentUser").(*models.User)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized learner not found"})
	}

	db, err := middlewares.GetDB(c)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	// Get the latest class session for the learner's enrolled classes
	var classSession models.ClassSession
	err = db.Joins("JOIN Enrollment ON Enrollment.ClassSessionID = ClassSession.ClassID").
		Where("Enrollment.LearnerID = ?", user.Learner.UserID).
		Order("ClassSession.CreatedAT DESC").
		First(&classSession).Error
	if err != nil {
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
