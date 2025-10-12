package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB DBConfig
}

type DBConfig struct {
	Dsn string
}

func Load() *Config {
	if err := godotenv.Load(".env"); err != nil {
		panic(err)
	}

	return &Config{
		DB: DBConfig{
			os.Getenv("DB_DSN"),
		},
	}
}
