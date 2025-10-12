package order

type CreateOrderRequest struct {
	ProductID uint `json:"product_id"`
	Quantity  int  `json:"quantity"`
}
