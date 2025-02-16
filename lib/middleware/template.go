package middleware

import (
	"bhbdev/jam/lib/logger"
	"context"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"time"
)

const (
	ContextTemplate = ContextKey("template")
	//ContextPartial  = ContextKey("partial")
)

type Template struct {
	handler http.Handler
	path    string
}

func NewTemplate(path string) *Template {
	return &Template{
		path: path,
	}
}

func (m *Template) SetHandler(handler http.Handler) {
	m.handler = handler
}

func (m *Template) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	patterns := []string{
		filepath.Join(m.path, "*.html"),
		filepath.Join(m.path, "*", "*.html"),
	}
	errs := make([]error, 0, len(patterns))

	// add custom functions to the template
	funcs := template.FuncMap{
		"now": time.Now,
		"add": func(a, b int) int {
			return a + b
		},
		"sub": func(a, b int) int {
			return a - b
		},
		"pathext": func(h string, n string) bool {
			logger.Info("pathext", "h", h, "n", n, "ext", filepath.Ext(h))
			return filepath.Ext(h) == n
		},
	}

	var err error
	tpl := template.New("")
	for _, pattern := range patterns {
		if tpl, err = tpl.Funcs(funcs).ParseGlob(pattern); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		allErrs := errors.Join(errs...)
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "Errors while processing templates:\n%v\n", allErrs)
		return
	}

	ctx := context.WithValue(r.Context(), ContextTemplate, tpl)
	m.handler.ServeHTTP(w, r.WithContext(ctx))
}
