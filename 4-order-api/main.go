package main

import (
	"4-order-api/config"
	"4-order-api/pkg/db"
	"fmt"
)

func main() {
	config := config.Load()
	_ = db.New(config.DB)
	fmt.Println("OK")
}
