package middleware_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"

	"github.com/reyhardy/exs-prescription-app/account"
	"github.com/reyhardy/exs-prescription-app/internal/ctx"
	"github.com/reyhardy/exs-prescription-app/internal/handler/middleware"
	"github.com/reyhardy/exs-prescription-app/session"
)

func TestAuthMiddleware(t *testing.T) {
	tests := map[string]struct {
		sessionId         string
		svMock            func(ctx context.Context, id string) (*session.Session, error)
		ugMock            func(ctx context.Context, id string) (*account.Detail, error)
		wantCode          int
		wantUserId        string
		wantRedicrectPath string
		wantNextCalled    bool
	}{
		"auth success": {
			sessionId: "sessionexist12345",
			svMock: func(ctx context.Context, id string) (*session.Session, error) {
				if id != "sessionexist12345" {
					return nil, fmt.Errorf("incorrect session id")
				}
				return &session.Session{Id: "sessionexist12345", AccountId: "userexist12345"}, nil
			},
			ugMock: func(ctx context.Context, id string) (*account.Detail, error) {
				if id != "userexist12345" {
					return nil, fmt.Errorf("incorrect user id")
				}
				return &account.Detail{Id: "userexist12345"}, nil
			},
			wantCode:       http.StatusOK,
			wantUserId:     "userexist12345",
			wantNextCalled: true,
		},
		"redirect on session not exist": {
			sessionId: "sessionexist12345",
			svMock: func(ctx context.Context, id string) (*session.Session, error) {
				return nil, fmt.Errorf("session not exist")
			},
			wantCode:          http.StatusSeeOther,
			wantRedicrectPath: "/login",
			wantNextCalled:    false,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			var userId string
			nextCalled := false
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				nextCalled = true
				u := ctx.GetContext(r.Context())
				if u == nil {
					t.Fatal("error getting user from context")
				}
				userId = u.Id
				w.WriteHeader(http.StatusOK)
			})
			req := httptest.NewRequest(http.MethodGet, "/welcome", nil)
			rr := httptest.NewRecorder()
			sv := &mockSessionValidator{
				svFn: tt.svMock,
			}
			ug := &mockUserGetter{
				ugFn: tt.ugMock,
			}
			req.AddCookie(&http.Cookie{
				Name:  "session_token",
				Value: tt.sessionId,
			})
			mw := middleware.New(ug, sv)
			testedMw := mw.Authentication(nextHandler)
			testedMw.ServeHTTP(rr, req)
			if got, want := rr.Code, tt.wantCode; got != want {
				t.Errorf("got %v, want %v", got, want)
			}
			if got, want := userId, tt.wantUserId; got != want {
				t.Errorf("got %v, want %v", got, want)
			}
			if got, want := rr.Header().Get("Location"), tt.wantRedicrectPath; got != want {
				t.Errorf("got %v, want %v", got, want)
			}
			if got, want := nextCalled, tt.wantNextCalled; got != want {
				t.Errorf("got %v, want %v", got, want)
			}
		})
	}
}

func TestRBACMiddleware(t *testing.T) {
	tests := map[string]struct {
		requiredRole     string
		user             *account.Detail
		wantCode         int
		wantIsNextCalled bool
		wantUserRole     []account.Role
	}{
		"authorized user": {
			requiredRole: "therapist",
			user: &account.Detail{
				Id:    "user12345",
				Roles: []account.Role{account.RoleTherapist},
			},
			wantCode:         http.StatusOK,
			wantIsNextCalled: true,
			wantUserRole:     []account.Role{account.RoleTherapist},
		},
		"unauthorized user": {
			requiredRole: "therapist",
			user: &account.Detail{
				Id:    "user12345",
				Roles: []account.Role{account.RolePatient},
			},
			wantCode:         http.StatusForbidden,
			wantIsNextCalled: false,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			var roles []account.Role
			nextCalled := false
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				nextCalled = true
				u := ctx.GetContext(r.Context())
				if u == nil {
					t.Fatal("error getting user from context")
				}
				roles = u.Roles
				w.WriteHeader(http.StatusOK)
			})
			mw := middleware.New(nil, nil)
			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			ctx := ctx.WithContext(req.Context(), ctx.MapUserCtx(tt.user))
			testedMW := mw.RBAC(tt.requiredRole)(nextHandler)
			testedMW.ServeHTTP(rr, req.WithContext(ctx))
			if got, want := rr.Code, tt.wantCode; got != want {
				t.Errorf("got %v, want %v", got, want)
			}
			if got, want := nextCalled, tt.wantIsNextCalled; got != want {
				t.Errorf("got %v, want %v", got, want)
			}
			if got, want := slices.Equal(roles, tt.wantUserRole), true; got != want {
				t.Errorf("got %v, want %v", got, want)
			}
		})
	}
}

type mockSessionValidator struct {
	svFn func(ctx context.Context, id string) (*session.Session, error)
}

func (sv *mockSessionValidator) ValidateSession(ctx context.Context, id string) (*session.Session, error) {
	return sv.svFn(ctx, id)
}

type mockUserGetter struct {
	ugFn func(ctx context.Context, id string) (*account.Detail, error)
}

func (ug *mockUserGetter) GetUser(ctx context.Context, id string) (*account.Detail, error) {
	return ug.ugFn(ctx, id)
}
