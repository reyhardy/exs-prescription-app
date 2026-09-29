package session

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"time"
)

var (
	ErrInvalidUSerId        = errors.New("invalid user id")
	ErrZeroDurationProvided = errors.New("d must not be zero")
)

type Session struct {
	Id        string
	AccountId string
	ExpiredAt time.Time
}

func New(userId string, d time.Duration) (*Session, error) {
	if userId == "" {
		return nil, ErrInvalidUSerId
	}
	if d == 0 {
		return nil, ErrZeroDurationProvided
	}
	id, err := generateToken()
	if err != nil {
		return nil, err
	}
	return &Session{
		Id:        id,
		AccountId: userId,
		ExpiredAt: time.Now().Add(d),
	}, nil
}

func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiredAt)
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	id := base64.RawURLEncoding.EncodeToString(b)
	return id, nil
}
