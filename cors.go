package utils

import (
	"net/http"

	"github.com/rs/cors"
)

var AllowAllCorsMiddleware = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "" {
			origin = "*"
		}
		config := cors.New(cors.Options{
			AllowedOrigins: []string{
				origin,
			},
			AllowedMethods: []string{
				http.MethodHead,
				http.MethodGet,
				http.MethodPost,
				http.MethodPut,
				http.MethodPatch,
				http.MethodDelete,
			},
			AllowedHeaders:   []string{"*"},
			AllowCredentials: true,
		})

		config.Handler(next).ServeHTTP(w, r)
	})
}
