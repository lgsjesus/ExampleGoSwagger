package database

import (
	"challenge.go.lgsjesus/domain"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

type Database struct {
	Db            *gorm.DB
	Dsn           string
	DbType        string
	AutoMigrateDb bool
	Env           string
}

func NewDb() *Database {
	return &Database{}
}

func (d *Database) Connect() (*gorm.DB, error) {
	var err error
	d.Db, err = gorm.Open(d.DbType, d.Dsn)
	if err != nil {
		return nil, err
	}
	if d.AutoMigrateDb {
		d.Db.AutoMigrate(&domain.Product{}, &domain.Customer{}, &domain.FavoriteProduct{}, &domain.User{})
	}

	return d.Db, nil

}
