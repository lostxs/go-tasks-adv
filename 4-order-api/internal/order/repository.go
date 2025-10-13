package order

import (
	"4-order-api/internal/models"
	"4-order-api/pkg/db"

	"errors"

	"gorm.io/gorm"
)

type OrderRepository struct {
	Database *db.Db
}

func NewOrderRepository(database *db.Db) *OrderRepository {
	return &OrderRepository{
		Database: database,
	}
}
func (repo *OrderRepository) CreateOrder(bodyUser uint, products []*models.Product, quantProd []QuantProductID) (*models.Order, error) {
	orderResp := &models.Order{
		UserId:   bodyUser,
		Products: products,
	}
	repo.Database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Debug().Create(orderResp).Error; err != nil {
			return err
		}

		for _, prodquant := range quantProd {
			var product models.Product
			if err := tx.First(&product, prodquant.ProductID).Error; err != nil {
				return errors.New("товар не найден")
			}
			if product.Quantity < prodquant.Quantity {
				return errors.New("недостаточно товара")
			}
			orderProduct := models.OrderProduct{
				OrderID:   orderResp.ID,
				ProductID: prodquant.ProductID,
				Quantity:  prodquant.Quantity,
			}

			if err := tx.Debug().Save(&orderProduct).Error; err != nil {
				return err
			}

		}
		return nil
	})

	return orderResp, nil
}

func (repo *OrderRepository) FindOrderId(userId uint) ([]*models.Order, error) {
	var orders []*models.Order

	err := repo.Database.DB.Debug().
		Where("user_id = ?", userId).
		Preload("Products").
		Preload("Products.OrderProduct", func(db *gorm.DB) *gorm.DB {
			return db.Where("order_id IN (?)", repo.Database.DB.
				Table("orders").
				Select("id").
				Where("user_id = ?", userId))
		}).
		Find(&orders).Error

	if err != nil {
		return nil, err
	}
	if len(orders) == 0 {
		return nil, errors.New("заказов не найдено")
	}

	return orders, nil
}
func (repo *OrderRepository) GetOrder(orderId uint) (*models.Order, error) {
	var order models.Order
	result := repo.Database.DB.Debug().
		Where("id =?", orderId).
		Preload("Products").
		Preload("Products.OrderProduct", "order_id =?", orderId).
		First(&order)
	if result.Error != nil {
		return nil, result.Error
	}
	return &order, nil
}
