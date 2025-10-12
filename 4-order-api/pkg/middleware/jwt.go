package middleware

import (
	"4-order-api/pkg/jwt"
	"context"
	"net/http"
	"strings"
)

type userKey struct{}

type JWTAuthMiddleware struct {
	jwt *jwt.JWT
}

func NewJWTAuthMiddleware(jwt *jwt.JWT) *JWTAuthMiddleware {
	return &JWTAuthMiddleware{jwt: jwt}
}

func (m *JWTAuthMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		claims, err := m.jwt.ParseToken(tokenStr)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), userKey{}, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *JWTAuthMiddleware) GetUser(ctx context.Context) *jwt.Payload {
	claims, ok := ctx.Value(userKey{}).(*jwt.Payload)
	if !ok {
		return nil
	}
	return claims
}
