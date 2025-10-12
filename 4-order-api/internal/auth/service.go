package auth

import (
	"4-order-api/pkg/jwt"
	"fmt"
	"time"
)

type Service struct {
	repo *Repository
	jwt  *jwt.JWT
}

func NewService(repo *Repository, jwt *jwt.JWT) *Service {
	return &Service{repo: repo, jwt: jwt}
}

func (s *Service) SendCode(phone string) (string, error) {
	session := NewSession(phone)
	_, err := s.repo.Create(session)
	if err != nil {
		return "", err
	}

	// SMS-service mock call
	fmt.Printf("SMS sent to %s: code=%s\n", phone, session.Code)

	return session.ID, nil
}

func (s *Service) VerifyCode(sessionID, code string) (string, error) {
	existingSession, _ := s.repo.FindByID(sessionID)
	if existingSession == nil {
		return "", ErrSessionNotFound
	}

	if time.Now().After(existingSession.ExpiresAt) {
		return "", ErrSessionExpired
	}

	if existingSession.Code != code {
		return "", ErrInvalidCode
	}

	token, err := s.jwt.GenerateToken(existingSession.Phone)
	if err != nil {
		return "", err
	}

	return token, nil
}
