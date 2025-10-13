package di

import "4-order-api/internal/models"

//можно добавить order interface чтобы package order не зависил от package models

type IUserRepository interface {
	FindUserByNum(number string) (*models.User, error)
}

type IProductRepository interface {
	FindProductArrayById(productIds []uint) ([]*models.Product, error)
}
