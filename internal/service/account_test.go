package service_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/reyhardy/exs-prescription-app/account"
	"github.com/reyhardy/exs-prescription-app/internal/service"
	"golang.org/x/crypto/bcrypt"
)

var (
	errAccountStore = errors.New("err db")
)

func TestAuthorizeUser(t *testing.T) {
	hashed, err := bcrypt.GenerateFromPassword([]byte("test12345"), bcrypt.MinCost)
	if err != nil {
		t.Fatal(err)
	}

	tests := map[string]struct {
		params  *account.DetailParams
		dbErr   bool
		wantErr error
		wantId  string
	}{
		"successful user login": {
			params: &account.DetailParams{
				Email:    "test@mail.com",
				Password: "test12345",
			},
			wantErr: nil,
			wantId:  "1234567890",
		},
		"email not exist": {
			params: &account.DetailParams{
				Email:    "notexist@mail.com",
				Password: "test12345",
			},
			wantErr: service.ErrWrongEmailOrPassword,
		},
		"wrong password": {
			params: &account.DetailParams{
				Email:    "test@mail.com",
				Password: "wrongpassword12345",
			},
			wantErr: service.ErrWrongEmailOrPassword,
		},
		"db error": {
			params: &account.DetailParams{
				Email:    "test@mail.com",
				Password: "wrongpassword12345",
			},
			dbErr:   true,
			wantErr: errAccountStore,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ur := &fakeAccountRepo{
				accountMap: map[string]*account.Detail{
					"test@mail.com": {Id: "1234567890", Email: "test@mail.com", Password: string(hashed)},
				},
				dbErr: tt.dbErr,
			}
			us := service.NewUser(ur)
			u, err := us.AuthorizeUser(t.Context(), tt.params)
			if got, want := err, tt.wantErr; got != want {
				t.Errorf("got %v, want %v", got, want)
			}
			if u == nil {
				return
			}
			if got, want := u.Id, tt.wantId; got != want {
				t.Errorf("got %v, want %v", got, want)
			}
		})
	}
}

func TestGetUser(t *testing.T) {
	tests := map[string]struct {
		id            string
		wantErr       error
		wantFirstName string
	}{
		"success": {
			id:            "1234567890",
			wantErr:       nil,
			wantFirstName: "test",
		},
		"error wrong id": {
			id:      "wrongid",
			wantErr: sql.ErrNoRows,
		},
		"error id not provided": {
			id:      "",
			wantErr: service.ErrIdNotProvided,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ur := &fakeAccountRepo{
				accountMap: map[string]*account.Detail{
					"1234567890": {Id: "1234567890", FirstName: "test"},
				},
			}
			us := service.NewUser(ur)
			u, err := us.GetUser(t.Context(), tt.id)
			if got, want := err, tt.wantErr; got != want {
				t.Errorf("got %v want %v", got, want)
			}
			if err != nil {
				t.Log("err not nil")
				return
			}
			if got, want := u.FirstName, tt.wantFirstName; got != want {
				t.Errorf("got %v want %v", got, want)
			}
		})
	}
}

type fakeAccountRepo struct {
	accountMap map[string]*account.Detail
	dbErr      bool
}

func (r *fakeAccountRepo) InsertAccount(ctx context.Context, u *account.Detail) error {
	r.accountMap[u.Id] = u
	return nil
}

func (r *fakeAccountRepo) GetAccountByEmail(ctx context.Context, email string) (*account.Detail, error) {
	if r.dbErr {
		return nil, errAccountStore
	}
	u, ok := r.accountMap[email]
	if !ok {
		return nil, sql.ErrNoRows
	}
	return u, nil
}

func (r *fakeAccountRepo) UpdateAccountLastLogin(ctx context.Context, id string) error {
	if r.dbErr {
		return errAccountStore
	}
	return nil
}

func (r *fakeAccountRepo) GetAccountById(ctx context.Context, id string) (*account.Detail, error) {
	u, ok := r.accountMap[id]
	if !ok {
		return nil, sql.ErrNoRows
	}
	return u, nil
}
