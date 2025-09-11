package services

import (
	"errors"
	"strconv"

	omise "github.com/omise/omise-go"
	"github.com/omise/omise-go/operations"
	"gorm.io/gorm"

	"github.com/a2n2k3p4/tutorium-backend/models"
)

// PaymentService encapsulates payment logic and DB access.
type PaymentService struct {
	DB     *gorm.DB
	Client *omise.Client
}

func NewPaymentService(db *gorm.DB, client *omise.Client) *PaymentService {
	return &PaymentService{DB: db, Client: client}
}

// CreateCharge creates a charge depending on PaymentRequest.PaymentType.
func (s *PaymentService) CreateCharge(req models.PaymentRequest) (*omise.Charge, error) {
	switch req.PaymentType {
	case "credit_card":
		return s.processCreditCard(req)
	case "promptpay":
		return s.processPromptPay(req)
	case "internet_banking":
		return s.processInternetBanking(req)
	default:
		return nil, errors.New("unsupported paymentType: " + req.PaymentType)
	}
}

// createCharge executes the Omise CreateCharge operation.
func (s *PaymentService) createCharge(op *operations.CreateCharge) (*omise.Charge, error) {
	ch := &omise.Charge{}
	if err := s.Client.Do(ch, op); err != nil {
		return nil, err
	}
	return ch, nil
}

// ListTransactions returns transactions with total count using optional filters.
func (s *PaymentService) ListTransactions(f TxFilters, limit, offset int) ([]models.Transaction, int64, error) {
	var totalCount int64
	if err := s.DB.Model(&models.Transaction{}).Scopes(helpersApplyTxFilters(f)).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}
	var transactions []models.Transaction
	if err := s.DB.Model(&models.Transaction{}).
		Scopes(helpersApplyTxFilters(f)).
		Preload("User").
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&transactions).Error; err != nil {
		return nil, 0, err
	}
	return transactions, totalCount, nil
}

// GetTransaction retrieves a transaction by numeric PK or charge_id.
func (s *PaymentService) GetTransaction(id string) (models.Transaction, error) {
	var tx models.Transaction
	// If numeric, treat as internal PK; else treat as ChargeID
	if n, err := strconv.ParseUint(id, 10, 64); err == nil {
		if e := s.DB.Preload("User").First(&tx, uint(n)).Error; e == nil {
			return tx, nil
		} else if !errors.Is(e, gorm.ErrRecordNotFound) {
			return models.Transaction{}, e
		}
	}
	if err := s.DB.Preload("User").Where("charge_id = ?", id).First(&tx).Error; err != nil {
		return models.Transaction{}, err
	}
	return tx, nil
}

// RefundByIDOrCharge refunds a charge by internal tx id or charge_id. If amount is nil, refunds full amount.
func (s *PaymentService) RefundByIDOrCharge(id string, amount *int64) (*omise.Refund, *omise.Charge, error) {
	// Resolve charge ID
	chargeID := id
	if !(len(id) > 5 && id[:5] == "chrg_") {
		// Not a charge id, look up transaction
		tx, err := s.GetTransaction(id)
		if err != nil {
			return nil, nil, err
		}
		chargeID = tx.ChargeID
	}

	// Create refund
	create := &operations.CreateRefund{ChargeID: chargeID}
	if amount != nil && *amount > 0 {
		create.Amount = *amount
	}
	rf := &omise.Refund{}
	if err := s.Client.Do(rf, create); err != nil {
		return nil, nil, err
	}

	// Retrieve updated charge and upsert locally
	ch := &omise.Charge{}
	if err := s.Client.Do(ch, &operations.RetrieveCharge{ChargeID: chargeID}); err != nil {
		return rf, nil, err
	}
	if err := s.UpsertTransactionFromCharge(ch, nil); err != nil {
		// Log-only in service layer: let caller still get refund info
		// but bubble up error for visibility.
		return rf, ch, err
	}
	return rf, ch, nil
}
