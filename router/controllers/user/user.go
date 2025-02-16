package user

import (
	"bhbdev/jam/lib/logger"
	"bhbdev/jam/lib/page"
	"bhbdev/jam/lib/pagination"
	"bhbdev/jam/models"
	"bhbdev/jam/services"
	"net/http"
	"strconv"
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

	db := p.Database()
	s := services.NewUserService(db.Instance())
	users, err := s.List(r.Context(), params)
	if err != nil {
		logger.Error("UserList error", "error", err)
	}
	logger.Info("UserList", "users", users)

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

	s := services.NewUserService(p.Database().Instance())
	var user *models.User

	if isEditing {
		id, err := strconv.Atoi(urlId)
		if err != nil {
			logger.Error("invalid user id", "value", urlId)
			http.Redirect(w, r, "/admin/user", http.StatusSeeOther)
			return
		}

		user, err = s.Get(r.Context(), id)
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

	// handle POST
	user.FirstName = r.PostFormValue("firstName")
	user.LastName = r.PostFormValue("lastName")
	user.Email = r.PostFormValue("email")
	user.Status = models.UserStatus(r.PostFormValue("status"))

	oldPassword := user.Password
	if newPassword := r.PostFormValue("password"); newPassword != "" {
		user.Password = newPassword
	}

	//todo: handle file upload better.. through context? middleware?
	// 	file, err := p.fs.GetUpload("profile_image", r)
	// 	if err != nil {
	// 		p.Render(w, tpl)
	// 		return
	// 	}
	// 	user.ProfileImage = p.Data["FileUpload"].(*FileUpload).FileName
	// }

	// validate
	for key, err := range user.Validate() {
		p.AddError(key, err)
	}

	if p.HasErrors() {
		p.Render(w, tpl)
		return
	}

	err := s.Save(r.Context(), user)
	if err != nil {
		logger.Error("JobAppForm error saving jobapp", "error", err)
		p.AddError("alert", "Error saving job application")
		p.Render(w, tpl)
		return
	}

	if oldPassword != user.Password {
		err := s.SetPassword(r.Context(), user.ID, user.Password)
		if err != nil {
			logger.Error("failed to update password", "error", err)
			p.AddError("Database", err.Error())
			p.Render(w, tpl)
			return
		}
	}

	http.Redirect(w, r, "/admin/user", http.StatusSeeOther)
}
