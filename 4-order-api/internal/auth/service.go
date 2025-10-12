package auth

import (
	"4-order-api/internal/session"
	"4-order-api/internal/user"
	"4-order-api/pkg/jwt"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type ServiceDeps struct {
	SessionRepository *session.Repository
	UserRepository    *user.Repository
	Jwt               *jwt.JWT
}

type Service struct {
	sessionRepository *session.Repository
	userRepository    *user.Repository
	jwt               *jwt.JWT
}

func NewService(deps ServiceDeps) *Service {
	return &Service{
		sessionRepository: deps.SessionRepository,
		userRepository:    deps.UserRepository,
		jwt:               deps.Jwt,
	}
}

func (s *Service) SendCode(phone string) (string, error) {
	session := session.NewSession(phone)
	_, err := s.sessionRepository.Create(session)
	if err != nil {
		return "", err
	}

	fmt.Printf("SMS sent to %s: code=%s\n", phone, session.Code)
	return session.ID, nil
}

func (s *Service) VerifyCode(sessionID, code string) (string, error) {
	session, _ := s.sessionRepository.FindByID(sessionID)
	if session == nil {
		return "", ErrSessionNotFound
	}

	if time.Now().After(session.ExpiresAt) {
		return "", ErrSessionExpired
	}

	if session.Code != code {
		return "", ErrInvalidCode
	}

	u, err := s.userRepository.FindByPhone(session.Phone)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			u, err = s.userRepository.Create(&user.User{Phone: session.Phone})
			if err != nil {
				return "", err
			}
		} else {
			return "", err
		}
	}

	token, err := s.jwt.GenerateToken(jwt.Payload{
		UserID: u.ID,
		Phone:  u.Phone,
	})
	if err != nil {
		return "", err
	}

	return token, nil
}
