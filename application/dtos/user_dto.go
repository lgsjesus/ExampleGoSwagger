package dtos

import (
	"challenge.go.lgsjesus/domain"
	"github.com/asaskevich/govalidator"
	"golang.org/x/crypto/bcrypt"
)

type UserDto struct {
	ID       int    `json:"id" valid:"-"` // ID is not required for creation, so no validation
	UserName string `json:"username" valid:"required;stringlength(3|60)"`
	NickName string `json:"nickname" valid:"required;stringlength(3|20)"`
	Password string `json:"password" valid:"required;stringlength(3|16)"`
}

func (u *UserDto) Validate() error {
	_, err := govalidator.ValidateStruct(u)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserDto) MapToUser() (*domain.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	// Validate the UserDto before mapping
	if err := u.Validate(); err != nil {
		return nil, err
	}

	return &domain.User{
		Username: u.UserName,
		Password: string(hashedPassword),
		NickName: u.NickName,
	}, nil
}

func (d *UserDto) NewUserDto(u *domain.User) *UserDto {
	return &UserDto{
		UserName: u.Username,
		NickName: u.NickName, // Note: Password should not be exposed in DTOs in production
		ID:       u.ID,
	}
}
