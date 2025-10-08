package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"sync"

	"github.com/jordan-wright/email"
)

type Config struct {
	Email    string
	Password string
	Address  string
}

type User struct {
	Email string `json:"email"`
	Hash  string `json:"hash"`
	Valid bool   `json:"valid"`
}

var (
	config Config
	users  []User
	mu     sync.Mutex
)

func loadConfig() {
	config = Config{
		Email:    os.Getenv("SMTP_EMAIL"),
		Password: os.Getenv("SMTP_PASSWORD"),
		Address:  os.Getenv("SMTP_ADDRESS"),
	}
	if config.Email == "" || config.Password == "" || config.Address == "" {
		log.Fatal("Не заданы переменные окружения: SMTP_EMAIL, SMTP_PASSWORD, SMTP_ADDRESS")
	}
}

func loadUsers() {
	data, err := os.ReadFile("users.json")
	if err == nil {
		_ = json.Unmarshal(data, &users)
	}
}

func saveUsers() {
	data, _ := json.MarshalIndent(users, "", "  ")
	_ = os.WriteFile("users.json", data, 0644)
}

func generateHash(email string) string {
	h := sha256.Sum256([]byte(email))
	return hex.EncodeToString(h[:])
}

func userExists(emailAddr string) *User {
	for i := range users {
		if users[i].Email == emailAddr {
			return &users[i]
		}
	}
	return nil
}

func sendHandler(w http.ResponseWriter, r *http.Request) {
	emailAddr := r.FormValue("email")
	if emailAddr == "" {
		http.Error(w, "email required", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	if existing := userExists(emailAddr); existing != nil {
		if existing.Valid {
			http.Error(w, "Email уже подтвержден ✅", http.StatusConflict)
			return
		}
		http.Error(w, "Email уже ожидает подтверждения", http.StatusConflict)
		return
	}

	hash := generateHash(emailAddr)

	e := email.NewEmail()
	e.From = config.Email
	e.To = []string{emailAddr}
	e.Subject = "Подтверждение Email"
	verifyLink := fmt.Sprintf("http://localhost:8080/verify/%s", hash)
	e.Text = []byte(fmt.Sprintf("Перейдите по ссылке для подтверждения: %s", verifyLink))

	auth := smtp.PlainAuth("", config.Email, config.Password, config.Address)
	if err := e.Send(config.Address+":587", auth); err != nil {
		http.Error(w, "Ошибка отправки письма: "+err.Error(), http.StatusInternalServerError)
		return
	}

	users = append(users, User{Email: emailAddr, Hash: hash, Valid: false})
	saveUsers()

	w.Write([]byte("Письмо отправлено ✅"))
}

func verifyHandler(w http.ResponseWriter, r *http.Request) {
	hash := r.URL.Path[len("/verify/"):]
	if hash == "" {
		http.Error(w, "hash required", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()
	for i := range users {
		if users[i].Hash == hash {
			users[i].Valid = true
			saveUsers()
			w.Write([]byte("Email подтверждён ✅"))
			return
		}
	}

	http.Error(w, "Хэш не найден", http.StatusNotFound)
}

func main() {
	loadConfig()
	loadUsers()

	http.HandleFunc("/send", sendHandler)
	http.HandleFunc("/verify/", verifyHandler)

	fmt.Println("Сервер запущен на :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
