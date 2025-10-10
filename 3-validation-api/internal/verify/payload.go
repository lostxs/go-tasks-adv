package verify

type SendRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type VerifyRequest struct {
	Hash string `json:"hash" validate:"required"`
}
