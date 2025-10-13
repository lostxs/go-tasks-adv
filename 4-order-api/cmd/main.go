package main

import (
	"4-order-api/configs"
	"4-order-api/internal/auth"
	"4-order-api/internal/order"
	"4-order-api/internal/product"
	"4-order-api/internal/user"
	"4-order-api/pkg/db"
	"4-order-api/pkg/logger"
	"4-order-api/pkg/middleware"
	"fmt"
	"net/http"
)

func App() *http.ServeMux {
	logger.LogInit()
	router := http.NewServeMux()
	conf := configs.LoadConfig()
	db := db.NewDb(conf)
	userRepository := user.NewUserRepository(db)
	productRepository := product.NewProductRepository(db)
	OrderRepository := order.NewOrderRepository(db)

	authService := auth.NewAuthRepository(userRepository)
	orderService := order.NewOrderService(userRepository, productRepository)
	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		Config:      conf,
		AuthService: authService,
	})
	product.NewProductHandler(router, &product.ProductHandDeps{
		ProductRepository: productRepository,
		Config:            conf,
	})
	order.NewOrderHandler(router, order.OrderHandlerDeps{
		Config:          conf,
		OrderRepository: OrderRepository,
		OrderService:    orderService,
	})
	return router
}

func main() {
	app := App()
	server := http.Server{
		Addr:    ":8085",
		Handler: middleware.Logging(app),
	}
	fmt.Println("БД Создана")
	fmt.Printf("Listen port%v\n", server.Addr)
	server.ListenAndServe()

}
