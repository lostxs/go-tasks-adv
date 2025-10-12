package session

import (
	"errors"
	"sync"
)

type Repository struct {
	sessions map[string]*Session
	mu       sync.Mutex
}

func NewRepository() *Repository {
	return &Repository{
		sessions: make(map[string]*Session),
	}
}

func (r *Repository) Create(session *Session) (*Session, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.sessions[session.ID] = session
	return session, nil
}

func (r *Repository) FindByID(id string) (*Session, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	session, ok := r.sessions[id]
	if !ok {
		return nil, errors.New("session not found")
	}

	return session, nil
}
