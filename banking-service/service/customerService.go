package service

import (
	"banking/domain"
	"banking/dto"
	"local.packages/lib/errs"
)

const (
	statusInActive = "0"
	statusActive   = "1"
)

//go:generate mockgen -destination=../mocks/service/mockCustomerService.go -package=service banking/service CustomerService
type CustomerService interface {
	GetAllCustomers(string) ([]dto.CustomerResponse, *errs.AppError)
	GetCustomer(string) (*dto.CustomerResponse, *errs.AppError)
}

type DefaultCustomerService struct {
	repo domain.CustomerRepository
}

var _ CustomerService = (*DefaultCustomerService)(nil)

func (s DefaultCustomerService) GetAllCustomers(status string) ([]dto.CustomerResponse, *errs.AppError) {
	var statusQuery string
	if status == "inactive" {
		statusQuery = statusInActive
	} else if status == "active" {
		statusQuery = statusActive
	} else {
		statusQuery = ""
	}

	customers, err := s.repo.FindAll(statusQuery)
	if err != nil {
		return nil, err
	}

	response := make([]dto.CustomerResponse, 0)
	for _, c := range customers {
		response = append(response, c.ToDto())
	}

	return response, nil
}
func (s DefaultCustomerService) GetCustomer(id string) (*dto.CustomerResponse, *errs.AppError) {
	c, err := s.repo.ByID(id)
	if err != nil {
		return nil, err
	}

	response := c.ToDto()

	return &response, nil
}

func NewCustomerService(repository domain.CustomerRepository) DefaultCustomerService {
	return DefaultCustomerService{repo: repository}
}
