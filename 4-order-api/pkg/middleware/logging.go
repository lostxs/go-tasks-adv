package middleware

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapper := &WrapperWriter{
			ResponseWriter: w,
			StatusCode:     http.StatusOK,
		}
		next.ServeHTTP(wrapper, r)
		//log.Println(wrapper.StatusCode, r.Method, r.URL.Path, time.Since(start))
		log.WithFields(log.Fields{

			"statuscode": wrapper.StatusCode,
			"method":     r.Method,
			"urlpath":    r.URL.Path,
			"timework":   time.Since(start),
		}).Info("Получен ")
	})
	//обновление
}
