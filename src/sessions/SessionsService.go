package sessions

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type SessionService struct {
	Store SessionStore
}

func NewSessionService(store SessionStore) *SessionService {
	return &SessionService{
		Store: store,
	}
}

func (s *SessionService) Authenticate(ctx context.Context, userId int) Session {
	sid, _ := uuid.NewRandom()
	session := Session{
		ID:        sid.String(),
		UserID:    userId,
		ExpiresAt: time.Now().Add(time.Hour * 24),
	}
	return s.Store.CreateNewSession(ctx, session)
}

func (s *SessionService) VerifySession(ctx context.Context, sessionId string) (Session, bool) {
	sesh, err := s.Store.GetSession(ctx, sessionId)
	if err != nil {
		return Session{}, false
	}
	return sesh, true
}
