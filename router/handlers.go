package router

import (
	"bhbdev/jam/router/controllers"
	"bhbdev/jam/router/controllers/chat"
	"bhbdev/jam/router/controllers/jobapp"
	"bhbdev/jam/router/controllers/user"
	"context"
	"net/http"
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
	//mux.HandleFunc("GET /logout", user.Logout)

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

	ws := chat.NewWebSocketHandler()
	mux.HandleFunc("/ws", ws.HandleFunc)
	mux.HandleFunc("/chat", controllers.Chat)
}

func Assets(mux Mux, path string) {
	mux.Handle("GET /assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(path))))
}
