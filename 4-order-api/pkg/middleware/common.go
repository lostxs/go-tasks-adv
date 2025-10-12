package middleware

import "net/http"

// Custom response writer to capture status code and response size
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

// Override WriteHeader to capture status code
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Override Write to capture response size
func (rw *responseWriter) Write(b []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(b)
	rw.size += size
	return size, err
}
