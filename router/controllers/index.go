package controllers

import (
	"bhbdev/jam/lib/page"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	p := page.New(r.Context())
	p.SetPageActive(r.URL.Path)
	p.Render(w, "index")
}
