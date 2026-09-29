package ctx

import (
	"context"
	"slices"
	"testing"

	"github.com/reyhardy/exs-prescription-app/account"
)

func TestWithContext(t *testing.T) {
	expectedUser := &account.Detail{
		Id:    "test12345",
		Roles: []account.Role{account.RoleTherapist},
	}
	uCtx := MapUserCtx(expectedUser)
	ctx := WithContext(context.Background(), uCtx)
	u, ok := ctx.Value(contextKey{}).(*UserCtx)
	if !ok {
		t.Fatalf("got not ok, want ok")
	}
	if got, want := u.Id, "test12345"; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := slices.Contains(u.Roles, account.RoleTherapist), true; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestGetUserContext(t *testing.T) {
	expectedUser := &account.Detail{
		Id:    "test12345",
		Roles: []account.Role{account.RoleTherapist},
	}
	uCtx := MapUserCtx(expectedUser)
	ctx := context.WithValue(context.Background(), contextKey{}, uCtx)
	u := GetContext(ctx)
	if got, want := u.Id, "test12345"; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	if got, want := slices.Contains(u.Roles, account.RoleTherapist), true; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}
