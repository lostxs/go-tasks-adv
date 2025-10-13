package middleware

import (
	"4-order-api/configs"
	jwte "4-order-api/pkg/JWTE"
	"context"
	"net/http"
	"strings"
)

type key string

const (
	ContextPhoneNumber key = "ContextPhoneNumber"
)

func Auth(next http.HandlerFunc, config *configs.Config) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer") {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
			return
		}
		//парсим токен
		token := strings.TrimPrefix(auth, "Bearer ")
		isValid, data := jwte.NewJWT(config.Auth.Secret).Parse(token)
		if !isValid {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
			return
		}
		ctx := context.WithValue(r.Context(), ContextPhoneNumber, data.Number)
		req := r.WithContext(ctx)
		next.ServeHTTP(w, req)
	})

}
