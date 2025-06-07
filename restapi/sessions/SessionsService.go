package sessions

type SessionService struct {
	Store SessionStore
}

func NewSessionService(store SessionStore) *SessionService {
	return &SessionService{
		Store: store,
	}
}

func (s *SessionService) CreateNewSession(userId int) {
	s.Store.CreateNewSession(userId)
}
