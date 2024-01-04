package utils

import (
	"io"
	"log"
	"net/http"
	"os"
)

// InterceptMiddleware prints request & response info to stdout
func InterceptMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("RequestURI:", r.RequestURI)
		log.Println("RequestHeader", r.Header)
		iw := &Interceptor{ResponseWriter: w}
		r.Body = &Body{r.Body, os.Stderr}
		next.ServeHTTP(iw, r)
		log.Println("StatusCode:", iw.StatusCode)
		log.Println("ResponseHeader:", iw.Header())
	})
}

type Interceptor struct {
	http.ResponseWriter
	StatusCode int
}

func (w *Interceptor) WriteHeader(statusCode int) {
	w.StatusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

type Body struct {
	io.ReadCloser
	Writer io.Writer
}

func (b *Body) Read(p []byte) (int, error) {
	t := io.TeeReader(b.ReadCloser, b.Writer)
	return t.Read(p)
}
