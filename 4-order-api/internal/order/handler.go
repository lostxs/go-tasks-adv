package order

import (
	"4-order-api/pkg/middleware"
	"4-order-api/pkg/request"
	"4-order-api/pkg/response"
	"net/http"
	"strconv"
)

type HandlerDeps struct {
	OrderService *Service
	AuthMW       *middleware.JWTAuthMiddleware
}

type Handler struct {
	orderService *Service
	authMW       *middleware.JWTAuthMiddleware
}

func NewHandler(router *http.ServeMux, deps HandlerDeps) {
	h := &Handler{
		orderService: deps.OrderService,
		authMW:       deps.AuthMW,
	}

	router.Handle("POST /order", h.authMW.Handler(http.HandlerFunc(h.Create())))
	router.Handle("GET /order/{id}", h.authMW.Handler(http.HandlerFunc(h.GetByID())))
	router.Handle("GET /my-orders", h.authMW.Handler(http.HandlerFunc(h.GetMyOrders())))
}

func (h *Handler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := h.authMW.GetUser(r.Context())

		body, err := request.ParseBody[CreateOrderRequest](w, r)
		if err != nil {
			return
		}

		order, err := h.orderService.Create(claims.UserID, body.ProductID, body.Quantity)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response.WriteJSON(w, http.StatusCreated, order)
	}

}

func (h *Handler) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := h.authMW.GetUser(r.Context())

		idStr := r.PathValue("id")
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		order, err := h.orderService.GetByID(uint(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		if order.UserID != claims.UserID {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		response.WriteJSON(w, http.StatusOK, order)
	}
}

func (h *Handler) GetMyOrders() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := h.authMW.GetUser(r.Context())
		if claims == nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		orders, err := h.orderService.GetByUser(claims.UserID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response.WriteJSON(w, http.StatusOK, orders)
	}
}
