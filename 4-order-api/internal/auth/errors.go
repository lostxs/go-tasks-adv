package auth

import "errors"

var (
	ErrInvalidCode     = errors.New("invalid code")
	ErrSessionExpired  = errors.New("session expired")
	ErrSessionNotFound = errors.New("session not found")
)
