package session

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"
)

type Session struct {
	ID        string
	Phone     string
	Code      string
	ExpiresAt time.Time
}

func NewSession(phone string) *Session {
	session := &Session{
		ID:        generateID(),
		Phone:     phone,
		Code:      generateCode(),
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}
	return session
}

func generateID() string {
	idBytes := make([]byte, 16)
	rand.Read(idBytes)
	return hex.EncodeToString(idBytes)
}

func generateCode() string {
	codeNum, _ := rand.Int(rand.Reader, big.NewInt(10000))
	return fmt.Sprintf("%04d", codeNum)
}
