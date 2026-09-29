package middleware

import (
	"context"
	"net/http"
	"slices"

	"github.com/reyhardy/exs-prescription-app/account"
	"github.com/reyhardy/exs-prescription-app/internal/ctx"
	"github.com/reyhardy/exs-prescription-app/session"
)

type sessionValidator interface {
	ValidateSession(ctx context.Context, id string) (*session.Session, error)
}

type userGetter interface {
	GetUser(ctx context.Context, id string) (*account.Detail, error)
}

type Middleware struct {
	ug userGetter
	sv sessionValidator
}

func New(ug userGetter, sv sessionValidator) Middleware {
	return Middleware{
		ug: ug,
		sv: sv,
	}
}

func (mw Middleware) Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("session_token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		s, err := mw.sv.ValidateSession(r.Context(), c.Value)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		u, err := mw.ug.GetUser(r.Context(), s.AccountId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		uCtx := ctx.MapUserCtx(u)
		ctx := ctx.WithContext(r.Context(), uCtx)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (mw Middleware) RBAC(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			u := ctx.GetContext(r.Context())
			if u == nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			if !slices.Contains(u.Roles, account.Role(role)) {
				http.Error(w, "role not matched", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
