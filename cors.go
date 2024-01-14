package utils

import (
	"net/http"

	"github.com/rs/cors"
)

var AllowAllCorsMiddleware = func(next http.Handler) http.Handler {
	return cors.AllowAll().Handler(next)
}
