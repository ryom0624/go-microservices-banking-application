package service

import (
	"banking/domain"
	"banking/dto"
	"local.packages/lib/errs"
	"time"
)

type AccountService interface {
	NewAccount(dto.NewAccountRequest) (*dto.NewAccountResponse, *errs.AppError)
	MakeTransaction(dto.NewTransactionRequest) (*dto.NewTransactionResponse, *errs.AppError)
}

var _ AccountService = (*DefaultAccountService)(nil)

type DefaultAccountService struct {
	repo domain.AccountRepository
}

func (s DefaultAccountService) NewAccount(req dto.NewAccountRequest) (*dto.NewAccountResponse, *errs.AppError) {
	err := req.Validate()
	if err != nil {
		return nil, err
	}

	a := domain.Account{
		AccountID:   "",
		CustomerID:  req.CustomerID,
		OpeningDate: time.Now().Format("2006-01-02 15:04:05"),
		AccountType: req.AccountType,
		Amount:      req.Amount,
		Status:      "1",
	}

	newAccount, err := s.repo.Save(a)
	if err != nil {
		return nil, err
	}

	response := newAccount.ToNewAccountResponseDto()

	return &response, nil
}

func (s DefaultAccountService) MakeTransaction(req dto.NewTransactionRequest) (*dto.NewTransactionResponse, *errs.AppError) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	if req.IsTransactionTypeWithdrawal() {
		account, err := s.repo.Find(req.AccountID)
		if err != nil {
			return nil, err
		}
		if !account.CanWithDraw(req.Amount) {
			return nil, errs.NewValidationError("insufficient balance in the account")
		}
	}

	transaction := domain.Transaction{
		AccountID:       req.AccountID,
		Amount:          req.Amount,
		TransactionType: req.TransactionType,
		TransactionDate: time.Now().Format("2006-01-02 15:04:05"),
	}
	savedTransaction, appError := s.repo.SaveTransaction(transaction)
	if appError != nil {
		return nil, appError
	}

	response := savedTransaction.ToTransactionResponseDto()

	return &response, nil
}

func NewAccountService(repo domain.AccountRepository) AccountService {
	return DefaultAccountService{repo: repo}
}
