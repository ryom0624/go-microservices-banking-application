package app

import (
	"banking/dto"
	"banking/mocks/service"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"local.packages/lib/errs"
	"net/http"
	"net/http/httptest"
	"testing"
)

var router *mux.Router
var ch CustomerHandler
var mockService *service.MockCustomerService

func setup(t *testing.T) func() {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService = service.NewMockCustomerService(ctrl)
	ch = CustomerHandler{service: mockService}
	router = mux.NewRouter()
	router.HandleFunc("/customers", ch.getAllCustomers)

	return func() {
		router = nil
		defer ctrl.Finish()
	}
}

func TestShouldReturnCustomersWithStatusCode200(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	dummyCustomers := []dto.CustomerResponse{
		{"1001", "Ashish", "New Delhi", "110011", "2000-01-01", "1"},
		{"1002", "Rob", "New Delhi", "110011", "2000-01-01", "1"},
	}

	mockService.EXPECT().GetAllCustomers("").Return(dummyCustomers, nil)
	request, _ := http.NewRequest(http.MethodGet, "/customers", nil)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Error("failed while testing the status code")
	}
}

func TestShouldReturnStatusCode500WithErrorCode(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	mockService.EXPECT().GetAllCustomers("").Return(nil, errs.NewUnexpectedError("some data"))

	request, _ := http.NewRequest(http.MethodGet, "/customers", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusInternalServerError {
		t.Error("failed while testing the status code")
	}

}
