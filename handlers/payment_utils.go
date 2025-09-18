package handlers

import (
	"strconv"

	"github.com/a2n2k3p4/tutorium-backend/models"
	"github.com/gofiber/fiber/v2"
)

// getUserIDFromRequest inspects body, header, and query to resolve a user ID.
func getUserIDFromRequest(c *fiber.Ctx, req *models.PaymentRequest) *uint {
	if cu, ok := c.Locals("currentUser").(*models.User); ok && cu != nil && cu.ID != 0 {
		u := cu.ID
		return &u
	}
	if req != nil && req.UserID != nil {
		return req.UserID
	}
	if userIDHeader := c.Get("X-User-ID"); userIDHeader != "" {
		if userID, err := strconv.ParseUint(userIDHeader, 10, 32); err == nil {
			u := uint(userID)
			return &u
		}
	}
	if userIDQuery := c.Query("user_id"); userIDQuery != "" {
		if userID, err := strconv.ParseUint(userIDQuery, 10, 32); err == nil {
			u := uint(userID)
			return &u
		}
	}
	return nil
}
