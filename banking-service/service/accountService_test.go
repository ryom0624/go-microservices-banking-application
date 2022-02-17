package service

import (
	"banking/domain"
	"banking/dto"
	mockdomain "banking/mocks/domain"
	"github.com/golang/mock/gomock"
	"local.packages/lib/errs"
	"testing"
	"time"
)

const dbTSLayout = "2006-01-02 15:04:05"

var ctrl *gomock.Controller
var mockRepo *mockdomain.MockAccountRepository
var service AccountService

func setup(t *testing.T) func() {

	ctrl = gomock.NewController(t)
	mockRepo = mockdomain.NewMockAccountRepository(ctrl)
	service = NewAccountService(mockRepo)

	return func() {
		ctrl.Finish()
	}
}

func Test_should_return_a_validation_error_response_when_the_request_is_not_validated(t *testing.T) {
	request := dto.NewAccountRequest{
		CustomerID:  "100",
		AccountType: "saving",
		Amount:      0,
	}

	service := NewAccountService(nil)

	_, appError := service.NewAccount(request)

	if appError == nil {
		t.Error("failed while testing new account validation")
	}
}

func Test_should_return_an_error_from_the_server_side_if_the_new_account_cannot_be_created(t *testing.T) {

	teardown := setup(t)
	defer teardown()

	defer ctrl.Finish()

	request := dto.NewAccountRequest{
		CustomerID:  "100",
		AccountType: "saving",
		Amount:      6000,
	}

	account := domain.NewAccount(request.CustomerID, request.AccountType, request.Amount)

	mockRepo.EXPECT().Save(account).Return(nil, errs.NewUnexpectedError("Unexpected database error"))

	_, appError := service.NewAccount(request)
	if appError == nil {
		t.Error("test failed while validating error for new account")
	}

}

func Test_should_return_new_account_response_while_a_new_accunt_is_saved_successfully(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	defer ctrl.Finish()

	request := dto.NewAccountRequest{
		CustomerID:  "100",
		AccountType: "saving",
		Amount:      6000,
	}
	account := domain.Account{
		CustomerID:  request.CustomerID,
		OpeningDate: time.Now().Format(dbTSLayout),
		AccountType: request.AccountType,
		Amount:      request.Amount,
		Status:      "1",
	}
	accountWithId := account
	accountWithId.AccountID = "201"

	mockRepo.EXPECT().Save(account).Return(&accountWithId, nil)

	newAccount, appError := service.NewAccount(request)
	if appError != nil {
		t.Error("test failed while creating new account")
	}

	if newAccount.AccountID != accountWithId.AccountID {
		t.Error("failed while matching new accoun id")
	}

}
