package user

type User struct {
	Email string `json:"email"`
	Hash  string `json:"hash"`
	Valid bool   `json:"valid"`
}

func NewUser(email, hash string, valid bool) *User {
	return &User{
		Email: email,
		Hash:  hash,
		Valid: valid,
	}
}
