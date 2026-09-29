package handler

import (
	"context"
	"math"
	"net/http"
	"time"

	"github.com/reyhardy/exs-prescription-app/account"
	"github.com/reyhardy/exs-prescription-app/internal/template"
	"github.com/reyhardy/exs-prescription-app/session"
)

type userAuth interface {
	AuthorizeUser(context.Context, *account.DetailParams) (*account.Detail, error)
}
type sessionServicer interface {
	CreateSession(context.Context, string) (*session.Session, error)
	DeleteSession(context.Context, string) error
}

type Auth struct {
	us userAuth
	sm sessionServicer
	tr template.Renderer
}

func NewAuth(us userAuth, sm sessionServicer, tr template.Renderer) *Auth {
	return &Auth{
		us: us,
		sm: sm,
		tr: tr,
	}
}

func (h Auth) LoginPage(w http.ResponseWriter, r *http.Request) {
	h.renderLoginPage(w, r, nil)
}

func (h *Auth) UserLogin(w http.ResponseWriter, r *http.Request) {
	params := &account.DetailParams{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}
	u, err := h.us.AuthorizeUser(r.Context(), params)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		h.renderLoginPage(w, r, "Invalid Credentials")
		return
	}
	s, err := h.sm.CreateSession(r.Context(), u.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    s.Id,
		Path:     "/",
		MaxAge:   int(math.Ceil(time.Until(s.ExpiredAt).Seconds())),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
	http.Redirect(w, r, "/profile", http.StatusFound)
}

func (h *Auth) UserLogout(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := h.sm.DeleteSession(r.Context(), c.Value); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:   "session_token",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (h *Auth) renderLoginPage(w http.ResponseWriter, r *http.Request, data any) {
	if err := h.tr.Render(r.Context(), w, template.Login, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
