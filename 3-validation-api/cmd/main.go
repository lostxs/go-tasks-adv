package main

import (
	"3-validation-api/config"
	"3-validation-api/internal/user"
	"3-validation-api/internal/verify"
	"3-validation-api/pkg/file"
	"fmt"
	"log"
	"net/http"
)

func main() {
	cfg := config.Load()
	userRepo := user.NewRepository(file.NewJsonDb("users.json"))
	router := http.NewServeMux()
	verify.NewVerifyHandler(router, verify.VerifyHanlderDeps{
		Config:   &cfg.Mail,
		UserRepo: userRepo,
	})
	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	fmt.Println("Сервер запущен на :8080")
	log.Fatal(server.ListenAndServe())
}
