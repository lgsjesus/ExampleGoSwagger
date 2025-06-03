package domain

import "github.com/asaskevich/govalidator"

type User struct {
	ID       int    `gorm:"primaryKey;autoIncrement"`
	Username string `valid:"required" gorm:"type:varchar(60)"`
	NickName string `valid:"required" gorm:"type:varchar(30);uniqueIndex"`
	Password string `valid:"required" gorm:"type:varchar(250)"`
}

func NewUser(username, nickname, password string) *User {
	return &User{
		Username: username,
		NickName: nickname,
		Password: password,
	}
}
func (u *User) Validate() error {
	_, err := govalidator.ValidateStruct(u)
	if err != nil {
		return err
	}
	return nil
}
