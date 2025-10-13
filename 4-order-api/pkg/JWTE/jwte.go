package jwte

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTData struct {
	Number    string
	SessionId string
}
type JWTE struct {
	Secret string
}

func NewJWT(secret string) *JWTE {
	return &JWTE{
		Secret: secret,
	}
}
func (j *JWTE) Create(data JWTData) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sessionId": data.SessionId,
		"number":    data.Number,
		"exp":       time.Now().Add(72 * time.Hour).Unix(),
	})
	fmt.Println(data.Number, data.SessionId)
	str, err := token.SignedString([]byte(j.Secret))
	if err != nil {
		return "", err
	}
	return str, nil
}
func (j *JWTE) Parse(token string) (bool, *JWTData) {
	tok, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.Secret), nil
	})
	if err != nil {
		return false, nil
	}
	number, ok := tok.Claims.(jwt.MapClaims)["number"]
	if !ok {
		return false, nil
	}
	return tok.Valid, &JWTData{
		Number: number.(string),
	}
}
