package main

import (
	"4-order-api/config"
	"4-order-api/internal/auth"
	"4-order-api/internal/order"
	"4-order-api/internal/product"
	"4-order-api/internal/session"
	"4-order-api/internal/user"
	"4-order-api/pkg/db"
	"4-order-api/pkg/jwt"
	"4-order-api/pkg/middleware"
	"net/http"

	"github.com/sirupsen/logrus"
)

func main() {
	config := config.Load()
	db := db.NewDB(config.DB)
	router := http.NewServeMux()
	logrus.SetFormatter(&logrus.JSONFormatter{})
	jwt := jwt.NewJWT(config.Auth.Secret)

	authMW := middleware.NewJWTAuthMiddleware(jwt)

	sessionRepository := session.NewRepository()
	proudctRepository := product.NewRepository(db)
	userRepository := user.NewRepository(db)
	orderRepository := order.NewRepository(db)

	authService := auth.NewService(auth.ServiceDeps{
		SessionRepository: sessionRepository,
		UserRepository:    userRepository,
		Jwt:               jwt,
	})
	orderService := order.NewService(order.ServiceDeps{
		OrderRepository:   orderRepository,
		ProductRepository: proudctRepository,
	})

	auth.NewHandler(router, auth.HandlerDeps{
		AuthService: authService,
	})
	product.NewHandler(router, product.HandlerDeps{
		ProductRepository: proudctRepository,
	})
	order.NewHandler(router, order.HandlerDeps{
		OrderService: orderService,
		AuthMW:       authMW,
	})

	server := http.Server{
		Addr:    ":8080",
		Handler: middleware.Logging(router),
	}

	logrus.Info("Starting server on :8080")
	server.ListenAndServe()
}
