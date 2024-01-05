package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/mattn/go-isatty"
)

type GinLogger struct {
	next      http.Handler
	formatter LogFormatter
	logWriter io.Writer
	skipPaths []string
}

type LogFormatter func(params LogFormatterParams) string

type LogFormatterParams struct {
	Request      *http.Request
	TimeStamp    time.Time
	StatusCode   int
	Latency      time.Duration
	ClientIP     string
	Host         string
	Method       string
	Proto        string
	Path         string
	ErrorMessage string
	BodySize     int
	Keys         map[string]interface{}
}

var defaultLogFormatter = func(param LogFormatterParams) string {
	var statusColor, methodColor, resetColor string
	if param.IsOutputColor() {
		statusColor = param.StatusCodeColor()
		methodColor = param.MethodColor()
		resetColor = param.ResetColor()
	}

	if param.Latency > time.Minute {
		param.Latency = param.Latency.Truncate(time.Second)
	}
	return fmt.Sprintf("%v |%s %3d %s| %8v |%s %-7s %s| %8s | %15s | %20s | %#v\n",
		param.TimeStamp.Format("2006/01/02 15:04:05"),
		statusColor, param.StatusCode, resetColor,
		formatDuration(param.Latency),
		methodColor, param.Method, resetColor,
		param.Proto,
		param.ClientIP,
		param.Host,
		param.Path,
	)
}

func formatDuration(duration time.Duration) string {
	// Convert the duration to a string
	durationStr := duration.String()

	// Find the index of the first non-digit character
	index := strings.IndexFunc(durationStr, func(r rune) bool {
		return unicode.IsLetter(r)
	})

	// Extract the numeric value and unit
	valueStr := durationStr[:index]
	unitStr := durationStr[index:]

	// Parse the value part as a float64
	value, _ := strconv.ParseFloat(valueStr, 64)

	// Format the value with 2 decimal places
	formattedValue := fmt.Sprintf("%.2f", value)

	// Combine the formatted value and unit back into a string
	formattedDuration := formattedValue + unitStr

	return formattedDuration
}

const (
	green   = "\033[97;42m"
	white   = "\033[90;47m"
	yellow  = "\033[90;43m"
	red     = "\033[97;41m"
	blue    = "\033[97;44m"
	magenta = "\033[97;45m"
	cyan    = "\033[97;46m"
	reset   = "\033[0m"
)

func (p *LogFormatterParams) StatusCodeColor() string {
	code := p.StatusCode

	switch {
	case code >= http.StatusContinue && code < http.StatusOK:
		return white
	case code >= http.StatusOK && code < http.StatusMultipleChoices:
		return green
	case code >= http.StatusMultipleChoices && code < http.StatusBadRequest:
		return white
	case code >= http.StatusBadRequest && code < http.StatusInternalServerError:
		return yellow
	default:
		return red
	}
}

func (p *LogFormatterParams) MethodColor() string {
	method := p.Method

	switch method {
	case http.MethodConnect:
		return white
	case http.MethodGet:
		return blue
	case http.MethodPost:
		return cyan
	case http.MethodPut:
		return yellow
	case http.MethodDelete:
		return red
	case http.MethodPatch:
		return green
	case http.MethodHead:
		return magenta
	case http.MethodOptions:
		return white
	default:
		return reset
	}
}

func (p *LogFormatterParams) ResetColor() string {
	return reset
}

func (p *LogFormatterParams) IsOutputColor() bool {
	return isatty.IsTerminal(os.Stdin.Fd()) && isatty.IsTerminal(os.Stdout.Fd())
}

type LogWriter struct {
	writer io.Writer
}

func (lw *LogWriter) Write(p []byte) (n int, err error) {
	return lw.writer.Write(p)
}

func GinLoggerMiddleware(next http.Handler) http.Handler {
	formatter := defaultLogFormatter
	logWriter := &LogWriter{writer: os.Stdout}
	skipPaths := make([]string, 0)

	return &GinLogger{
		next:      next,
		formatter: formatter,
		logWriter: logWriter,
		skipPaths: skipPaths,
	}
}

func (m *GinLogger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	path := r.URL.Path
	raw := r.URL.RawQuery
	iw := &Interceptor{ResponseWriter: w}

	m.next.ServeHTTP(iw, r)

	status := iw.StatusCode

	if status == 0 {
		status = 200
	}

	if !m.shouldSkipPath(path) {
		param := LogFormatterParams{
			Request:    r,
			TimeStamp:  time.Now(),
			Latency:    time.Since(start),
			ClientIP:   getClientIP(r),
			Host:       r.Host,
			Method:     r.Method,
			Proto:      r.Proto,
			StatusCode: status,
			Path:       path,
		}

		if raw != "" {
			param.Path = path + "?" + raw
		}

		fmt.Fprint(m.logWriter, m.formatter(param))
	}
}

func (m *GinLogger) shouldSkipPath(path string) bool {
	for _, p := range m.skipPaths {
		if p == path {
			return true
		}
	}
	return false
}

func getClientIP(r *http.Request) string {
	// Retrieve the client IP address from the request headers
	clientIP := r.Header.Get("X-Real-IP")
	if clientIP == "" {
		clientIP = r.Header.Get("X-Forwarded-For")
		if clientIP == "" {
			clientIP = r.RemoteAddr
		}
	}
	return clientIP
}
