package user

import (
	"encoding/json"
	"sync"
	"time"
)

type Db interface {
	Read() ([]byte, error)
	Write([]byte) error
}

type Repository struct {
	Users     []User    `json:"users"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type RepositoryWithDb struct {
	Repository
	db Db
	mu sync.Mutex
}

func NewRepository(db Db) *RepositoryWithDb {
	data, err := db.Read()
	if err != nil {
		return &RepositoryWithDb{
			Repository: Repository{
				Users:     []User{},
				UpdatedAt: time.Now(),
			},
			db: db,
			mu: sync.Mutex{},
		}
	}
	var repository Repository
	err = json.Unmarshal(data, &repository)
	if err != nil {
		return &RepositoryWithDb{
			Repository: Repository{
				Users:     []User{},
				UpdatedAt: time.Now(),
			},
			db: db,
			mu: sync.Mutex{},
		}
	}
	return &RepositoryWithDb{
		Repository: repository,
		db:         db,
		mu:         sync.Mutex{},
	}
}

func (r *RepositoryWithDb) AddOne(user User) error {
	r.Users = append(r.Users, user)
	return r.save()
}

func (r *RepositoryWithDb) FindMany(str string, checker func(User, string) bool) []User {
	var users []User
	for _, user := range r.Users {
		isMatch := checker(user, str)
		if isMatch {
			users = append(users, user)
		}
	}
	return users
}

func (r *RepositoryWithDb) FindOne(str string, checker func(User, string) bool) *User {
	for i := range r.Users {
		u := &r.Users[i]
		if checker(*u, str) {
			return u
		}
	}
	return nil
}

func (r *RepositoryWithDb) DeleteOne(str string, checker func(User, string) bool) error {
	for i, u := range r.Users {
		isMatch := checker(u, str)
		if isMatch {
			r.Users = append(r.Users[:i], r.Users[i+1:]...)
			return r.save()
		}
	}
	return nil
}

func (r *RepositoryWithDb) save() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.UpdatedAt = time.Now()
	data, err := json.MarshalIndent(r.Repository, "", "  ")
	if err != nil {
		return err
	}
	return r.db.Write(data)
}
