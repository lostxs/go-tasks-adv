package product

import (
	"4-order-api/pkg/request"
	"4-order-api/pkg/response"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

type HandlerDeps struct {
	ProductRepository *Repository
}

type Handler struct {
	productRepository *Repository
}

func NewHandler(router *http.ServeMux, deps HandlerDeps) {
	hander := &Handler{
		productRepository: deps.ProductRepository,
	}

	router.HandleFunc("POST /product", hander.Create())
	router.HandleFunc("GET /product/{id}", hander.Get())
	router.HandleFunc("PATCH /product/{id}", hander.Update())
	router.HandleFunc("DELETE /product/{id}", hander.Delete())
}

func (h *Handler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := request.ParseBody[CreateRequest](w, r)
		if err != nil {
			return
		}
		createdProuct, err := h.productRepository.Create(&Product{
			Name:        body.Name,
			Description: body.Description,
			Images:      body.Images,
			Price:       body.Price,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		response.WriteJSON(w, http.StatusCreated, createdProuct)
	}
}

func (h *Handler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		product, err := h.productRepository.GetByID(uint(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		response.WriteJSON(w, http.StatusOK, product)
	}
}

func (h *Handler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := request.ParseBody[UpdateRequest](w, r)
		if err != nil {
			return
		}
		idStr := r.PathValue("id")
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		_, err = h.productRepository.GetByID(uint(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		updatedProduct, err := h.productRepository.Update(&Product{
			Model:       gorm.Model{ID: uint(id)},
			Name:        body.Name,
			Description: body.Description,
			Images:      body.Images,
			Price:       body.Price,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		response.WriteJSON(w, http.StatusOK, updatedProduct)
	}
}

func (h *Handler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = h.productRepository.Delete(uint(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		response.WriteJSON(w, http.StatusNoContent, nil)
	}
}
