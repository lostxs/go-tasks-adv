package main

import (
	"4-order-api/internal/models"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	db, err := gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.SetupJoinTable(&models.Order{}, "Products", &models.OrderProduct{})
	if err != nil {
		log.Fatal(err)
	}
	err = db.SetupJoinTable(&models.Product{}, "Orders", &models.OrderProduct{})
	if err != nil {
		log.Fatal(err)
	}
	db.Migrator().DropTable(&models.OrderProduct{}, &models.Product{}, &models.User{}, &models.Order{})
	db.AutoMigrate(&models.Product{}, &models.User{}, &models.Order{}, &models.OrderProduct{})
	log.Println("Новые поля добавлены.")
}
