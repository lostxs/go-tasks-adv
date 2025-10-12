package auth

type SendCodeRequest struct {
	Phone string `json:"phone" validate:"required"`
}

type VerifyCodeRequest struct {
	SessionID string `json:"sessionId" validate:"required"`
	Code      string `json:"code" validate:"required"`
}
