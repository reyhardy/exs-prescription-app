package service_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/reyhardy/exs-prescription-app/internal/service"
	"github.com/reyhardy/exs-prescription-app/session"
)

var (
	errSessionStore = errors.New("err db")
)

func TestCreateSession(t *testing.T) {
	tests := map[string]struct {
		userId         string
		d              time.Duration
		dbErr          bool
		wantErr        error
		wantNilSession bool
	}{
		"success": {
			userId:         "test12345",
			d:              3 * time.Second,
			wantErr:        nil,
			wantNilSession: false,
		},
		"id not provided": {
			userId:         "",
			d:              3 * time.Second,
			wantErr:        session.ErrInvalidUSerId,
			wantNilSession: true,
		},
		"session duration is zero": {
			userId:         "test12345",
			d:              0,
			wantErr:        session.ErrZeroDurationProvided,
			wantNilSession: true,
		},
		"db err": {
			userId:         "test12345",
			d:              3 * time.Second,
			dbErr:          true,
			wantErr:        errSessionStore,
			wantNilSession: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			sm := service.NewSession(&fakeSessionRepo{
				sessionMap: make(map[string]*session.Session),
				dbErr:      tt.dbErr,
			}, tt.d)
			s, err := sm.CreateSession(t.Context(), tt.userId)
			if got, want := err, tt.wantErr; got != want {
				t.Errorf("got %v, want %v", got, want)
			}
			if got, want := s == nil, tt.wantNilSession; got != want {
				t.Errorf("got %v, want %v", got, want)
			}
		})
	}
}

func TestValidateSession(t *testing.T) {
	tests := map[string]struct {
		sessionId      string
		ttl            time.Duration
		dbErr          bool
		wantErr        error
		wantNilSession bool
	}{
		"success": {
			sessionId:      "session12345",
			ttl:            10 * time.Second,
			wantErr:        nil,
			wantNilSession: false,
		},
		"error session expired": {
			sessionId:      "session12345",
			ttl:            3 * time.Second,
			wantErr:        service.ErrSessionAlreadyExpired,
			wantNilSession: true,
		},
		"db err": {
			sessionId:      "session12345",
			ttl:            3 * time.Second,
			dbErr:          true,
			wantErr:        errSessionStore,
			wantNilSession: true,
		},
	}
	for name, tt := range tests {
		sr := &fakeSessionRepo{
			sessionMap: map[string]*session.Session{
				"session12345": {
					Id:        "session12345",
					AccountId: "user12345",
					ExpiredAt: time.Now().Add(tt.ttl),
				},
			},
			dbErr: tt.dbErr,
		}
		t.Run(name, func(t *testing.T) {
			sm := service.NewSession(sr, time.Hour)

			time.Sleep(4 * time.Second)
			s, err := sm.ValidateSession(t.Context(), tt.sessionId)
			if got, want := err, tt.wantErr; got != want {
				t.Errorf("got %v want %v", got, want)
			}
			if got, want := s == nil, tt.wantNilSession; got != want {
				t.Errorf("got %v want %v", got, want)
			}
		})
	}
}

func TestDeleteSession(t *testing.T) {
	tests := map[string]struct {
		dbErr               bool
		sessionId           string
		wantErr             error
		wantSessionNotExist bool
	}{
		"no error session exist": {
			sessionId:           "session12345",
			wantErr:             nil,
			wantSessionNotExist: true,
		},
		"no error session not exists in store": {
			sessionId:           "sessionnotexists",
			wantErr:             nil,
			wantSessionNotExist: true,
		},
		"db error": {
			dbErr:               true,
			sessionId:           "session12345",
			wantErr:             errSessionStore,
			wantSessionNotExist: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			sessionRepo := &fakeSessionRepo{
				sessionMap: map[string]*session.Session{
					"session12345": {
						Id:        "session12345",
						AccountId: "user12345",
						ExpiredAt: time.Now().Add(3 * time.Second),
					},
				},
				dbErr: tt.dbErr,
			}
			sm := service.NewSession(sessionRepo, 3*time.Second)
			err := sm.DeleteSession(t.Context(), tt.sessionId)
			if got, want := err, tt.wantErr; got != want {
				t.Errorf("got %v want %v", got, want)
			}
			_, exists := sessionRepo.sessionMap[tt.sessionId]
			if got, want := !exists, tt.wantSessionNotExist; got != want {
				t.Errorf("got %v want %v", got, want)
			}
		})
	}
}

type fakeSessionRepo struct {
	sessionMap map[string]*session.Session
	dbErr      bool
}

func (r *fakeSessionRepo) InsertSession(ctx context.Context, s *session.Session) error {
	if r.dbErr {
		return errSessionStore
	}
	r.sessionMap[s.Id] = s
	return nil
}

func (r *fakeSessionRepo) GetSessionById(ctx context.Context, id string) (*session.Session, error) {
	if r.dbErr {
		return nil, errSessionStore
	}
	s, ok := r.sessionMap[id]
	if !ok {
		return nil, fmt.Errorf("session id not exists")
	}
	return s, nil
}

func (r *fakeSessionRepo) DeleteSession(ctx context.Context, id string) error {
	if r.dbErr {
		return errSessionStore
	}
	delete(r.sessionMap, id)
	return nil
}
