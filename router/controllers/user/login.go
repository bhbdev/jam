package user

import (
	"bhbdev/jam/lib/page"
	"bhbdev/jam/models"
	"bhbdev/jam/services"
	"net/http"
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

	userService := services.NewUserService(p.Database().Instance())
	var user models.User

	// handle POST
	user.Email = r.PostFormValue("email")
	user.Password = r.PostFormValue("password")

	errs := userService.Login(r.Context(), &user)
	// validate
	for key, err := range errs {
		p.AddError(key, err)
	}

	if p.HasErrors() {
		p.Render(w, tpl)
		return
	}
	w.Header().Set("HX-Redirect", "/apps")
	//http.Redirect(w, r, "/apps", http.StatusSeeOther)
}
