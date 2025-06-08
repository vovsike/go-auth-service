package sessions

import (
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

func (s *SessionService) Authenticate(userId int) Session {
	sid, _ := uuid.NewRandom()
	session := Session{
		sessionId: sid.String(),
		userId:    userId,
		expires:   time.Now().Add(time.Hour * 24),
	}
	return s.Store.CreateNewSession(session)
}
