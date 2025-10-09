// payment_handler_db.go contains GET and POST handlers for /payments and /transactions
package handlers

import (
	"errors"
	"log"

	"github.com/a2n2k3p4/tutorium-backend/config"
	"github.com/a2n2k3p4/tutorium-backend/models"
	"github.com/a2n2k3p4/tutorium-backend/services"
	"github.com/gofiber/fiber/v2"
	omise "github.com/omise/omise-go"
	"github.com/omise/omise-go/operations"
	"gorm.io/gorm"
)

// CreateCharge godoc
//
//	@Summary		Create a payment charge
//	@Description	Create an Omise charge. For credit cards, prefer using a token. For testing, server-side tokenization via card fields is supported.
//	@Tags			Payments
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		models.PaymentRequest	true	"Payment payload"
//	@Success		200		{object} string	"Omise charge response"
//	@Failure		400		{string} string	"Invalid request"
//	@Failure		500		{string} string	"Server error"
//	@Router			/payments/charge [post]
func (h *PaymentHandler) CreateCharge(c *fiber.Ctx) error {
	var req models.PaymentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(err.Error())
	}
	if req.Amount <= 0 {
		return c.Status(400).JSON("error amount is required")
	}
	// Fill defaults / enforce currency
	if req.Currency == "" {
		req.Currency = config.PAYMENTDefaultCurrency()
	} else if cfg := config.PAYMENTDefaultCurrency(); cfg != "" && req.Currency != cfg {
		return c.Status(400).JSON("error currency must be " + cfg)
	}
	if req.PaymentType == "internet_banking" && req.ReturnURI == "" {
		req.ReturnURI = config.PAYMENTReturnURI()
	}

	// Try to resolve user id from body/header/query
	userID := getUserIDFromRequest(c, &req)

	// Consider per-request idempotency key if provided: use a fresh client to avoid cross-request header bleeding.
	var client *omise.Client = h.Client
	if idk := c.Get("Idempotency-Key"); idk != "" {
		if pk, sk := config.OMISEPublicKey(), config.OMISESecretKey(); pk != "" && sk != "" {
			if cli, err := omise.NewClient(pk, sk); err == nil {
				cli.WithCustomHeaders(map[string]string{"Idempotency-Key": idk})
				client = cli
			}
		}
	}
	svc := services.NewPaymentService(h.DB, client)
	var (
		charge *omise.Charge
		err    error
	)
	charge, err = svc.CreateCharge(req)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	// Persist/Upsert a local transaction row (idempotent on charge_id)
	if err := svc.UpsertTransactionFromCharge(charge, userID); err != nil {
		log.Printf("Failed to save transaction: %v", err) // do not fail outward
	}

	return c.JSON(charge)
}

func (h *PaymentHandler) createCharge(op *operations.CreateCharge) (*omise.Charge, error) { // kept for backward-compat; unused by handlers after refactor
	ch := &omise.Charge{}
	if err := h.Client.Do(ch, op); err != nil {
		return nil, err
	}
	return ch, nil
}

// ListTransactions godoc
//
//	@Summary		List transactions
//	@Description	List transactions with optional filters and pagination.
//	@Tags			Payments
//	@Produce		json
//	@Param			user_id	query		string	false	"Filter by user ID"
//	@Param			status	query		string	false	"Filter by status (e.g. successful, failed)"
//	@Param			channel	query		string	false	"Filter by channel (e.g. card, promptpay)"
//	@Param			limit	query		int		false	"Limit (default 50)"
//	@Param			offset	query		int		false	"Offset (default 0)"
//	@Success		200		{object}	models.TransactionListResponse
//	@Failure		500		{string}	string	"Server error"
//	@Router			/payments/transactions [get]
func (h *PaymentHandler) ListTransactions(c *fiber.Ctx) error {
	f := services.TxFilters{
		UserID:  c.Query("user_id"),
		Status:  c.Query("status"),
		Channel: c.Query("channel"),
	}
	limit, offset := services.HelpersParseLimitOffset(c.Query("limit"), c.Query("offset"))

	svc := services.NewPaymentService(h.DB, h.Client)
	transactions, totalCount, err := svc.ListTransactions(f, limit, offset)
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.JSON(fiber.Map{
		"transactions": transactions,
		"pagination": fiber.Map{
			"total":  totalCount,
			"limit":  limit,
			"offset": offset,
		},
	})
}

// GetTransaction godoc
//
//	@Summary		Get a transaction
//	@Description	Get a transaction by internal ID or Omise charge_id.
//	@Tags			Payments
//	@Produce		json
//	@Param			id	path		string	true	"Transaction ID or charge_id"
//	@Success		200	{object}	models.Transaction
//	@Failure		400	{string}	string	"Invalid transaction ID"
//	@Failure		404	{string}	string	"Not found"
//	@Failure		500	{string}	string	"Server error"
//	@Router			/payments/transactions/{id} [get]
func (h *PaymentHandler) GetTransaction(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON("id is required")
	}

	svc := services.NewPaymentService(h.DB, h.Client)
	tx, err := svc.GetTransaction(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(404).JSON("Transaction not found")
		}
		return c.Status(500).JSON(err.Error())
	}
	return c.JSON(tx)
}

// RefundTransaction godoc
//
//	@Summary		Refund a transaction
//	@Description	Refund an Omise charge by transaction ID or charge_id. Partial refund if amount provided.
//	@Tags			Payments
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string				true	"Transaction ID or charge_id"
//	@Param			payload	body		map[string]int	false	"Refund payload, e.g. {\"amount\": 1000} satang"
//	@Success		200		{object}	map[string]interface{}
//	@Failure		400		{string}	string	"Invalid request"
//	@Failure		404		{string} string	"Transaction not found"
//	@Failure		500		{string} string	"Server error"
//	@Router			/payments/transactions/{id}/refund [post]
func (h *PaymentHandler) RefundTransaction(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON("id is required")
	}

	var body struct {
		Amount *int64 `json:"amount"`
	}
	_ = c.BodyParser(&body)

	svc := services.NewPaymentService(h.DB, h.Client)
	refund, updatedCharge, err := svc.RefundByIDOrCharge(id, body.Amount)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(404).JSON("Transaction not found")
		}
		return c.Status(500).JSON(err.Error())
	}
	return c.JSON(fiber.Map{"refund": refund, "charge": updatedCharge})
}
