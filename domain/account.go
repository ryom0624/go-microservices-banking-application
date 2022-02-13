package domain

import (
	"go-microservices-banking-application/dto"
	"go-microservices-banking-application/errs"
)

type Account struct {
	AccountID   string  `db:"account_id"`
	CustomerID  string  `db:"customer_id"`
	OpeningDate string  `db:"opening_date"`
	AccountType string  `db:"account_type"`
	Amount      float64 `db:"amount"`
	Status      string  `db:"status"`
}

func (a Account) ToNewAccountResponseDto() dto.NewAccountResponse {
	return dto.NewAccountResponse{AccountID: a.AccountID}
}

func (a Account) CanWithDraw(amount float64) bool {
	return a.Amount > amount
}

type AccountRepository interface {
	Save(Account) (*Account, *errs.AppError)
	SaveTransaction(transaction Transaction) (*Transaction, *errs.AppError)
	Find(accountID string) (*Account, *errs.AppError)
}
