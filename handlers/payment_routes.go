package handlers

import (
	"log"

	"github.com/a2n2k3p4/tutorium-backend/middlewares"
	"github.com/gofiber/fiber/v2"
)

// PaymentRoutes registers Omise payment-related endpoints.
func PaymentRoutes(app *fiber.App) {
    // Simple health check
    app.Get("/health", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{"status": "ok"})
    })

    // Optional return URI endpoint for redirect-based flows (3DS/IBANKING)
    app.Get("/payments/return", func(c *fiber.Ctx) error {
        // Omise will redirect users back here; you can enhance this to show a page.
        return c.JSON(fiber.Map{"status": "returned"})
    })

	// Charges
	app.Post("/payments/charge", func(c *fiber.Ctx) error {
		db, err := middlewares.GetDB(c)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "db not available"})
		}
		client, err := middlewares.GetOmise(c)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "payment client not available"})
		}
		h := NewPaymentHandler(db, client)
		return h.CreateCharge(c)
	})

	// Transactions
	app.Get("/payments/transactions", func(c *fiber.Ctx) error {
		db, err := middlewares.GetDB(c)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "db not available"})
		}
		client, err := middlewares.GetOmise(c)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "payment client not available"})
		}
		h := NewPaymentHandler(db, client)
		return h.ListTransactions(c)
	})

	app.Get("/payments/transactions/:id", func(c *fiber.Ctx) error {
		db, err := middlewares.GetDB(c)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "db not available"})
		}
		client, err := middlewares.GetOmise(c)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "payment client not available"})
		}
		h := NewPaymentHandler(db, client)
		return h.GetTransaction(c)
	})

	// Refund
	app.Post("/payments/transactions/:id/refund", func(c *fiber.Ctx) error {
		db, err := middlewares.GetDB(c)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "db not available"})
		}
		client, err := middlewares.GetOmise(c)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "payment client not available"})
		}
		h := NewPaymentHandler(db, client)
		return h.RefundTransaction(c)
	})

	// Webhook
	app.Post("/webhooks/omise", func(c *fiber.Ctx) error {
		db, err := middlewares.GetDB(c)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "db not available"})
		}
		client, err := middlewares.GetOmise(c)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "payment client not available"})
		}
		h := NewPaymentHandler(db, client)
		if err := h.HandleWebhook(c); err != nil {
			log.Printf("webhook handler error: %v", err)
			return err
		}
		return nil
	})
}
