package dto

import (
	"local.packages/errs"
	"strings"
)

type NewTransactionRequest struct {
	CustomerID      string  `json:"customer_id"`
	AccountID       string  `json:"account_id"`
	Amount          float64 `json:"amount"`
	TransactionType string  `json:"transaction_type"`
}

const (
	transactionTypeWithdrawal = "withdraw"
	transactionTypeDeposit    = "deposit"
)

func (r NewTransactionRequest) IsTransactionTypeWithdrawal() bool {
	return strings.ToLower(r.TransactionType) == transactionTypeWithdrawal
}

func (r NewTransactionRequest) IsTransactionTypeDeposit() bool {
	return strings.ToLower(r.TransactionType) == transactionTypeDeposit
}

func (r NewTransactionRequest) Validate() *errs.AppError {
	if strings.ToLower(r.TransactionType) != "withdraw" && strings.ToLower(r.TransactionType) != "deposit" {
		return errs.NewValidationError("choose deposit or withdraw")
	}
	if r.Amount < 0 {
		return errs.NewValidationError("Amount cannot be less than zero")
	}
	return nil
}
