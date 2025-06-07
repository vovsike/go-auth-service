package sessions

import "github.com/google/uuid"

type SessionService struct {
	Store SessionStore
}

func NewSessionService(store SessionStore) *SessionService {
	return &SessionService{
		Store: store,
	}
}

func (s *SessionService) CreateNewSession(userId int) {
	sid, _ := uuid.NewRandom()
	s.Store.CreateNewSession(userId, sid.String())
}
