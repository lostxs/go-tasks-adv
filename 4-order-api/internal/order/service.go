package order

import (
	"4-order-api/internal/product"
)

type ServiceDeps struct {
	OrderRepository   *Repository
	ProductRepository *product.Repository
}

type Service struct {
	orderRepository   *Repository
	productRepository *product.Repository
}

func NewService(deps ServiceDeps) *Service {
	return &Service{
		orderRepository:   deps.OrderRepository,
		productRepository: deps.ProductRepository,
	}
}

func (s *Service) Create(userID, productId uint, quantity int) (*Order, error) {
	var orderItems []OrderItem

	p, err := s.productRepository.GetByID(productId)
	if err != nil {
		return nil, err
	}

	orderItems = append(orderItems, OrderItem{
		ProductID: p.ID,
		Price:     p.Price,
		Quantity:  quantity,
	})

	order := &Order{
		UserID: userID,
		Status: StatusCreated,
		Total:  p.Price * float64(quantity),
		Items:  orderItems,
	}

	return s.orderRepository.Create(order)
}

func (s *Service) GetByID(id uint) (*Order, error) {
	return s.orderRepository.GetByID(id)
}

func (s *Service) GetByUser(userID uint) ([]Order, error) {
	return s.orderRepository.GetByUserID(userID)
}
