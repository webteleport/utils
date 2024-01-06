package utils

import (
	"bufio"
	"io"
	"log"
	"net"
	"net/http"
	"os"
)

// InterceptMiddleware prints request & response info to stdout
func InterceptMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("RequestURI:", r.RequestURI)
		log.Println("RequestHeader", r.Header)
		iw := &Interceptor{ResponseWriter: w, StatusCode: http.StatusOK}
		r.Body = &Body{r.Body, os.Stderr}
		next.ServeHTTP(iw, r)
		log.Println("StatusCode:", iw.StatusCode)
		log.Println("ResponseHeader:", iw.Header())
	})
}

func (w *Interceptor) Flush() {
	if flusher, ok := w.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

func (w *Interceptor) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	conn, rw, err := w.ResponseWriter.(http.Hijacker).Hijack()
	if err == nil && w.StatusCode == 0 {
		// The status will be StatusSwitchingProtocols if there was no error and
		// WriteHeader has not been called yet
		w.StatusCode = http.StatusSwitchingProtocols
	}
	return conn, rw, err
}

type Interceptor struct {
	http.ResponseWriter
	StatusCode int
}

func (w *Interceptor) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.StatusCode = statusCode
}

type Body struct {
	io.ReadCloser
	Writer io.Writer
}

func (b *Body) Read(p []byte) (int, error) {
	t := io.TeeReader(b.ReadCloser, b.Writer)
	return t.Read(p)
}
