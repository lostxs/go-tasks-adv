package auth

import (
	"4-order-api/configs"
	jwte "4-order-api/pkg/JWTE"
	"4-order-api/pkg/req"
	"4-order-api/pkg/resp"
	"net/http"
)

type AuthHandlerDeps struct {
	*configs.Config
	*AuthService
}
type AuthHandler struct {
	*configs.Config
	*AuthService
}

func NewAuthHandler(router *http.ServeMux, deps AuthHandlerDeps) {
	handler := &AuthHandler{
		Config:      deps.Config,
		AuthService: deps.AuthService,
	}
	router.HandleFunc("POST /auth", handler.register)
	router.HandleFunc("POST /auth/verify", handler.verify)
}
func (handler *AuthHandler) register(w http.ResponseWriter, request *http.Request) {

	// получаем тело запроса
	body, err := req.HandleBody[RegisterRequest](w, request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	findUser, _ := handler.AuthService.UserRepository.FindUserByNum(body.Number)
	if findUser == nil {
		user, _ := handler.AuthService.Register(body.Number)
		resp.Json(w, user, 201)

	} else {

		userWithNewSession, err := handler.AuthService.Update(findUser)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		resp.Json(w, userWithNewSession, 200)
	}

	//запрашиваем код подтверждения
	// проверяем sesiodID и код
	//выдаем токен
}
func (handler *AuthHandler) verify(w http.ResponseWriter, request *http.Request) {
	body, err := req.HandleBody[VerifyRequest](w, request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	findUser, err := handler.AuthService.UserRepository.FindUserBySession(body.SessionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if body.SessionID == "" || body.Code == "" {
		http.Error(w, "Сессия или код истекли", http.StatusBadRequest)
		return
	}
	if findUser.SessionID == body.SessionID && findUser.Code == body.Code {

		//выдаем токен
		secret, err := jwte.NewJWT(handler.Auth.Secret).Create(jwte.JWTData{
			Number:    findUser.Number,
			SessionId: body.SessionID,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		data := VerifyResponse{
			Token: secret,
		}
		resp.Json(w, data, 200)

	} else {
		_, err := handler.AuthService.UpdateCode(findUser, body.Code)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		resp.Json(w, "Неверный код, Код отправлен", 200)
	}
}
