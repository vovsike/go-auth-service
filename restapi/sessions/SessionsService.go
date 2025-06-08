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
		SessionId: sid.String(),
		UserId:    userId,
		Expires:   time.Now().Add(time.Hour * 24),
	}
	return s.Store.CreateNewSession(session)
}

func (s *SessionService) VerifySession(sessionId string) (Session, bool) {
	sesh, err := s.Store.GetSession(sessionId)
	if err != nil {
		return Session{}, false
	}
	return sesh, true
}
