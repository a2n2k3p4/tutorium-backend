package handlers

import (
	"encoding/json"
	"log"

	"github.com/gofiber/fiber/v2"
	omise "github.com/omise/omise-go"
	"github.com/omise/omise-go/operations"
	"gorm.io/gorm"

	"github.com/a2n2k3p4/tutorium-backend/models"
	"github.com/a2n2k3p4/tutorium-backend/services"
)

type PaymentHandler struct {
	DB     *gorm.DB
	Client *omise.Client
}

func NewPaymentHandler(db *gorm.DB, client *omise.Client) *PaymentHandler {
	return &PaymentHandler{DB: db, Client: client}
}

// Health godoc
//
//	@Summary		Health check
//	@Description	Simple health check for payment service.
//	@Tags			Payments
//	@Produce		json
//	@Success		200	{object}	map[string]string
//	@Router			/health [get]
func (h *PaymentHandler) Health(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "ok"})
}

// HandleWebhook accepts either an Event payload (object:"event") or a Charge payload (object:"charge").
// Flow:
//   - if event: RetrieveEvent -> extract charge.id -> RetrieveCharge -> upsert
//   - if charge: RetrieveCharge -> upsert
//
// Return 5xx on transient failure (so Omise retries); 200 when processed or intentionally ignored.
// HandleWebhook godoc
//
//	@Summary		Omise webhook
//	@Description	Handles Omise events by verifying and upserting transaction status. Accepts either an Event payload (object:"event") or a Charge payload (object:"charge").
//	@Tags			Payments
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		models.OmiseWebhookPayload	true	"Omise webhook payload (event or charge object)"
//	@Success		200		{string}	string						"OK"
//	@Failure		400		{object}	map[string]string			"Bad request"
//	@Failure		500		{string}	string						"Retryable server error"
//	@Router			/webhooks/omise [post]
func (h *PaymentHandler) HandleWebhook(c *fiber.Ctx) error {
	var envelope struct {
		Object string `json:"object"`
		ID     string `json:"id"`
	}
	if err := json.Unmarshal(c.Body(), &envelope); err != nil || envelope.ID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload: missing object or id"})
	}

	var chargeID string

	switch envelope.Object {
	case "event":
		// Verify the event by retrieving it from Omise
		ev := &omise.Event{}
		if err := h.Client.Do(ev, &operations.RetrieveEvent{EventID: envelope.ID}); err != nil {
			log.Printf("webhook: verify event failed id=%s err=%v", envelope.ID, err)
			// Returning 5xx allows the sender to retry (useful for transient network issues).
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		// Extract the embedded object; only handle charge
		var embedded struct {
			ID     string `json:"id"`
			Object string `json:"object"`
		}
		raw, err := json.Marshal(ev.Data)
		if err != nil {
			log.Printf("webhook: marshal ev.Data failed id=%s err=%v", envelope.ID, err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		if err := json.Unmarshal(raw, &embedded); err != nil || embedded.Object != "charge" || embedded.ID == "" {
			// Not a charge-related event → acknowledge and exit.
			return c.SendStatus(fiber.StatusOK)
		}
		chargeID = embedded.ID

	case "charge":
		// Some dashboard/testing tools show the charge payload directly.
		chargeID = envelope.ID

	default:
		// Ignore other payload types.
		return c.SendStatus(fiber.StatusOK)
	}

	// Retrieve the charge to independently verify status, then upsert locally.
	ch := &omise.Charge{}
	if err := h.Client.Do(ch, &operations.RetrieveCharge{ChargeID: chargeID}); err != nil {
		log.Printf("webhook: retrieve charge failed charge=%s err=%v", chargeID, err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	svc := services.NewPaymentService(h.DB, h.Client)
	if err := svc.UpsertTransactionFromCharge(ch, nil); err != nil {
		log.Printf("webhook: upsert failed charge=%s err=%v", ch.ID, err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	log.Printf("webhook: processed charge=%s status=%s amount=%d source=%v", ch.ID, ch.Status, ch.Amount, ch.Source)
	return c.SendStatus(fiber.StatusOK)
}
