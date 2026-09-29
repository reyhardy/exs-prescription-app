package handler_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/reyhardy/exs-prescription-app/account"
	"github.com/reyhardy/exs-prescription-app/internal/ctx"
	"github.com/reyhardy/exs-prescription-app/internal/handler"
	"github.com/reyhardy/exs-prescription-app/internal/template"
)

func TestProfilePage(t *testing.T) {
	tests := map[string]struct {
		user           *account.Detail
		getUserFn      func(ctx context.Context, id string) (*account.Detail, error)
		renderFn       func(ctx context.Context, w io.Writer, name template.PageName, data any) error
		wantCode       int
		wantUserIdData string
	}{
		"success": {
			user: &account.Detail{
				Id:        "test12345",
				FirstName: "test",
				Roles:     []account.Role{account.RoleTherapist},
			},
			getUserFn: func(ctx context.Context, id string) (*account.Detail, error) {
				if id != "test12345" {
					return nil, fmt.Errorf("id %v not exist", id)
				}
				return &account.Detail{Id: "test12345", FirstName: "test", Roles: []account.Role{account.RoleTherapist}}, nil
			},
			renderFn: func(ctx context.Context, w io.Writer, name template.PageName, data any) error {
				return nil
			},
			wantCode:       http.StatusOK,
			wantUserIdData: "test12345",
		},
		"failure": {
			user: &account.Detail{
				Id:    "nonexist",
				Roles: []account.Role{account.RolePatient},
			},
			getUserFn: func(ctx context.Context, id string) (*account.Detail, error) {
				return nil, fmt.Errorf("user with id %v not found", id)
			},
			wantCode: http.StatusInternalServerError,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ugMock := mockUserGetter{
				getUserFn: tt.getUserFn,
			}
			renderer := &profilePageRenderer{
				renderFn: tt.renderFn,
			}
			uh := handler.NewUser(ugMock, renderer)
			req := httptest.NewRequest(http.MethodGet, "/profile", nil)
			uCtx := ctx.MapUserCtx(tt.user)
			ctx := ctx.WithContext(req.Context(), uCtx)
			rr := httptest.NewRecorder()
			uh.UserProfilePage(rr, req.WithContext(ctx))
			if got, want := rr.Code, tt.wantCode; got != want {
				t.Errorf("got %v want %v", got, want)
			}
			data, ok := renderer.data.(*account.Detail)
			if !ok {
				t.Log("no data found")
				return
			}
			if got, want := data.Id, tt.wantUserIdData; got != want {
				t.Errorf("got %v want %v", got, want)
			}
		})
	}
}

type mockUserGetter struct {
	getUserFn func(context.Context, string) (*account.Detail, error)
}

func (m mockUserGetter) GetUser(ctx context.Context, id string) (*account.Detail, error) {
	return m.getUserFn(ctx, id)
}

type profilePageRenderer struct {
	renderFn func(context.Context, io.Writer, template.PageName, any) error
	data     any
}

func (r *profilePageRenderer) Render(ctx context.Context, w io.Writer, name template.PageName, data any) error {
	r.data = data
	return r.renderFn(ctx, w, name, r.data)
}
