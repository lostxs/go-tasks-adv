package auth

import (
	"4-order-api/pkg/jwt"
	"4-order-api/pkg/middleware"
	"4-order-api/pkg/request"
	"4-order-api/pkg/response"
	"errors"
	"net/http"
)

type HandlerDeps struct {
	Service *Service
	JWT     *jwt.JWT
}

type Handler struct {
	service *Service
}

func NewHandler(router *http.ServeMux, deps HandlerDeps) {
	h := &Handler{
		service: deps.Service,
	}

	router.Handle("GET /auth/profile", middleware.JWTAuth(deps.JWT, h.Profile()))

	router.HandleFunc("POST /auth/send", h.SendCode())
	router.HandleFunc("POST /auth/verify", h.VerifyCode())
}

func (h *Handler) Profile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middleware.GetUser(r.Context())

		response.WriteJSON(w, http.StatusOK, map[string]any{
			"phone": claims.Phone,
			"exp":   claims.ExpiresAt,
		})
	}
}

func (h *Handler) SendCode() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := request.ParseBody[SendCodeRequest](w, r)
		if err != nil {
			return
		}

		sessionID, err := h.service.SendCode(body.Phone)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response.WriteJSON(w, http.StatusCreated, sessionID)
	}
}

func (h *Handler) VerifyCode() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := request.ParseBody[VerifyCodeRequest](w, r)
		if err != nil {
			return
		}

		token, err := h.service.VerifyCode(body.SessionID, body.Code)
		if err != nil {
			if errors.Is(err, ErrSessionNotFound) || errors.Is(err, ErrInvalidCode) {
				response.WriteJSON(w, http.StatusUnauthorized, err.Error())
			} else {
				response.WriteJSON(w, http.StatusInternalServerError, err.Error())
			}
			return
		}

		response.WriteJSON(w, http.StatusCreated, token)
	}
}
