package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Mail MailConfig
}

type MailConfig struct {
	Email    string
	Password string
	Host     string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла:", err)
	}

	config := &Config{
		Mail: MailConfig{
			Email:    os.Getenv("SMTP_EMAIL"),
			Password: os.Getenv("SMTP_PASSWORD"),
			Host:     os.Getenv("SMTP_HOST"),
		},
	}

	if config.Mail.Email == "" || config.Mail.Password == "" || config.Mail.Host == "" {
		log.Fatal("Не заданы переменные окружения: SMTP_EMAIL, SMTP_PASSWORD, SMTP_ADDRESS")
	}

	return config
}
