package middleware

import "net/http"

func RegisterLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		println(r) // before

		defer println(w) // after

		next.ServeHTTP(w, r)
	})
}
