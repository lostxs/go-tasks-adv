package product

import (
	"4-order-api/pkg/db"

	"gorm.io/gorm/clause"
)

type Repository struct {
	db *db.DB
}

func NewRepository(database *db.DB) *Repository {
	return &Repository{db: database}
}

func (r *Repository) Create(product *Product) (*Product, error) {
	result := r.db.Create(product)
	if result.Error != nil {
		return nil, result.Error
	}
	return product, nil
}

func (r *Repository) GetByID(id uint) (*Product, error) {
	var product Product
	result := r.db.First(&product, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &product, nil
}

func (r *Repository) Update(product *Product) (*Product, error) {
	result := r.db.Clauses(clause.Returning{}).Updates(product)
	if result.Error != nil {
		return nil, result.Error
	}
	return product, nil
}

func (r *Repository) Delete(id uint) error {
	result := r.db.Delete(&Product{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
