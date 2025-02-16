package controllers

import (
	"net/http"

	"github.com/bhbdev/jam/lib/page"
)

func Index(w http.ResponseWriter, r *http.Request) {
	p := page.New(r.Context())
	p.SetPageActive(r.URL.Path)
	p.Render(w, "index")
}
