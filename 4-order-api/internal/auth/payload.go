package auth

type RegisterRequest struct {
	Number string `json:"number"`
}
type VerifyRequest struct {
	SessionID string `json:"sessionid"`
	Code      string `json:"code"`
}
type VerifyResponse struct {
	Token string `json:"token"`
}
