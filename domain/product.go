package domain

import "github.com/asaskevich/govalidator"

type Product struct {
	ID     int     `gorm:"primaryKey;autoIncrement"`
	Title  string  `valid:"not_null" gorm:"type:varchar(150)"`
	Review *string `gorm:"type:varchar(300)"`
	Image  *string `gorm:"type:varchar(400)"`
	Price  float64 `valid:"not_null" gorm:"type:decimal(10,2)"`
}

func NewProduct(title string, review *string, image *string, price float64) *Product {
	return &Product{
		Title:  title,
		Review: review,
		Image:  image,
		Price:  price,
	}
}
func (p *Product) Validate() error {
	_, err := govalidator.ValidateStruct(p)
	if err != nil {
		return err
	}
	return nil
}
