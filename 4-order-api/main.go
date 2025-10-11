package main

import (
	"4-order-api/config"
	"4-order-api/internal/product"
	"4-order-api/pkg/db"
	"fmt"
	"log"
	"net/http"
)

func main() {
	config := config.Load()
	db := db.NewDB(config.DB)
	router := http.NewServeMux()

	proudctRepo := product.NewRepository(db)

	product.NewHandler(router, product.HandlerDeps{
		Repo: proudctRepo,
	})

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	fmt.Println("Server listening on port 8080")
	log.Fatal(server.ListenAndServe())
}
