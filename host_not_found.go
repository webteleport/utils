package utils

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
)

const NOT_FOUND_TEMPLATE = `<!DOCTYPE html>
<html lang="en" charset="utf-8">
<head>
        <title>{{.Title}}</title>
        <meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1, user-scalable=no">
</head>
<body>
        <h2>{{.Title}}</h2>
        <p>{{.Message}}</p>
        <code>{{.Code}}</code>
        <p>
          <a href="{{.Link}}">webteleport/ufo</a>
        </p>
</body>
</html>
`

type NotFoundData struct {
	Title   string
	Message string
	Code    string
	Link    string
}

func HostNotFoundHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.New("404").Parse(NOT_FOUND_TEMPLATE)
		if err != nil {
			log.Fatalln(err)
		}
		data := NotFoundData{
			Title:   fmt.Sprintf("🙈 host not found: %s", r.Host),
			Message: `You can teleport your local app to this site. Try:`,
			Code:    `$ ufo teleport https://ufo.k0s.io http://127.0.0.1:3000`,
			Link:    "https://github.com/webteleport/ufo",
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
		tmpl.Execute(w, data)
	})
}
