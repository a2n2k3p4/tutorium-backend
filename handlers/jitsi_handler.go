package handlers

import (
	"errors"
	"fmt"
	"time"

	"github.com/a2n2k3p4/tutorium-backend/middlewares"
	"github.com/a2n2k3p4/tutorium-backend/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func JitsiRoutes(app *fiber.App) {
	jitsi_meeting := NewJitsiHandler("https://meet.jit.si")
	jitsi_teacher := app.Group("/jitsi", middlewares.ProtectedMiddleware(), middlewares.TeacherRequired())
	jitsi_teacher.Get("/start/:id", jitsi_meeting.StartMeeting)
	jitsi_learner := app.Group("/jitsi", middlewares.ProtectedMiddleware(), middlewares.LearnerRequired())
	jitsi_learner.Get("/link/:id", jitsi_meeting.GetMeetingLink)
}

// Create baseURL variable type
type JitsiHandler struct {
	BaseURL string `json:"base_url"`
}

func NewJitsiHandler(baseURL string) *JitsiHandler {
	return &JitsiHandler{BaseURL: baseURL}
}

//! For both StartMeeting and GetMeetingLink, we need to enter the class session ID

func (h *JitsiHandler) StartMeeting(c *fiber.Ctx) error {
	user, ok := c.Locals("currentUser").(*models.User)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized teacher not found"})
	}

	// Generate a unique room name using teacher ID and current timestamp
	room := fmt.Sprintf("KUtutorium_%d_%d", user.Teacher.UserID, time.Now().Unix())
	link := fmt.Sprintf("%s/%s", h.BaseURL, room)

	// Save the meeting link to the database
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

	classSession.ClassURL = link
	err = db.Save(&classSession).Error
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.JSON(fiber.Map{
		"meeting_link": link,
	})
}

func (h *JitsiHandler) GetMeetingLink(c *fiber.Ctx) error {
	_, ok := c.Locals("currentUser").(*models.User)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized learner not found"})
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
