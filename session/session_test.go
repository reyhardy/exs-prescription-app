package session_test

import (
	"testing"
	"time"

	"github.com/reyhardy/exs-prescription-app/session"
)

func TestNewSession(t *testing.T) {
	tests := map[string]struct {
		userId  string
		d       time.Duration
		wantErr error
	}{
		"success": {
			userId:  "user12345",
			d:       3 * time.Second,
			wantErr: nil,
		},
		"error empty id": {
			userId:  "",
			d:       3 * time.Second,
			wantErr: session.ErrInvalidUSerId,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := session.New(tt.userId, tt.d)
			if got, want := err, tt.wantErr; got != want {
				t.Errorf("got %v, want %v", got, want)
			}
		})
	}
}

func TestSessionExpired(t *testing.T) {
	tests := map[string]struct {
		expiredAt     time.Time
		wantIsExpired bool
	}{
		"is expired": {
			expiredAt:     time.Now().Add(-3 * time.Second),
			wantIsExpired: true,
		},
		"is not expired": {
			expiredAt:     time.Now().Add(3 * time.Second),
			wantIsExpired: false,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			s := &session.Session{
				ExpiredAt: tt.expiredAt,
			}
			if got, want := s.IsExpired(), tt.wantIsExpired; got != want {
				t.Errorf("got %v, want %v", got, want)
			}
		})
	}
}
