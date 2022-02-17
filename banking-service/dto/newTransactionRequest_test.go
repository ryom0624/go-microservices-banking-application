package dto

import (
	"net/http"
	"testing"
)

func TestShouldReturnErrorWhenTransactionTypeIsNotDepositOrWithdrawal(t *testing.T) {
	request := NewTransactionRequest{
		TransactionType: "invalid type",
	}
	appError := request.Validate()

	if appError.Message != "choose deposit or withdraw" {
		t.Errorf("invalid message while testing transaction type")
	}
	if appError.Code != http.StatusUnprocessableEntity {
		t.Errorf("invalid status code while testing transaction type")
	}
}

func TestShouldReturnErrorWhenAmountIsLessThanZero(t *testing.T) {
	request := NewTransactionRequest{Amount: -100, TransactionType: transactionTypeDeposit}

	appError := request.Validate()

	if appError.Message != "Amount cannot be less than zero" {
		t.Errorf("invalid message while validating amount")
	}
	if appError.Code != http.StatusUnprocessableEntity {
		t.Errorf("invalid status code while validating amount")
	}

}
