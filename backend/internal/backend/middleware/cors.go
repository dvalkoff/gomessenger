package middleware

import (
	"net/http"
)

type CorsMiddleware interface {
	HandleCors(next http.Handler) http.Handler
}

type corsMiddleware struct {
	frontendOrigin string
}

func NewCorsMiddleware(frontendOrigin string) CorsMiddleware {
	return &corsMiddleware{frontendOrigin: frontendOrigin}
}

func (cm *corsMiddleware) HandleCors(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			headers := w.Header()
			headers.Set("Access-Control-Allow-Origin", cm.frontendOrigin)
			if r.Method == http.MethodOptions {
				headers.Set("Access-Control-Allow-Methods", "*")
				headers.Set("Access-Control-Allow-Headers", "*")
				headers.Set("Access-Control-Allow-Credentials", "true")
				headers.Set("Access-Control-Max-Age", "86400")
				w.WriteHeader(200)
				return
			}
			next.ServeHTTP(w, r)
		},
	)
}
