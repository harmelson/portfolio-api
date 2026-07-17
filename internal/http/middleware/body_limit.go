package middleware

import (
	"net/http"
)

const (
	MaxBodySize int64 = 1 << 20 // 1MB
)

func BodyLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, MaxBodySize)
		next.ServeHTTP(w, r)
	})
}
