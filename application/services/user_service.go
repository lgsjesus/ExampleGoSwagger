package services

import (
	"challenge.go.lgsjesus/application/dtos"
	"challenge.go.lgsjesus/framework/repositories"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repository *repositories.UserRepositoryDb
}

func NewUserService(repository *repositories.UserRepositoryDb) *UserService {
	return &UserService{
		repository: repository,
	}
}

func (s *UserService) CreateUser(userDto *dtos.UserDto) error {
	if err := userDto.Validate(); err != nil {
		return err
	}

	user, err := userDto.MapToUser()
	if err != nil {
		return err
	}
	_, err = s.repository.Insert(user)
	if err != nil {
		return err
	}
	return nil
}
func (s *UserService) GetUser(userId int) (*dtos.UserDto, error) {
	var userDto dtos.UserDto
	user, err := s.repository.Find(userId)
	if err != nil {
		return nil, err
	}
	return userDto.NewUserDto(user), nil
}
func (s *UserService) UpdateUser(userDto *dtos.UserDto) (*dtos.UserDto, error) {
	if err := userDto.Validate(); err != nil {
		return nil, err
	}

	user, err := userDto.MapToUser()
	if err != nil {
		return nil, err
	}
	user, err = s.repository.Update(user)
	if err != nil {
		return nil, err
	}
	return userDto.NewUserDto(user), nil
}
func checkPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
