package middleware

import (
	"context"
	"net/http"
	"strings"

	"4-order-api/pkg/jwt"
)

type contextKey string

const userKey = contextKey("user")

func jwtAuthFailed(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
}

func JWTAuth(jwt *jwt.JWT, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authedHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authedHeader, "Bearer ") {
			jwtAuthFailed(w)
			return
		}
		token := strings.TrimPrefix(authedHeader, "Bearer ")

		claims, err := jwt.VerifyToken(token)
		if err != nil {
			jwtAuthFailed(w)
			return
		}

		ctx := context.WithValue(r.Context(), userKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUser(ctx context.Context) *jwt.Claims {
	if ctx == nil {
		return nil
	}
	if claims, ok := ctx.Value(userKey).(*jwt.Claims); ok {
		return claims
	}
	return nil
}
