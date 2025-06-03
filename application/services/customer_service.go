package services

import (
	"challenge.go.lgsjesus/application/dtos"
	"challenge.go.lgsjesus/framework/repositories"
)

type CustomerService struct {
	repository *repositories.CustomerRepositoryDb
}

func NewCustomerService(repository *repositories.CustomerRepositoryDb) *CustomerService {
	return &CustomerService{
		repository: repository,
	}
}

func (s *CustomerService) CreateCustomer(customerDto *dtos.CustomerDto) (*dtos.CustomerDto, error) {
	if err := customerDto.Validate(); err != nil {
		return nil, err
	}

	customer := customerDto.MapToCustomer()
	customer, err := s.repository.Insert(customer)
	if err != nil {
		return nil, err
	}

	return customerDto.NewCustomerDto(customer), nil
}

func (s *CustomerService) GetCustomer(customerId int) (*dtos.CustomerDto, error) {

	var customerDto dtos.CustomerDto
	customer, err := s.repository.Find(customerId)
	if err != nil {
		return nil, err
	}
	return customerDto.NewCustomerDto(customer), nil
}

func (s *CustomerService) UpdateCustomer(customerDto *dtos.CustomerDto) (*dtos.CustomerDto, error) {
	if err := customerDto.Validate(); err != nil {
		return nil, err
	}

	customer := customerDto.MapToCustomer()
	customer, err := s.repository.Update(customer)
	if err != nil {
		return nil, err
	}
	return customerDto.NewCustomerDto(customer), nil
}
