package auth

import (
	"4-order-api/internal/models"
	"4-order-api/internal/user"
	"4-order-api/pkg/rand"
	"errors"
)

type AuthService struct {
	UserRepository *user.UserRepository
}

func NewAuthRepository(userRepository *user.UserRepository) *AuthService {
	return &AuthService{
		UserRepository: userRepository,
	}
}

func (service *AuthService) Register(number string) (*models.User, error) {
	// идем в БД проверям на наличие номер

	//если нет то генерим sesionID и выдаем иначе генерим новый

	user := &models.User{
		Number:    number,
		SessionID: rand.RandSession(),
	}
	createdUser, err := service.UserRepository.CreateUser(user)
	if err != nil {
		return nil, errors.New("ошибка создания пользователя")
	}
	return createdUser, nil
}
func (service *AuthService) Update(user *models.User) (*models.User, error) {

	user.SessionID = rand.RandSession()

	updateUser, err := service.UserRepository.UpdateSessionId(user)
	if err != nil {
		return nil, errors.New("ошибка обновления пользователя")
	}
	return updateUser, nil
}
func (service *AuthService) UpdateCode(user *models.User, code string) (*models.User, error) {
	updateUser, err := service.UserRepository.UpdateCode(user, code)
	if err != nil {
		return nil, errors.New("ошибка обновления пользователя")
	}
	return updateUser, nil
}
