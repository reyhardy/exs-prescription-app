package repo

import (
	"context"

	"github.com/reyhardy/exs-prescription-app/internal/db"
	"github.com/reyhardy/exs-prescription-app/session"
)

type Session struct {
	db *db.DB
}

func NewSession(db *db.DB) *Session {
	return &Session{
		db: db,
	}
}

func (r *Session) InsertSession(ctx context.Context, s *session.Session) error {
	q := `
		INSERT INTO sessions
		(id, account_id, expired_at)
		VALUES (?, ?, ?)
	`
	return r.db.Exec(ctx, q, s.Id, s.AccountId, s.ExpiredAt)
}

func (r *Session) GetSessionById(ctx context.Context, id string) (*session.Session, error) {
	q := `
		SELECT id, account_id, expired_at FROM sessions 
		WHERE id = ?
	`
	var s session.Session
	if err := r.db.QueryRow(ctx, q, &s, id); err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *Session) DeleteSession(ctx context.Context, id string) error {
	q := `
		DELETE FROM sessions WHERE id = ?;
	`
	if err := r.db.Exec(ctx, q, id); err != nil {
		return err
	}
	return nil
}
