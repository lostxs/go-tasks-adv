package verify

import (
	"3-validation-api/config"
	"3-validation-api/internal/user"
	"3-validation-api/pkg/request"
	"3-validation-api/pkg/response"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/smtp"
	"strings"

	"github.com/jordan-wright/email"
)

type VerifyHanlderDeps struct {
	Config   *config.MailConfig
	UserRepo *user.RepositoryWithDb
}

type VerifyHandler struct {
	config   *config.MailConfig
	userRepo *user.RepositoryWithDb
}

func NewVerifyHandler(router *http.ServeMux, deps VerifyHanlderDeps) {
	handler := &VerifyHandler{
		config:   deps.Config,
		userRepo: deps.UserRepo,
	}

	router.HandleFunc("POST /send", handler.Send())
	router.HandleFunc("GET /verify/", handler.Verify())
}

func (handler *VerifyHandler) Send() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := request.HandleBody[SendRequest](w, r)
		if err != nil {
			return
		}

		existing := handler.userRepo.FindOne(body.Email, func(u user.User, s string) bool {
			return u.Email == s
		})
		if existing != nil {
			if existing.Valid {
				response.Json(w, "Email уже подтверждён", http.StatusConflict)
				return
			}
			response.Json(w, "Email уже ожидает подтверждения", http.StatusConflict)
			return
		}

		hash, err := generateHash(body.Email)
		if err != nil {
			response.Json(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = handler.sendVerificationEmail(body.Email, hash)
		if err != nil {
			response.Json(w, err.Error(), http.StatusInternalServerError)
			return
		}

		newUser := user.NewUser(body.Email, hash, false)
		err = handler.userRepo.AddOne(*newUser)
		if err != nil {
			response.Json(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response.Json(w, true, http.StatusOK)
	}
}

func (handler *VerifyHandler) Verify() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hash := strings.TrimPrefix(r.URL.Path, "/verify/")
		if hash == "" {
			response.Json(w, "Хэш не указан", http.StatusBadRequest)
			return
		}

		existing := handler.userRepo.FindOne(hash, func(u user.User, s string) bool {
			return u.Hash == s
		})
		if existing == nil {
			response.Json(w, false, http.StatusNotFound)
			return
		}

		err := handler.userRepo.DeleteOne(hash, func(u user.User, s string) bool {
			return u.Hash == s
		})
		if err != nil {
			response.Json(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response.Json(w, true, http.StatusOK)
	}
}

func (handler *VerifyHandler) sendVerificationEmail(toEmail, hash string) error {
	e := &email.Email{
		From:    handler.config.Email,
		To:      []string{toEmail},
		Subject: "Подтверждение Email",
		Text:    fmt.Appendf(nil, "Перейдите по ссылке для подтверждения: http://localhost:8080/verify/%s", hash),
		HTML: fmt.Appendf(nil, `
		<h1>Подтверждение Email</h1>
		<p>Перейдите по <a href="http://localhost:8080/verify/%s">ссылке</a> для подтверждения</p>`, hash),
	}

	auth := smtp.PlainAuth("", handler.config.Email, handler.config.Password, handler.config.Host)
	return e.Send(handler.config.Host+":587", auth)
}

func generateHash(email string) (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	h := sha256.Sum256(append([]byte(email), b...))
	return hex.EncodeToString(h[:]), nil
}
