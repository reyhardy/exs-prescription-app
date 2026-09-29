package repo

import (
	"context"

	"github.com/reyhardy/exs-prescription-app/account"
	"github.com/reyhardy/exs-prescription-app/internal/db"
)

type Account struct {
	db *db.DB
}

func NewAccount(db *db.DB) *Account {
	return &Account{
		db: db,
	}
}

func (r *Account) InsertAccount(ctx context.Context, u *account.Detail) error {
	q := `
		INSERT INTO accounts 
		(id, first_name, last_name, email, password, birthday, age, gender, roles) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);
	`
	return r.db.Exec(ctx, q, u.Id, u.FirstName, u.LastName, u.Email, u.Password, u.Birthday, u.Age, u.Gender, u.Roles)
}

func (r *Account) GetAccountByEmail(ctx context.Context, email string) (*account.Detail, error) {
	q := `
		SELECT id, email, password FROM accounts WHERE email = ?
	`
	var u account.Detail
	if err := r.db.QueryRow(ctx, q, &u, email); err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *Account) UpdateAccountLastLogin(ctx context.Context, id string) error {
	q := `
		UPDATE accounts SET
		last_login_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`
	return r.db.Exec(ctx, q, id)
}

func (r *Account) GetAccountById(ctx context.Context, id string) (*account.Detail, error) {
	q := `
		SELECT a.id, a.first_name, a.last_name, a.email, a.birthday, a.age, a.gender, a.status, a.last_login_at, a.created_at, a.updated_at, group_concat(r.roles) as roles
		FROM accounts AS a
		LEFT JOIN account_roles AS r
		ON a.id = r.user_id
		WHERE a.id = ?;
	`
	var u account.Detail
	if err := r.db.QueryRow(ctx, q, &u, id); err != nil {
		return nil, err
	}
	return &u, nil
}
