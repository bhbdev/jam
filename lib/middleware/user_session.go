package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/bhbdev/jam/config"
	"github.com/bhbdev/jam/lib/session"
	"github.com/bhbdev/jam/models"
	"github.com/bhbdev/jam/services"
)

const ContextUserSession = ContextKey("jam-user")

type UserSession struct {
	handler     http.Handler
	session     *session.Session
	userService *services.UserService
}

func NewUserSessionHandler(sess *session.Session, userService *services.UserService) *UserSession {
	return &UserSession{
		session:     sess,
		userService: userService,
	}
}

func (m *UserSession) SetHandler(handler http.Handler) {
	m.handler = handler
}

func (m *UserSession) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Skip user extraction for non-admin pages and the login pages
	path := r.URL.Path
	if strings.HasPrefix(path, "/assets/") || path == "/login" || path == "/register" {
		m.handler.ServeHTTP(w, r)
		return
	}

	var user *models.User
	// check for user session cookie
	sessionId, err := r.Cookie(config.SessionCookieName)
	if err != nil {
		http.Redirect(w, r, "/login#no-cookie", http.StatusTemporaryRedirect)
		return
	}
	// get user session
	sessionData, err := m.session.Get(r.Context(), sessionId.Value)
	if err != nil {
		http.Redirect(w, r, "/login#"+err.Error(), http.StatusTemporaryRedirect)
		return
	}
	// find user by id in session data
	user, err = m.userService.Get(r.Context(), int(sessionData.Id))
	if err != nil {
		http.Redirect(w, r, "/login#"+err.Error(), http.StatusTemporaryRedirect)
		return
	}
	if user.Status != models.UserStatusActive {
		http.Redirect(w, r, "/login#inactive", http.StatusTemporaryRedirect)
		return
	}

	// set expiration
	if err = m.session.Expire(r.Context(), sessionId.Value, config.SessionExpireTime); err != nil {
		http.Redirect(w, r, "/login#"+err.Error(), http.StatusTemporaryRedirect)
		return
	}

	ctx := context.WithValue(r.Context(), ContextUserSession, user)
	m.handler.ServeHTTP(w, r.WithContext(ctx))
}
