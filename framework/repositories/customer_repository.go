package repositories

import (
	"errors"

	"challenge.go.lgsjesus/domain"
	"github.com/jinzhu/gorm"
)

type CustomerRepository interface {
	Insert(customer *domain.Customer) (*domain.Customer, error)
	Find(id int) (*domain.Customer, error)
	Update(customer *domain.Customer) (*domain.Customer, error)
	//Delete(id int) error
	List() ([]*domain.Customer, error)
}

type CustomerRepositoryDb struct {
	Db *gorm.DB
}

// NewCustomerRepositoryDb creates a new instance of CustomerRepositoryDb
func NewCustomerRepositoryDb(db *gorm.DB) *CustomerRepositoryDb {
	return &CustomerRepositoryDb{Db: db}
}

func (repo CustomerRepositoryDb) Insert(customer *domain.Customer) (*domain.Customer, error) {
	error := customer.Validate()
	if error != nil {
		return nil, error
	}
	err := repo.Db.Create(customer).Error
	if err != nil {
		return nil, err
	}
	return customer, nil
}
func (repo CustomerRepositoryDb) Find(id int) (*domain.Customer, error) {
	var customer domain.Customer
	err := repo.Db.First(&customer, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	if customer.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &customer, nil
}
func (repo CustomerRepositoryDb) Update(customer *domain.Customer) (*domain.Customer, error) {
	var err error
	err = customer.Validate()
	if err != nil {
		return nil, err
	}
	if customer.ID == 0 {
		return nil, errors.New("customer ID cannot be zero")
	}
	err = repo.Db.Save(customer).Error
	if err != nil {
		return nil, err
	}
	return customer, nil
}

func (repo CustomerRepositoryDb) List() ([]*domain.Customer, error) {
	var customers []*domain.Customer
	err := repo.Db.Find(&customers).Error
	if err != nil {
		return nil, err
	}
	if len(customers) == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return customers, nil
}
