package controllers

import (
	"net/http"

	"github.com/bhbdev/jam/lib/page"
)

func Chat(w http.ResponseWriter, r *http.Request) {
	p := page.New(r.Context())
	p.Render(w, "chat")
}
