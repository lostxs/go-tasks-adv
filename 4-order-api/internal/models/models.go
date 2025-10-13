package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Number    string `json:"number" gorm:"index"`
	SessionID string `json:"sessionId"`
	Code      string `json:"code"`
	Orders    []*Order
}

type Product struct {
	gorm.Model
	Name         string         `json:"name"`
	Description  string         `json:"description"`
	Images       pq.StringArray `json:"images" gorm:"type:text[]"`
	Orders       []*Order       `gorm:"many2many:order_products;"`
	Quantity     uint           `json:"quantity"`
	OrderProduct []OrderProduct `gorm:"foreignKey:ProductID"`
}

type Order struct {
	gorm.Model
	UserId   uint       `json:"user_id"`
	Products []*Product `gorm:"many2many:order_products;"`
}

// промежуточная таблица
type OrderProduct struct {
	OrderID   uint `gorm:"primaryKey"`
	ProductID uint `gorm:"primaryKey"`
	Quantity  uint `json:"quantity"`
}

func NewProduct(name, des string, images []string, quantity uint) *Product {
	return &Product{
		Name:        name,
		Description: des,
		Images:      pq.StringArray(images),
		Quantity:    quantity,
	}
}
