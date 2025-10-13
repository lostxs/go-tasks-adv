package product

import (
	"4-order-api/internal/models"
	"4-order-api/pkg/db"
	"errors"
)

type ProductRepository struct {
	Database *db.Db
}

func NewProductRepository(database *db.Db) *ProductRepository {
	return &ProductRepository{
		Database: database,
	}
}
func (repos *ProductRepository) Create(product *models.Product) (*models.Product, error) {
	result := repos.Database.DB.Create(product)
	if result.Error != nil {
		return nil, result.Error
	}
	return product, nil
}
func (repos *ProductRepository) Update(product *models.Product) (*models.Product, error) {
	result := repos.Database.DB.Updates(product)
	if result.Error != nil {
		return nil, result.Error
	}
	return product, nil
}
func (repos *ProductRepository) FindId(id uint) (*models.Product, error) {
	var prod models.Product

	res := repos.Database.DB.First(&prod, id)
	if res.Error != nil {
		return nil, res.Error
	}
	return &prod, nil
}
func (repos *ProductRepository) Delete(id string) error {
	res := repos.Database.DB.Delete(&models.Product{}, id)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
func (repos *ProductRepository) GetId(id string) error {
	var prod models.Product
	res := repos.Database.DB.First(&prod, id)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
func (repos *ProductRepository) GetAllProd() ([]models.Product, error) {
	var allProd []models.Product
	res := repos.Database.DB.Find(&allProd)
	if res.Error != nil {
		return nil, res.Error
	}
	return allProd, nil
}
func (repo *ProductRepository) FindProductArrayById(productIDs []uint) ([]*models.Product, error) {
	var products []*models.Product

	if err := repo.Database.DB.Debug().Where("id IN ?", productIDs).Find(&products).Error; err != nil {

		return nil, errors.New("ошибка добавления продуктов")
	}
	if len(products) == 0 {
		return nil, errors.New("продукты не найдены")
	}

	return products, nil
}
