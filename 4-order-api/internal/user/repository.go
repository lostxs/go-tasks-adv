package user

import (
	"4-order-api/internal/models"
	"4-order-api/pkg/db"
	"fmt"
)

type UserRepository struct {
	Database *db.Db
}

func NewUserRepository(database *db.Db) *UserRepository {
	return &UserRepository{
		Database: database,
	}
}
func (repo *UserRepository) CreateUser(user *models.User) (*models.User, error) {
	result := repo.Database.DB.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}
func (repo *UserRepository) FindUserByNum(number string) (*models.User, error) {
	var user models.User
	result := repo.Database.DB.First(&user, "number= ? ", number)
	if result.Error != nil {
		return nil, result.Error
	}
	fmt.Printf("Имя пользователя: %v ID Пользователя: %v\n", user.Number, user.ID)
	return &user, nil
}
func (repo *UserRepository) UpdateSessionId(user *models.User) (*models.User, error) {
	//fmt.Printf("Имя пользователя: %v Пароль: %v\n", user.Number, user.SessionID)
	result := repo.Database.DB.Model(&models.User{}).Where("id = ?", user.ID).Update("session_id", user.SessionID)
	if result.Error != nil {
		return nil, result.Error
	}
	//fmt.Printf("Имя пользователя: %v Пароль: %v\n", user.Number, user.SessionID)
	return user, nil

}
func (repo *UserRepository) FindUserBySession(session string) (*models.User, error) {
	var user models.User
	result := repo.Database.DB.First(&user, "session_id = ? ", session)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
func (repo *UserRepository) UpdateCode(user *models.User, code string) (*models.User, error) {

	result := repo.Database.DB.Model(&models.User{}).Where("id = ?", user.ID).Update("code", code)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil

}
