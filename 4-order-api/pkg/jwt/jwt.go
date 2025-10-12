package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Payload struct {
	UserID uint   `json:"user_id"`
	Phone  string `json:"phone"`
	jwt.RegisteredClaims
}

type JWT struct {
	secret string
}

func NewJWT(secret string) *JWT {
	j := &JWT{secret: secret}
	return j
}

func (j *JWT) GenerateToken(payload Payload) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": payload.UserID,
		"phone":   payload.Phone,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})
	return token.SignedString([]byte(j.secret))
}

func (j *JWT) ParseToken(tokenStr string) (*Payload, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &Payload{
			UserID: uint(claims["user_id"].(float64)),
			Phone:  claims["phone"].(string),
		}, nil
	}
	return nil, err
}
