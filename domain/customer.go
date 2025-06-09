package domain

import "github.com/asaskevich/govalidator"

type Customer struct {
	ID      int    `gorm:"primaryKey;autoIncrement"`
	Name    string `valid:"required" gorm:"type:varchar(60)"`
	Email   string `valid:"required" gorm:"type:varchar(250)"`
	Phone   string `valid:"required" gorm:"type:varchar(20)"`
	Address string `valid:"required" gorm:"type:varchar(250)"`
}

func NewCustomer(name, email, phone, address string) *Customer {
	return &Customer{
		Name:    name,
		Email:   email,
		Phone:   phone,
		Address: address,
	}
}
func (c *Customer) Validate() error {
	_, err := govalidator.ValidateStruct(c)
	if err != nil {
		return err
	}
	return nil
}
