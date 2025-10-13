package order

import (
	"time"

	"github.com/lib/pq"
)

type QuantProductID struct {
	ProductID uint `json:"product_id"`
	Quantity  uint `json:"quantity"`
}
type OrderRequest struct {
	UserID   uint             `json:"user_id"`
	Products []QuantProductID `json:"quantproduct_id"`
}

type OrderResponse struct {
	ID        uint              `json:"ID"`
	UserID    uint              `json:"user_id"`
	CreatedAt time.Time         `json:"CreatedAt"`
	UpdatedAt time.Time         `json:"UpdatedAt"`
	Products  []ProductResponse `json:"Products"`
}
type ProductResponse struct {
	ID          uint           `json:"ID"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Images      pq.StringArray `json:"images" gorm:"type:text[]"`
	Quantity    uint           `json:"quantity"` //количество в заказе

}
type AllGetOrderByUserResponse struct {
	Orders []OrderResponse
}
