package ctx

import (
	"context"

	"github.com/reyhardy/exs-prescription-app/account"
)

type contextKey struct{}

type UserCtx struct {
	Id    string
	Roles account.Roles
}

func MapUserCtx(u *account.Detail) *UserCtx {
	return &UserCtx{
		Id:    u.Id,
		Roles: u.Roles,
	}
}

func WithContext(ctx context.Context, u *UserCtx) context.Context {
	return context.WithValue(ctx, contextKey{}, u)
}

func GetContext(ctx context.Context) *UserCtx {
	u, ok := ctx.Value(contextKey{}).(*UserCtx)
	if !ok {
		return nil
	}
	return u
}
