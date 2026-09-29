package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/reyhardy/exs-prescription-app/account"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrWrongEmailOrPassword = errors.New("wrong email or password")
	ErrInactiveUser         = errors.New("inactive user")
	ErrIdNotProvided        = errors.New("user Id not provided")
	ErrUnauthorized         = errors.New("unauthorized id")
)

type UserRepo interface {
	InsertAccount(ctx context.Context, u *account.Detail) error
	GetAccountByEmail(ctx context.Context, email string) (*account.Detail, error)
	UpdateAccountLastLogin(ctx context.Context, id string) error
	GetAccountById(ctx context.Context, id string) (*account.Detail, error)
}

type User struct {
	r UserRepo
}

func NewUser(r UserRepo) *User {
	return &User{
		r: r,
	}
}

func (s *User) CreateNewUser(ctx context.Context, params *account.DetailParams) error {
	u, err := account.New(params)
	if err != nil {
		return err
	}
	if err := s.r.InsertAccount(ctx, u); err != nil {
		return err
	}
	return nil
}

func (s *User) AuthorizeUser(ctx context.Context, params *account.DetailParams) (*account.Detail, error) {
	u, err := s.r.GetAccountByEmail(ctx, params.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrWrongEmailOrPassword
		}
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(params.Password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, ErrWrongEmailOrPassword
		}
		return nil, err
	}
	if err := s.r.UpdateAccountLastLogin(ctx, u.Id); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *User) GetUser(ctx context.Context, id string) (*account.Detail, error) {
	if id == "" {
		return nil, ErrIdNotProvided
	}
	u, err := s.r.GetAccountById(ctx, id)
	if err != nil {
		return nil, err
	}
	return u, nil
}
