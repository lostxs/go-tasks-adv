package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Db   DbConfig
	Auth AuthConfig
}
type DbConfig struct {
	Dsn string
}
type AuthConfig struct {
	Secret string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Ошибка загрузки .env файла")
	}
	secret, ok := os.LookupEnv("SECRET")
	if !ok {
		log.Fatalln("Отсутствует секрет в переменной окружения")
	}
	return &Config{
		Db: DbConfig{
			Dsn: os.Getenv("DSN"),
		},
		Auth: AuthConfig{
			Secret: secret,
		},
	}
}
