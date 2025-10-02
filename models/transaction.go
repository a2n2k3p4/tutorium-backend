package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Transaction represents an Omise charge persisted locally for auditing and balance updates.
type Transaction struct {
	ID             uint              `gorm:"primaryKey" json:"id"`
	CreatedAt      time.Time         `gorm:"index" json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
	DeletedAt      gorm.DeletedAt    `gorm:"index" json:"-"`
	UserID         *uint             `gorm:"index" json:"user_id,omitempty"`
	ChargeID       string            `gorm:"uniqueIndex;size:64" json:"charge_id"`
	AmountSatang   int64             `json:"amount_satang"`
	Currency       string            `gorm:"size:8;index" json:"currency"`
	Channel        string            `gorm:"size:64;index" json:"channel"`
	Status         string            `gorm:"size:32;index" json:"status"`
	FailureCode    *string           `json:"failure_code,omitempty"`
	FailureMessage *string           `json:"failure_message,omitempty"`
	RawPayload     []byte            `json:"-"`
	Meta           datatypes.JSONMap `gorm:"type:jsonb" json:"meta,omitempty" swaggertype:"object"`

	User *User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"-"`
}

// PaymentRequest is the payload from your frontend to initiate a charge.
type PaymentRequest struct {
	Amount      int64                  `json:"amount"`               // (satang unit : 100 satang = 1 THB)
	Currency    string                 `json:"currency"`             // "THB"
	PaymentType string                 `json:"paymentType"`          // "credit_card" | "promptpay" | "internet_banking"
	Token       string                 `json:"token,omitempty"`      // for card charges (preferred)
	ReturnURI   string                 `json:"return_uri,omitempty"` // required for some redirects (3DS/internet banking)
	Description string                 `json:"description,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"` // free-form, attached to the Omise charge
	Card        map[string]interface{} `json:"card,omitempty"`     // server-side tokenization (TESTING ONLY)
	Bank        string                 `json:"bank,omitempty"`     // e.g. "bbl", "bay", "scb"
	UserID      *uint                  `json:"user_id,omitempty"`  // FK to users.id
}

// TransactionListResponse is a doc helper for Swagger representing the list response.
type TransactionListResponse struct {
	Transactions []Transaction `json:"transactions"`
	Pagination   struct {
		Total  int64 `json:"total"`
		Limit  int   `json:"limit"`
		Offset int   `json:"offset"`
	} `json:"pagination"`
}

// OmiseWebhookPayload represents the webhook payload from Omise.
// It can be either an Event object or a Charge object.
type OmiseWebhookPayload struct {
	Object string `json:"object" example:"charge"` // "event" or "charge"
	ID     string `json:"id" example:"chrg_test_658q8luocil7hlhd07n"`
}
