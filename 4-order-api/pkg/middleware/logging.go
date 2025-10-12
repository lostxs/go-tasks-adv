package middleware

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a logger with request details
		log := logrus.WithFields(logrus.Fields{
			"method":      r.Method,
			"path":        r.URL.Path,
			"query":       r.URL.RawQuery,
			"remote_addr": r.RemoteAddr,
			"user_agent":  r.UserAgent(),
			"referer":     r.Referer(),
		})

		// Create a custom response writer to capture the status code
		rw := &responseWriter{w, http.StatusOK, 0}

		start := time.Now()
		log.Info("HTTP request started")

		next.ServeHTTP(rw, r)

		duration := time.Since(start)

		// Log the completed request with additional information
		log.WithFields(logrus.Fields{
			"status_code": rw.statusCode,
			"duration_ms": duration.Milliseconds(),
			"size_bytes":  rw.size,
		}).Info("HTTP request completed")
	})
}
