package router

import (
	"context"
	"net/http"

	"github.com/bhbdev/jam/router/controllers"
	"github.com/bhbdev/jam/router/controllers/chat"
	"github.com/bhbdev/jam/router/controllers/jobapp"
	"github.com/bhbdev/jam/router/controllers/resume"
	"github.com/bhbdev/jam/router/controllers/user"
)

type Mux interface {
	Handle(pattern string, handler http.Handler)
	HandleFunc(pattern string, handler http.HandlerFunc)
}

func Apply(ctx context.Context, mux Mux) {
	mux.HandleFunc("/", controllers.Index)

	mux.HandleFunc("GET /register", user.Register)
	mux.HandleFunc("POST /register", user.Register)
	//mux.HandleFunc("GET /register/verify/{token}", user.RegisterVerify)

	mux.HandleFunc("GET /login", user.Login)
	mux.HandleFunc("POST /login", user.Login)
	mux.HandleFunc("GET /logout", user.Logout)

	mux.HandleFunc("GET /admin/user/list", user.UserList)
	mux.HandleFunc("GET /admin/user/add", user.UserForm)
	mux.HandleFunc("POST /admin/user/add", user.UserForm)
	mux.HandleFunc("GET /admin/user/edit/{id}", user.UserForm)
	mux.HandleFunc("POST /admin/user/edit/{id}", user.UserForm)

	mux.HandleFunc("GET /apps", jobapp.JobAppList)
	mux.HandleFunc("GET /apps/new", jobapp.JobAppForm)
	mux.HandleFunc("POST /apps/new", jobapp.JobAppForm)
	mux.HandleFunc("GET /apps/edit/{id}", jobapp.JobAppForm)
	mux.HandleFunc("POST /apps/edit/{id}", jobapp.JobAppForm)
	mux.HandleFunc("DELETE /apps/delete/{id}", jobapp.JobAppDelete)

	mux.HandleFunc("GET /resumes", resume.ResumeList)
	mux.HandleFunc("GET /resume/upload", resume.Upload)
	mux.HandleFunc("POST /resume/upload", resume.Upload)
	mux.HandleFunc("DELETE /resume/upload/{id}", resume.UploadDelete)

	ws := chat.NewWebSocketHandler()
	mux.HandleFunc("/ws", ws.HandleFunc)
	mux.HandleFunc("/chat", chat.Chat)
}

func Assets(mux Mux, path string) {
	mux.Handle("GET /assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(path))))
}
