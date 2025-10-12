package order

import (
	"4-order-api/internal/product"

	"gorm.io/gorm"
)

type Status string

const (
	StatusCreated Status = "created"
)

type Order struct {
	gorm.Model
	UserID uint
	Status Status `gorm:"type:varchar(20);default:'created'"`
	Total  float64
	Items  []OrderItem `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type OrderItem struct {
	gorm.Model
	OrderID   uint
	ProductID uint
	Product   product.Product `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Quantity  int
	Price     float64
}
