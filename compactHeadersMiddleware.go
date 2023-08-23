package utils

import (
	"net/http"

	"golang.org/x/exp/slices"
)

func CompactHeadersMiddleware(next http.Handler, keys []string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Wrap the response writer to filter headers
		wrappedWriter := &headerFilterResponseWriter{
			ResponseWriter: w,
			HeaderKeys:     keys,
		}

		next.ServeHTTP(wrappedWriter, r)
	})
}

// Custom ResponseWriter to filter headers
type headerFilterResponseWriter struct {
	http.ResponseWriter
	HeaderKeys []string
}

func (hw *headerFilterResponseWriter) WriteHeader(code int) {
	// Filter out duplicate headers based on values before writing the response header
	headers := hw.ResponseWriter.Header()
	for key, values := range headers {
		if slices.Contains(hw.HeaderKeys, key) {
			compact := slices.Compact(values)
			headers[key] = compact
		}
	}

	hw.ResponseWriter.WriteHeader(code)
}
