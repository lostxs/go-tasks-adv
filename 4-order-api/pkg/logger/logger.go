package logger

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func LogInit() {
	log.SetFormatter(&log.JSONFormatter{})
	file, err := os.OpenFile("logrus.json", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Не удалось открыть файл логов %v", err)
	}
	log.SetOutput(file)
	log.SetLevel(log.InfoLevel)
}
