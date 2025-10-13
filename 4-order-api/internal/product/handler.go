package product

import (
	"4-order-api/configs"
	"4-order-api/internal/models"
	"4-order-api/pkg/middleware"
	"4-order-api/pkg/req"
	"4-order-api/pkg/resp"
	"fmt"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

type ProductHandDeps struct {
	ProductRepository *ProductRepository
	*configs.Config
}
type ProductHandler struct {
	ProductRepository *ProductRepository
}

func NewProductHandler(router *http.ServeMux, product *ProductHandDeps) {
	handler := &ProductHandler{
		ProductRepository: product.ProductRepository,
	}
	router.HandleFunc("POST /prod/create", handler.create)

	router.HandleFunc("PATCH /prod/update/{id}", handler.update)
	router.HandleFunc("DELETE /prod/delete/{id}", handler.delete)
	router.HandleFunc("GET /prod/{id}", middleware.Auth(handler.getById, product.Config))
	router.HandleFunc("GET /all", middleware.Auth(handler.getAllProduct, product.Config))
}
func (handler *ProductHandler) create(w http.ResponseWriter, request *http.Request) {
	body, err := req.HandleBody[ProductCreate](w, request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	product := models.NewProduct(body.Name, body.Description, body.Images, body.Quantity)
	createProd, err := handler.ProductRepository.Create(product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

	}
	resp.Json(w, createProd, 201)

}
func (handler *ProductHandler) update(w http.ResponseWriter, request *http.Request) {
	body, err := req.HandleBody[ProductUpdate](w, request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	idStr := request.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	product, err := handler.ProductRepository.Update(&models.Product{
		Model:       gorm.Model{ID: uint(id)},
		Name:        body.Name,
		Description: body.Description,
		Images:      body.Images,
		Quantity:    body.Quantity,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	resp.Json(w, product, 201)

}
func (handler *ProductHandler) delete(w http.ResponseWriter, request *http.Request) {
	idStr := request.PathValue("id")
	id, _ := strconv.Atoi(idStr)

	_, err := handler.ProductRepository.FindId(uint(id))
	if err != nil {
		resp.Json(w, "Карточка не найдена", http.StatusOK)
		return
	}
	err = handler.ProductRepository.Delete(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp.Json(w, "Карточка удалена", http.StatusOK)
}

func (handler *ProductHandler) getById(w http.ResponseWriter, request *http.Request) {
	phonNumber, ok := request.Context().Value(middleware.ContextPhoneNumber).(string)
	if ok {
		//редирект на страничку register
		fmt.Println(phonNumber)
	}
	idStr := request.PathValue("id")
	id, _ := strconv.Atoi(idStr)
	getProduct, err := handler.ProductRepository.FindId(uint(id))
	if err != nil {
		resp.Json(w, "Карточка не найдена", http.StatusOK)
		return
	}
	resp.Json(w, getProduct, http.StatusOK)

}
func (handler *ProductHandler) getAllProduct(w http.ResponseWriter, request *http.Request) {
	phonNumber, ok := request.Context().Value(middleware.ContextPhoneNumber).(string)
	if ok {
		fmt.Println(phonNumber)
	}

	allprod, err := handler.ProductRepository.GetAllProd()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	resp.Json(w, allprod, http.StatusOK)
}
