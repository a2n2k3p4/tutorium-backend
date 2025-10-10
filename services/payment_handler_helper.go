// payment_handler_helper.go contains helper functions for Omise payment handling
package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/a2n2k3p4/tutorium-backend/models"
	omise "github.com/omise/omise-go"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TxFilters struct {
	UserID  string
	Status  string
	Channel string
}

// ---------------------- payment helpers ----------------------
// (helper for ListTransactions) GORM scope for queries with optional filters: user, status, and channel.
func helpersApplyTxFilters(f TxFilters) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if f.UserID != "" {
			db = db.Where("user_id = ?", f.UserID)
		}
		if f.Status != "" {
			db = db.Where("status = ?", f.Status)
		}
		if f.Channel != "" {
			db = db.Where("channel = ?", f.Channel)
		}
		return db
	}
}

// (helper for ListTransactions) safe pagination defaults.
func HelpersParseLimitOffset(limitStr, offsetStr string) (int, int) {
	limit, offset := 50, 0
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		limit = l
	}
	if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
		offset = o
	}
	return limit, offset
}

// ---------------------- webhook helpers ----------------------
// (HandleWebhook helper) update-insert a local transaction row from Omise Charge
// upsertTransactionFromCharge updates/creates the local transaction and adjusts user balance
// only on status transitions across the "successful" boundary.
func (s *PaymentService) UpsertTransactionFromCharge(charge *omise.Charge, userID *uint) error {
	if charge == nil {
		return fmt.Errorf("nil charge")
	}
	userID = extractUserIDFromCharge(charge, userID)
	channel := determineChannel(charge)
	rawPayload, _ := json.Marshal(charge)

	var meta datatypes.JSONMap
	if charge.Metadata != nil {
		meta = datatypes.JSONMap(charge.Metadata)
	}

	tx := s.DB.Begin()
	if err := tx.Error; err != nil {
		return err
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	var prev models.Transaction
	if err := tx.
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("charge_id = ?", charge.ID).
		Take(&prev).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()
		return err
	}
	prevWasSuccessful := prev.Status == "successful"

	if userID == nil && prev.UserID != nil {
		userID = prev.UserID
	}

	newTx := models.Transaction{
		UserID:         userID,
		ChargeID:       charge.ID,
		AmountSatang:   charge.Amount,
		Currency:       charge.Currency,
		Channel:        channel,
		Status:         string(charge.Status),
		FailureCode:    charge.FailureCode,
		FailureMessage: charge.FailureMessage,
		RawPayload:     rawPayload,
		Meta:           meta,
	}
	if err := tx.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "charge_id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"status", "failure_code", "failure_message",
			"amount_satang", "currency", "channel",
			"raw_payload", "meta", "updated_at", "user_id",
		}),
	}).Create(&newTx).Error; err != nil {
		tx.Rollback()
		return err
	}

	if userID != nil {
		if err := s.adjustUserBalanceOnStatusTransition(tx, charge, userID, prevWasSuccessful); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// adjustUserBalanceOnStatusTransition handles user balance adjustment logic for status transitions.
func (s *PaymentService) adjustUserBalanceOnStatusTransition(tx *gorm.DB, charge *omise.Charge, userID *uint, prevWasSuccessful bool) error {
	nowSuccessful := string(charge.Status) == "successful"
	switch {
	case !prevWasSuccessful && nowSuccessful:
		amountTHB := float64(charge.Amount) / 100.0 // convert satang to THB
		if err := tx.Model(&models.User{}).
			Where("id = ?", *userID).
			Update("balance", gorm.Expr("balance + ?", amountTHB)).Error; err != nil {
			log.Printf("Failed to credit user balance: %v", err)
			return err
		}
	case prevWasSuccessful && !nowSuccessful:
		// Debit balance if a previously successful charge became non-successful (e.g., refund/reversal)
		amountTHB := float64(charge.Amount) / 100.0
		if err := tx.Model(&models.User{}).
			Where("id = ?", *userID).
			Update("balance", gorm.Expr("balance - ?", amountTHB)).Error; err != nil {
			log.Printf("Failed to debit user balance: %v", err)
			return err
		}
	}
	return nil
}

func determineChannel(charge *omise.Charge) string {
	if charge == nil {
		return "card"
	}
	if charge.Source != nil && charge.Source.Type != "" {
		return charge.Source.Type
	}
	return "card"
}

// (helper for helpersParseLimitOffset)
func extractUserIDFromCharge(charge *omise.Charge, userID *uint) *uint {
	if userID != nil {
		return userID
	}
	if charge == nil || charge.Metadata == nil {
		return userID
	}
	if v, ok := charge.Metadata["user_id"]; ok {
		switch vv := v.(type) {
		case string:
			if n, err := strconv.ParseUint(vv, 10, 32); err == nil {
				u := uint(n)
				return &u
			}
		case float64:
			u := uint(vv)
			return &u
		case json.Number:
			if n, err := vv.Int64(); err == nil {
				u := uint(n)
				return &u
			}
		}
	}
	return userID
}
