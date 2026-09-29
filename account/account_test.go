package account_test

import (
	"testing"
	"time"

	"github.com/reyhardy/exs-prescription-app/account"
)

func TestCreateNewAccount(t *testing.T) {
	tests := map[string]struct {
		params         *account.DetailParams
		wantIdNotEmpty bool
		wantFirstName  string
		wantEmail      string
		wantErr        error
	}{
		"success create account": {
			params: &account.DetailParams{
				FirstName: "Juan",
				Email:     "juan@mail.com",
				BirthDay:  time.Date(1990, 10, 10, 0, 0, 0, 0, time.UTC),
				Gender:    string(account.GenderMale),
				Roles:     []string{string(account.RoleTherapist)},
			},
			wantIdNotEmpty: true,
			wantFirstName:  "Juan",
			wantEmail:      "juan@mail.com",
			wantErr:        nil,
		},
		"error invalid params": {
			params: &account.DetailParams{
				FirstName: "Juan",
				Email:     "",
				BirthDay:  time.Date(1990, 10, 10, 0, 0, 0, 0, time.UTC),
				Gender:    string(account.GenderMale),
				Roles:     []string{string(account.RoleTherapist)},
			},
			wantIdNotEmpty: false,
			wantErr:        account.ErrInvalidParams,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			u, err := account.New(tt.params)
			if got, want := err, tt.wantErr; got != want {
				t.Errorf(" got %v, want %v", got, want)
			}
			if u != nil {
				if got, want := u.Id != "", tt.wantIdNotEmpty; got != want {
					t.Errorf(" got %v, want %v", got, want)
				}
				if got, want := u.FirstName, tt.wantFirstName; got != want {
					t.Errorf(" got %v, want %v", got, want)
				}
				if got, want := u.Email, tt.wantEmail; got != want {
					t.Errorf(" got %v, want %v", got, want)
				}
			}
		})
	}
}
