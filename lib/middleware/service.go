package middleware

import (
	"bhbdev/jam/services"
	"context"
	"net/http"
)

const ContextServices = ContextKey("services")

type PageServices struct {
	svc     *services.Services
	handler http.Handler
}

func NewPageServices(services *services.Services) *PageServices {
	return &PageServices{
		svc: services,
	}
}

func (m *PageServices) SetHandler(handler http.Handler) {
	m.handler = handler
}

func (m *PageServices) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := context.WithValue(r.Context(), ContextServices, m.svc)
	m.handler.ServeHTTP(w, r.WithContext(ctx))
}
