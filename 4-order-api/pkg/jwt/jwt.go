package jwt

import (
	"maps"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWT struct {
	secret string
}

func NewJWT(secret string) *JWT {
	return &JWT{
		secret: secret,
	}
}

func (j *JWT) GenerateToken(claims map[string]any) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claimsMap := token.Claims.(jwt.MapClaims)

	maps.Copy(claimsMap, claims)
	claimsMap["exp"] = time.Now().Add(24 * time.Hour).Unix()

	return token.SignedString([]byte(j.secret))
}
