package handler_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/reyhardy/exs-prescription-app/account"
	"github.com/reyhardy/exs-prescription-app/internal/handler"
	"github.com/reyhardy/exs-prescription-app/internal/template"
	"github.com/reyhardy/exs-prescription-app/session"
	"golang.org/x/crypto/bcrypt"
)

func TestLoginUser(t *testing.T) {
	renderer := mockRenderer{
		renderFn: func(context.Context, io.Writer, template.PageName, any) error {
			return nil
		},
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte("test12345"), bcrypt.MinCost)
	if err != nil {
		t.Fatal(err)
	}
	tests := map[string]struct {
		authorizeUserFn  func(ctx context.Context, params *account.DetailParams) (*account.Detail, error)
		createSessionFn  func(ctx context.Context, userId string) (*session.Session, error)
		wantStatusCode   int
		wantRedirectTo   string
		wantCookieExist  bool
		wantCookieMaxAge int
	}{
		"success login": {
			authorizeUserFn: func(ctx context.Context, params *account.DetailParams) (*account.Detail, error) {
				if params.Email != "test@mail.com" || params.Password != "test12345" {
					return nil, fmt.Errorf("email %v or password %v is incorrect", params.Email, params.Password)
				}
				return &account.Detail{Id: "1234567890", Email: "test@mail.com", Password: string(hashed)}, nil
			},
			createSessionFn: func(ctx context.Context, userId string) (*session.Session, error) {
				if userId != "1234567890" {
					return nil, fmt.Errorf("user Id %v is incorrect", userId)
				}
				return &session.Session{Id: "sess12345678", AccountId: "1234567890", ExpiredAt: time.Now().Add(3 * time.Second)}, nil
			},
			wantStatusCode:   http.StatusFound,
			wantRedirectTo:   "/profile",
			wantCookieExist:  true,
			wantCookieMaxAge: 3,
		},
		"failure login": {
			authorizeUserFn: func(ctx context.Context, params *account.DetailParams) (*account.Detail, error) {
				return nil, fmt.Errorf("email %v or password %v is incorrect", params.Email, params.Password)
			},
			wantStatusCode:  http.StatusUnprocessableEntity,
			wantCookieExist: false,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			us := mockUserService{
				authFn: tt.authorizeUserFn,
			}
			sm := mockSessionService{
				createFn: tt.createSessionFn,
			}
			form := url.Values{}
			form.Set("email", "test@mail.com")
			form.Set("password", "test12345")
			req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(form.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			rr := httptest.NewRecorder()
			h := handler.NewAuth(us, sm, renderer)
			h.UserLogin(rr, req)
			if got, want := rr.Code, tt.wantStatusCode; got != want {
				t.Errorf("got %v, want %v", got, want)
			}
			c := getCookie(rr, "session_token")
			if got, want := c != nil, tt.wantCookieExist; got != want {
				t.Errorf("got %v want %v", got, want)
			}
			if got, want := rr.Header().Get("Location"), tt.wantRedirectTo; got != want {
				t.Errorf("got %v want %v", got, want)
			}
			if c == nil {
				return
			}
			if got, want := c.MaxAge, tt.wantCookieMaxAge; got != want {
				t.Errorf("got %v want %v", got, want)
			}
		})
	}
}

func TestLogoutUser(t *testing.T) {
	tests := map[string]struct {
		deleteFn         func(context.Context, string) error
		wantStatusCode   int
		wantRedirectPath string
		wantCookieValue  string
	}{
		"success": {
			deleteFn: func(ctx context.Context, id string) error {
				return nil
			},
			wantStatusCode:   http.StatusSeeOther,
			wantCookieValue:  "",
			wantRedirectPath: "/login",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			sm := mockSessionService{
				deleteFn: tt.deleteFn,
			}
			h := handler.NewAuth(mockUserService{}, sm, mockRenderer{})
			req := httptest.NewRequest(http.MethodPost, "/logout", nil)
			req.AddCookie(&http.Cookie{
				Name:  "session_token",
				Value: "session_id_12345",
			})
			rr := httptest.NewRecorder()
			h.UserLogout(rr, req)
			if got, want := rr.Code, tt.wantStatusCode; got != want {
				t.Errorf("got %v want %v", got, want)
			}
			if got, want := rr.Header().Get("Location"), tt.wantRedirectPath; got != want {
				t.Errorf("got %v want %v", got, want)
			}
			c := getCookie(rr, "session_token")
			if got, want := c.Value, tt.wantCookieValue; got != want {
				t.Errorf("got %v want %v", got, want)
			}
		})
	}
}

type mockUserService struct {
	authFn func(context.Context, *account.DetailParams) (*account.Detail, error)
}

func (m mockUserService) AuthorizeUser(ctx context.Context, u *account.DetailParams) (*account.Detail, error) {
	return m.authFn(ctx, u)
}

type mockSessionService struct {
	createFn func(context.Context, string) (*session.Session, error)
	deleteFn func(context.Context, string) error
}

func (m mockSessionService) CreateSession(ctx context.Context, userId string) (*session.Session, error) {
	return m.createFn(ctx, userId)
}

func (m mockSessionService) DeleteSession(ctx context.Context, id string) error {
	return m.deleteFn(ctx, id)
}

type mockRenderer struct {
	renderFn func(context.Context, io.Writer, template.PageName, any) error
}

func (m mockRenderer) Render(ctx context.Context, w io.Writer, name template.PageName, data any) error {
	return m.renderFn(ctx, w, name, data)
}

func getCookie(rr *httptest.ResponseRecorder, name string) *http.Cookie {
	for _, cookie := range rr.Result().Cookies() {
		if cookie.Name == name {
			return cookie
		}
	}
	return nil
}
