package user

import (
	"net/http"

	"github.com/bhbdev/jam/lib/page"
	"github.com/bhbdev/jam/models"
)

func Register(w http.ResponseWriter, r *http.Request) {
	tpl := "user/register"

	p := page.New(r.Context())
	p.SetPageActive("/register")
	p.SetAutoWrap(false)
	p.Data["FormAction"] = r.URL.Path

	if r.Method == http.MethodGet {
		p.Render(w, tpl)
		return
	}

	var user models.User

	// handle POST
	user.Username = r.PostFormValue("username")
	user.Email = r.PostFormValue("email")

	errs := p.Services.UserService.Register(r.Context(), &user)
	// validate
	for key, err := range errs {
		p.AddError(key, err)
	}

	if p.HasErrors() {
		p.Render(w, tpl)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
