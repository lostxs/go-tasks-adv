package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB   DBConfig
	Auth AuthConfig
}

type DBConfig struct {
	Dsn string
}

type AuthConfig struct {
	Secret string
}

func Load() *Config {
	if err := godotenv.Load(".env"); err != nil {
		panic(err)
	}

	return &Config{
		DB: DBConfig{
			os.Getenv("DB_DSN"),
		},
		Auth: AuthConfig{
			Secret: os.Getenv("JWT_SECRET"),
		},
	}
}
