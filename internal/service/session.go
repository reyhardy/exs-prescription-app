package service

import (
	"context"
	"errors"
	"time"

	"github.com/reyhardy/exs-prescription-app/session"
)

var (
	ErrSessionAlreadyExpired = errors.New("session already expired")
	ErrSessionNotExists      = errors.New("session not exists")
)

type sessionRepo interface {
	InsertSession(ctx context.Context, s *session.Session) error
	GetSessionById(ctx context.Context, id string) (*session.Session, error)
	DeleteSession(ctx context.Context, id string) error
}

type SessionManager struct {
	r   sessionRepo
	ttl time.Duration
}

func NewSession(r sessionRepo, ttl time.Duration) *SessionManager {
	return &SessionManager{
		r:   r,
		ttl: ttl,
	}
}

func (s *SessionManager) CreateSession(ctx context.Context, userId string) (*session.Session, error) {
	session, err := session.New(userId, s.ttl)
	if err != nil {
		return nil, err
	}
	if err := s.r.InsertSession(ctx, session); err != nil {
		return nil, err
	}
	return session, nil
}

func (s *SessionManager) ValidateSession(ctx context.Context, id string) (*session.Session, error) {
	session, err := s.r.GetSessionById(ctx, id)
	if err != nil {
		return nil, err
	}
	if session.IsExpired() {
		if err := s.DeleteSession(ctx, id); err != nil {
			return nil, err
		}
		return nil, ErrSessionAlreadyExpired
	}
	return session, nil
}

func (s *SessionManager) DeleteSession(ctx context.Context, id string) error {
	return s.r.DeleteSession(ctx, id)
}
