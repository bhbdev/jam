package middleware

import (
	"net/http"
	"strings"
)

var _ Middleware = new(ContentType)

type ContentType struct {
	handler http.Handler
}

func NewContentType() *ContentType {
	return &ContentType{}
}

func (m *ContentType) SetHandler(handler http.Handler) {
	m.handler = handler
}

func (m *ContentType) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.HasSuffix(r.URL.Path, ".css") {
		w.Header().Set("Content-Type", "text/css")
	} else if strings.HasSuffix(r.URL.Path, ".js") {
		w.Header().Set("Content-Type", "application/javascript")
	} else if strings.HasSuffix(r.URL.Path, ".pdf") {
		w.Header().Set("Content-Type", "application/pdf")
		w.Header().Set("Content-Disposition", "inline;filename=resume.pdf")
	} else {
		w.Header().Set("Content-Type", "text/html; charset=utf=8")
	}
	m.handler.ServeHTTP(w, r)
}
