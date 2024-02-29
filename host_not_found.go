package utils

import (
	"net/http"
	"text/template"
)

const NOT_FOUND_TEMPLATE = `Not found: {{.Host}}`

type NotFoundData struct {
	Host string
}

func HostNotFoundHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.New("HostNotFound").Parse(NOT_FOUND_TEMPLATE))
		data := NotFoundData{
			Host: r.Host,
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
		tmpl.Execute(w, data)
	})
}
