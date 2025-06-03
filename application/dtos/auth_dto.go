package dtos

import "github.com/asaskevich/govalidator"

type AuthDto struct {
	NickName string `json:"nickname" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (a *AuthDto) Validate() error {
	_, err := govalidator.ValidateStruct(a)
	if err != nil {
		return err
	}
	return nil
}
