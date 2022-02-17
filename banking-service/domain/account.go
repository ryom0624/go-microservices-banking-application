package domain

import (
	"banking/dto"
	"local.packages/lib/errs"
	"time"
)

type Account struct {
	AccountID   string  `db:"account_id"`
	CustomerID  string  `db:"customer_id"`
	OpeningDate string  `db:"opening_date"`
	AccountType string  `db:"account_type"`
	Amount      float64 `db:"amount"`
	Status      string  `db:"status"`
}

func (a Account) ToNewAccountResponseDto() *dto.NewAccountResponse {
	return &dto.NewAccountResponse{AccountID: a.AccountID}
}

func (a Account) CanWithDraw(amount float64) bool {
	return a.Amount > amount
}

//go:generate mockgen -destination=../mocks/domain/mockAccountRepository.go -package=domain banking/domain AccountRepository
type AccountRepository interface {
	Save(Account) (*Account, *errs.AppError)
	SaveTransaction(transaction Transaction) (*Transaction, *errs.AppError)
	Find(accountID string) (*Account, *errs.AppError)
}

const dbTSLayout = "2006-01-02 15:04:05"

func NewAccount(customerID, accountType string, amount float64) Account {
	return Account{
		CustomerID:  customerID,
		AccountType: accountType,
		Amount:      amount,
		OpeningDate: time.Now().Format(dbTSLayout),
		Status:      "1",
	}
}
