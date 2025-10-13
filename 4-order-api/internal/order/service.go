package order

import (
	"4-order-api/internal/models"
	"4-order-api/pkg/di"
	"errors"
	"fmt"
)

type OrderService struct {
	UserRepository    di.IUserRepository
	ProductRepository di.IProductRepository
}

func NewOrderService(UserRepository di.IUserRepository, ProductRepository di.IProductRepository) *OrderService {
	return &OrderService{
		UserRepository:    UserRepository,
		ProductRepository: ProductRepository,
	}
}
func (service *OrderService) CreateOrderServ(productsResp []QuantProductID, userID uint, phoneNumber string) ([]*models.Product, error) {
	dbUser, err := service.UserRepository.FindUserByNum(phoneNumber)
	if err != nil || dbUser == nil {

		return nil, errors.New("ошибка поиска номера телефона")
	}
	fmt.Printf("DBUser %d\n", dbUser.ID)
	fmt.Printf("UserId %d\n", userID)
	if dbUser.ID != userID {

		return nil, errors.New("пользователь не найден")
	}
	productIDs := make([]uint, len(productsResp))
	for i, p := range productsResp {
		productIDs[i] = p.ProductID
	}

	products, err := service.ProductRepository.FindProductArrayById(productIDs)
	if err != nil {

		return nil, errors.New("ошибка поиска продукта")
	}
	return products, nil

}
