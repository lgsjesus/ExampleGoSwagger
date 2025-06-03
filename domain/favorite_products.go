package domain

import "github.com/asaskevich/govalidator"

type FavoriteProduct struct {
	ID         int `gorm:"primaryKey;autoIncrement"`
	ProductID  int `gorm:"not null column:product_id;"`
	CustomerID int `gorm:"not null column:customer_id;"`
}

func NewFavoriteProduct(productID, customerID int) *FavoriteProduct {
	return &FavoriteProduct{
		ProductID:  productID,
		CustomerID: customerID,
	}
}
func (fp *FavoriteProduct) Validate() error {
	_, error := govalidator.ValidateStruct(fp)
	if error != nil {
		return error
	}
	return nil
}
