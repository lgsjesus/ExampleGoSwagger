package dtos

import (
	"challenge.go.lgsjesus/domain"
	"github.com/asaskevich/govalidator"
)

type CustomerDto struct {
	ID      int    `json:"id" valid:"-"`
	Name    string `json:"name" valid:"required"`
	Email   string `json:"email" valid:"required"`
	Phone   string `json:"phone" valid:"required"`
	Address string `json:"address" valid:"required"`
}

func (c *CustomerDto) Validate() error {
	_, err := govalidator.ValidateStruct(c)
	if err != nil {
		return err
	}
	return nil
}

func (c *CustomerDto) MapToCustomer() *domain.Customer {
	return &domain.Customer{
		ID:      c.ID,
		Name:    c.Name,
		Email:   c.Email,
		Phone:   c.Phone,
		Address: c.Address,
	}
}

func (d *CustomerDto) NewCustomerDto(c *domain.Customer) *CustomerDto {
	return &CustomerDto{
		ID:      c.ID,
		Name:    c.Name,
		Email:   c.Email,
		Phone:   c.Phone,
		Address: c.Address,
	}
}
