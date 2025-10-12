package db

import (
	"4-order-api/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	*gorm.DB
}

func NewDB(config config.DBConfig) *DB {
	db, err := gorm.Open(postgres.Open(config.Dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return &DB{db}
}
