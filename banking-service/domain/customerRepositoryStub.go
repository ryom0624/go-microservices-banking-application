package domain

import (
	"local.packages/lib/errs"
)

var _ CustomerRepository = (*CustomerRepositoryStub)(nil)

type CustomerRepositoryStub struct {
	customers []Customer
}

func (s CustomerRepositoryStub) FindAll(status string) ([]Customer, *errs.AppError) {
	return s.customers, nil
}

func (d CustomerRepositoryStub) ByID(id string) (*Customer, *errs.AppError) {
	if id != "1" {
		return nil, errs.NewNotFoundError("customer not found")
	}
	return &Customer{Id: "1001", Name: "RYO", City: "Tokyo", Zipcode: "1500001", DateOfBirth: "1992/02/12", Status: "1"}, nil

}

func NewCustomerRepositoryStub() CustomerRepositoryStub {
	return CustomerRepositoryStub{
		customers: []Customer{
			{Id: "1001", Name: "RYO", City: "Tokyo", Zipcode: "1500001", DateOfBirth: "1992/02/12", Status: "1"},
			{Id: "1002", Name: "Alice", City: "Tokyo", Zipcode: "1500001", DateOfBirth: "1992/02/12", Status: "1"},
			{Id: "1003", Name: "Bob", City: "Tokyo", Zipcode: "1500001", DateOfBirth: "1992/02/12", Status: "1"},
		},
	}
}
