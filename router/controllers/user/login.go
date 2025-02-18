package user

import (
	"net/http"

	"github.com/bhbdev/jam/config"
	"github.com/bhbdev/jam/lib/page"
	"github.com/bhbdev/jam/models"
)

func Login(w http.ResponseWriter, r *http.Request) {
	tpl := "user/login"

	p := page.New(r.Context())
	p.SetPageActive("/login")
	p.SetAutoWrap(false)
	p.Data["FormAction"] = r.URL.Path

	if r.Method == http.MethodGet {
		w.Header().Set("HX-Redirect", "/login")
		p.Render(w, tpl)
		return
	}

	var user models.User

	// handle POST
	user.Email = r.PostFormValue("email")
	user.Password = r.PostFormValue("password")

	sessionId, errs := p.Services.UserService.Login(w, r, &user)
	// validate
	for key, err := range errs {
		p.AddError(key, err)
	}

	if p.HasErrors() {
		p.Render(w, tpl)
		return
	}

	cookie := &http.Cookie{
		Name:  config.SessionCookieName,
		Value: sessionId,
		Path:  "/",
	}
	http.SetCookie(w, cookie)
	w.Header().Set("HX-Redirect", "/apps")
	//http.Redirect(w, r, "/apps", http.StatusSeeOther)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	p := page.New(r.Context())
	p.Services.UserService.Logout(w, r)

	cookie := &http.Cookie{
		Name:   config.SessionCookieName,
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
	//w.Header().Set("HX-Redirect", "/login")
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
