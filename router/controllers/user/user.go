package user

import (
	"net/http"
	"strconv"

	"github.com/bhbdev/jam/lib/logger"
	"github.com/bhbdev/jam/lib/page"
	"github.com/bhbdev/jam/lib/pagination"
	"github.com/bhbdev/jam/models"
)

// const (
// 	UserPath = "/user"
// )

func UserList(w http.ResponseWriter, r *http.Request) {
	p := page.New(r.Context())
	p.SetPageActive("/admin/user")

	params := &models.ListParams{
		Pagination: pagination.NewPagination(r),
	}

	users, err := p.Services.UserService.List(r.Context(), params)
	if err != nil {
		logger.Error("UserList error", "error", err)
	}
	logger.Info("UserList", "users!!!", len(users))

	p.Data["Users"] = users
	p.Data["Params"] = params
	p.Render(w, "user/list")
}

func UserForm(w http.ResponseWriter, r *http.Request) {
	tpl := "user/form"

	urlId := r.PathValue("id")
	isEditing := urlId != ""
	p := page.New(r.Context())
	p.SetPageActive("/admin/user")
	p.Data["FormAction"] = r.URL.Path
	p.Data["IsEditing"] = isEditing

	user := &models.User{}

	if isEditing {
		id, err := strconv.Atoi(urlId)
		if err != nil {
			logger.Error("invalid user id", "value", urlId)
			http.Redirect(w, r, "/admin/user", http.StatusSeeOther)
			return
		}

		user, err = p.Services.UserService.Get(r.Context(), id)
		if err != nil {
			logger.Error("user does not exist", "id", id)
			http.Redirect(w, r, "/admin/user", http.StatusSeeOther)
			return
		}
	}

	p.Data["User"] = user

	if r.Method == http.MethodGet {
		p.Render(w, tpl)
		return
	}

	logger.Info("UserForm", "POST", r.PostFormValue("first_name"))
	// handle POST
	user.FirstName = r.PostFormValue("first_name")
	user.LastName = r.PostFormValue("last_name")
	user.Email = r.PostFormValue("email")
	user.Status = models.UserStatus(r.PostFormValue("status"))
	user.ProfileImage = r.PostFormValue("profile_image")

	oldPassword := user.Password
	if newPassword := r.PostFormValue("password"); newPassword != "" {
		user.Password = newPassword
	}

	// validate
	for key, err := range user.Validate() {
		p.AddError(key, err)
	}

	if p.HasErrors() {
		p.Render(w, tpl)
		return
	}

	err := p.Services.UserService.Save(r.Context(), user)
	if err != nil {
		logger.Error("UserForm error saving user", "error", err)
		p.AddError("alert", "Error saving user")
		p.Render(w, tpl)
		return
	}

	if oldPassword != user.Password {
		err := p.Services.UserService.SetPassword(r.Context(), user.ID, user.Password)
		if err != nil {
			logger.Error("failed to update password", "error", err)
			p.AddError("Database", err.Error())
			p.Render(w, tpl)
			return
		}
	}

	http.Redirect(w, r, "/admin/user/list", http.StatusSeeOther)
}
