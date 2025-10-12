package order

import (
	"4-order-api/pkg/db"
)

type Repository struct {
	db *db.DB
}

func NewRepository(database *db.DB) *Repository {
	return &Repository{db: database}
}

func (r *Repository) Create(order *Order) (*Order, error) {
	result := r.db.Create(order)
	if result.Error != nil {
		return nil, result.Error
	}
	return order, nil
}

func (r *Repository) GetByID(id uint) (*Order, error) {
	var order Order
	result := r.db.Preload("Items.Product").First(&order, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &order, nil
}

func (r *Repository) GetByUserID(userID uint) ([]Order, error) {
	var orders []Order
	result := r.db.Preload("Items.Product").Where("user_id = ?", userID).Find(&orders)
	if result.Error != nil {
		return nil, result.Error
	}
	return orders, nil
}
