package main

import (
	"4-order-api/config"
	"4-order-api/internal/auth"
	"4-order-api/internal/product"
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

	proudctRepo := product.NewRepository(db)
	authRepo := auth.NewRepository()

	authService := auth.NewService(authRepo, jwt)

	product.NewHandler(router, product.HandlerDeps{
		Repo: proudctRepo,
	})
	auth.NewHandler(router, auth.HandlerDeps{
		Service: authService,
	})

	server := http.Server{
		Addr:    ":8080",
		Handler: middleware.Logging(router),
	}

	logrus.Info("Starting server on :8080")
	server.ListenAndServe()
}
