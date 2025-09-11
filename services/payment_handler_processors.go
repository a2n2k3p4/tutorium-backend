package services

import (
	"fmt"
	"strconv"
	"time"

	"github.com/a2n2k3p4/tutorium-backend/models"
	omise "github.com/omise/omise-go"
	"github.com/omise/omise-go/operations"
)

// ---------------------- processors ----------------------
func (s *PaymentService) processCreditCard(req models.PaymentRequest) (*omise.Charge, error) {
	// Attach user_id to metadata if present (Omise supports custom metadata).
	metadata := req.Metadata
	if req.UserID != nil {
		if metadata == nil {
			metadata = make(map[string]interface{})
		}
		metadata["user_id"] = fmt.Sprintf("%d", *req.UserID)
	}

	// Preferred flow: card token already created by frontend (Omise.js / mobile SDK).
	if req.Token != "" {
		return s.createCharge(&operations.CreateCharge{
			Amount:      req.Amount,
			Currency:    req.Currency,
			Card:        req.Token,
			ReturnURI:   req.ReturnURI,
			Description: req.Description,
			Metadata:    metadata,
		})
	}

	// Server-side tokenization (testing only)
	if req.Card == nil {
		return nil, fmt.Errorf("missing token; either provide token or card for tokenization")
	}
	name, _ := req.Card["name"].(string)
	number, _ := req.Card["number"].(string)

	var expMonth, expYear int
	var securityCode string

	switch v := req.Card["expiration_month"].(type) {
	case float64:
		expMonth = int(v)
	case string:
		n, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("invalid expiration_month: %v", v)
		}
		expMonth = n
	default:
		return nil, fmt.Errorf("unexpected type for expiration_month: %T", v)
	}
	switch v := req.Card["expiration_year"].(type) {
	case float64:
		expYear = int(v)
	case string:
		n, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("invalid expiration_year: %v", v)
		}
		expYear = n
	default:
		return nil, fmt.Errorf("unexpected type for expiration_year: %T", v)
	}
	switch v := req.Card["security_code"].(type) {
	case string:
		securityCode = v
	case float64:
		securityCode = strconv.Itoa(int(v))
	default:
		return nil, fmt.Errorf("unexpected type for security_code: %T", v)
	}

	token := &omise.Token{}
	if err := s.Client.Do(token, &operations.CreateToken{
		Name:            name,
		Number:          number,
		ExpirationMonth: time.Month(expMonth),
		ExpirationYear:  expYear,
		SecurityCode:    securityCode,
	}); err != nil {
		return nil, fmt.Errorf("failed to create token: %v", err)
	}

	return s.createCharge(&operations.CreateCharge{
		Amount:      req.Amount,
		Currency:    req.Currency,
		Card:        token.ID,
		ReturnURI:   req.ReturnURI,
		Description: req.Description,
		Metadata:    metadata,
	})
}

func (s *PaymentService) processPromptPay(req models.PaymentRequest) (*omise.Charge, error) {
	// Create a source with type "promptpay", then create a charge from it.
	metadata := req.Metadata
	if req.UserID != nil {
		if metadata == nil {
			metadata = make(map[string]interface{})
		}
		metadata["user_id"] = fmt.Sprintf("%d", *req.UserID)
	}

	src := &omise.Source{}
	if err := s.Client.Do(src, &operations.CreateSource{
		Type:     "promptpay",
		Amount:   req.Amount,
		Currency: req.Currency,
	}); err != nil {
		return nil, fmt.Errorf("failed to create promptpay source: %v", err)
	}

	return s.createCharge(&operations.CreateCharge{
		Amount:      req.Amount,
		Currency:    req.Currency,
		Source:      src.ID,
		Description: req.Description,
		Metadata:    metadata,
	})
}

func (s *PaymentService) processInternetBanking(req models.PaymentRequest) (*omise.Charge, error) {
	// Internet banking requires a source like "internet_banking_bbl", "internet_banking_scb", etc.
	if req.Bank == "" {
		return nil, fmt.Errorf(`bank is required for internet_banking (e.g. "bay", "bbl", "scb")`)
	}
	if req.ReturnURI == "" {
		return nil, fmt.Errorf("return_uri is required for internet_banking")
	}

	metadata := req.Metadata
	if req.UserID != nil {
		if metadata == nil {
			metadata = make(map[string]interface{})
		}
		metadata["user_id"] = fmt.Sprintf("%d", *req.UserID)
	}

	src := &omise.Source{}
	if err := s.Client.Do(src, &operations.CreateSource{
		Type:     "internet_banking_" + req.Bank,
		Amount:   req.Amount,
		Currency: req.Currency,
	}); err != nil {
		return nil, fmt.Errorf("failed to create internet banking source: %v", err)
	}

	return s.createCharge(&operations.CreateCharge{
		Amount:      req.Amount,
		Currency:    req.Currency,
		Source:      src.ID,
		ReturnURI:   req.ReturnURI,
		Description: req.Description,
		Metadata:    metadata,
	})
}
