package user

import "4-order-api/pkg/db"

type Repository struct {
	db *db.DB
}

func NewRepository(database *db.DB) *Repository {
	return &Repository{db: database}
}

func (r *Repository) FindByPhone(phone string) (*User, error) {
	var u User
	result := r.db.Where("phone = ?", phone).First(&u)
	if result.Error != nil {
		return nil, result.Error
	}
	return &u, nil
}

func (r *Repository) Create(user *User) (*User, error) {
	result := r.db.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}
