package order

import (
	"4-order-api/configs"
	"4-order-api/internal/models"
	"4-order-api/pkg/middleware"
	"4-order-api/pkg/req"
	"4-order-api/pkg/resp"
	"fmt"
	"net/http"
	"strconv"
)

// сделать интерфейс для того чтобы увести зависимости
type OrderHandlerDeps struct {
	*OrderRepository
	*configs.Config
	*OrderService
}
type OrderHandler struct {
	*configs.Config
	*OrderRepository
	*OrderService
}

func NewOrderHandler(router *http.ServeMux, deps OrderHandlerDeps) {
	handler := &OrderHandler{
		Config:          deps.Config,
		OrderRepository: deps.OrderRepository,
		OrderService:    deps.OrderService,
	}
	router.HandleFunc("POST /order", middleware.Auth(handler.Order, deps.Config))
	router.HandleFunc("GET /order/{id}", middleware.Auth(handler.getOrder, deps.Config))
	router.HandleFunc("GET /my-orders/{id}", middleware.Auth(handler.getOrderByUser, deps.Config))
}
func (handler *OrderHandler) Order(w http.ResponseWriter, request *http.Request) {
	body, err := req.HandleBody[OrderRequest](w, request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//либо реализовать вместо поиска по ИД в юзер репозитории
	phonNumber, ok := request.Context().Value(middleware.ContextPhoneNumber).(string)
	if !ok {
		fmt.Println(phonNumber)
	}

	products, err := handler.CreateOrderServ(body.Products, body.UserID, phonNumber)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	orderResp, err := handler.OrderRepository.CreateOrder(body.UserID, products, body.Products)
	fmt.Println(orderResp.ID)

	newOrder := GetResponseOrder(orderResp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	productMap := make(map[uint]*ProductResponse)
	for i := range newOrder.Products {
		productMap[newOrder.Products[i].ID] = &newOrder.Products[i]
	}

	for _, reqProduct := range body.Products {
		if existing, ok := productMap[reqProduct.ProductID]; ok {

			existing.Quantity = reqProduct.Quantity

		}
	}
	fmt.Println(newOrder)

	resp.Json(w, newOrder, 201)
}

func (handler *OrderHandler) getOrder(w http.ResponseWriter, request *http.Request) {
	idStr := request.PathValue("id")
	resId, _ := strconv.Atoi(idStr)

	getOrders, err := handler.OrderRepository.GetOrder(uint(resId))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	order := GetResponseOrder(getOrders)
	resp.Json(w, order, http.StatusOK)
}
func (handler *OrderHandler) getOrderByUser(w http.ResponseWriter, request *http.Request) {
	idStr := request.PathValue("id")
	resId, _ := strconv.Atoi(idStr)
	getOrders, err := handler.OrderRepository.FindOrderId(uint(resId))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	orders := AllGetOrderByUserResponse{
		Orders: make([]OrderResponse, 0, len(getOrders)),
	}
	for _, order := range getOrders {

		orders.Orders = append(orders.Orders, GetResponseOrder(order))
	}

	resp.Json(w, orders, http.StatusOK)
}
func GetResponseOrder(getOrders *models.Order) OrderResponse {
	order := OrderResponse{
		ID:        getOrders.ID,
		UserID:    getOrders.UserId,
		CreatedAt: getOrders.CreatedAt,
		UpdatedAt: getOrders.UpdatedAt,
		Products:  make([]ProductResponse, 0, len(getOrders.Products)),
	}
	for _, product := range getOrders.Products {

		var quantity uint
		for _, q := range product.OrderProduct {
			if order.ID == q.OrderID {
				quantity = q.Quantity
				break
			}
		}
		order.Products = append(order.Products, ProductResponse{
			ID:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Images:      product.Images,
			Quantity:    quantity,
		})

	}
	return order
}
