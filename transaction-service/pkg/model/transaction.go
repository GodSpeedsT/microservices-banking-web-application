package model

import (
	"github.com/google/uuid"
	"time"
)

type TransactionType string

const (
	TransactionTypeDeposit    TransactionType = "DEPOSIT"
	TransactionTypeWithdrawal TransactionType = "WITHDRAWAL"
	TransactionTypeTransfer   TransactionType = "TRANSFER"
	TransactionTypeInterest   TransactionType = "INTEREST"
)

type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "PENDING"
	TransactionStatusCompleted TransactionStatus = "COMPLETED"
	TransactionStatusFailed    TransactionStatus = "FAILED"
	TransactionStatusCancelled TransactionStatus = "CANCELLED"
)

type Transaction struct {
	ID          string                 `json:"id" db:"id"`
	UserID      string                 `json:"user_id" db:"user_id"`
	AccountID   string                 `json:"account_id" db:"account_id"`
	Amount      float64                `json:"amount" db:"amount"`
	Currency    string                 `json:"currency" db:"currency"`
	Type        TransactionType        `json:"type" db:"type"`
	Status      TransactionStatus      `json:"status" db:"status"`
	Description string                 `json:"description" db:"description"`
	Reference   string                 `json:"reference" db:"reference"`
	Metadata    map[string]interface{} `json:"metadata" db:"metadata"`
	CreatedAt   time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at" db:"updated_at"`
}

type TransactionRequest struct {
	UserID      string          `json:"user_id" binding:"required"`
	AccountID   string          `json:"account_id" binding:"required"`
	Amount      float64         `json:"amount" binding:"required,gt=0"`
	Currency    string          `json:"currency" binding:"required"`
	Type        TransactionType `json:"type" binding:"required"`
	Description string          `json:"description"`
	Reference   string          `json:"reference"`
}

type TransactionResponse struct {
	Transaction *Transaction `json:"transaction"`
	Message     string       `json:"message"`
}

func NewTransaction(req *TransactionRequest) *Transaction {
	now := time.Now()
	return &Transaction{
		ID:          uuid.New().String(),
		UserID:      req.UserID,
		AccountID:   req.AccountID,
		Amount:      req.Amount,
		Currency:    req.Currency,
		Type:        req.Type,
		Status:      TransactionStatusPending,
		Description: req.Description,
		Reference:   req.Reference,
		Metadata:    make(map[string]interface{}),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
