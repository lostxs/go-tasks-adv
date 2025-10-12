package auth

import (
	"4-order-api/pkg/request"
	"4-order-api/pkg/response"
	"errors"
	"net/http"
)

type Handler struct {
	service *Service
}

type HandlerDeps struct {
	Service *Service
}

func NewHandler(router *http.ServeMux, deps HandlerDeps) {
	h := &Handler{
		service: deps.Service,
	}

	router.HandleFunc("POST /auth/send", h.SendCode)
	router.HandleFunc("POST /auth/verify", h.VerifyCode)
}

func (h *Handler) SendCode(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler) VerifyCode(w http.ResponseWriter, r *http.Request) {
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
