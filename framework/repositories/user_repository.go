package repositories

import (
	"errors"

	"challenge.go.lgsjesus/domain"
	"github.com/jinzhu/gorm"
)

type UserRepository interface {
	Insert(user *domain.User) (*domain.User, error)
	Find(id int) (*domain.User, error)
	FindByNickName(nick_name string) (*domain.User, error)
	Update(user *domain.User) (*domain.User, error)
}

type UserRepositoryDb struct {
	Db *gorm.DB
}

// NewUserRepositoryDb creates a new instance of UserRepositoryDb
func NewUserRepositoryDb(db *gorm.DB) *UserRepositoryDb {
	return &UserRepositoryDb{Db: db}
}

func (repo UserRepositoryDb) Insert(user *domain.User) (*domain.User, error) {
	err := user.Validate()
	if err != nil {
		return nil, err
	}
	err = repo.verifyDuplicateNickName(user)
	if err != nil {
		return nil, err
	}

	err = repo.Db.Create(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (repo UserRepositoryDb) Find(id int) (*domain.User, error) {
	var user domain.User
	err := repo.Db.First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	if user.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &user, nil
}
func (repo UserRepositoryDb) Update(user *domain.User) (*domain.User, error) {
	var err error
	err = user.Validate()
	if err != nil {
		return nil, err
	}
	var existingUser domain.User
	repo.Db.Find(&existingUser, "id = ?", user.ID)
	if existingUser.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	err = repo.verifyDuplicateNickName(user)
	if err != nil {
		return nil, err
	}

	err = repo.Db.Save(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (repo UserRepositoryDb) FindByNickName(nickName string) (*domain.User, error) {
	var user domain.User
	err := repo.Db.Where("nick_name = ?", nickName).First(&user).Error
	if err != nil {
		return nil, err
	}
	if user.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &user, nil
}

func (repo UserRepositoryDb) verifyDuplicateNickName(user *domain.User) error {
	var existingUser domain.User
	err := repo.Db.Where("nick_name = ? and id <> ?", user.NickName, user.ID).First(&existingUser).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if existingUser.ID != 0 && existingUser.ID != user.ID {
		return errors.New("already user with this nickname, please choose another nickname") // or a custom error indicating duplicate nickname
	}
	return nil
}
