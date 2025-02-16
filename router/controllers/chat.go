package controllers

import (
	"bhbdev/jam/lib/page"
	"net/http"
)

func Chat(w http.ResponseWriter, r *http.Request) {
	p := page.New(r.Context())
	p.Render(w, "chat")
}
