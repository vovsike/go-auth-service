package sessions

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type ServiceDB struct {
	dbpool *pgxpool.Pool
}

func NewServiceDB(pool *pgxpool.Pool) *ServiceDB {
	return &ServiceDB{dbpool: pool}
}

func (s ServiceDB) Authenticate(ctx context.Context, userId uuid.UUID) (Session, error) {
	sid := uuid.New()
	expireAt := time.Now().Add(time.Hour * 24)
	_, err := s.dbpool.Exec(ctx, "INSERT INTO sessions (id, accountid, expiry) VALUES ($1, $2, $3)", sid, userId, expireAt)
	if err != nil {
		return Session{}, err
	}
	return Session{
		ID:        sid,
		UserID:    userId,
		ExpiresAt: expireAt,
	}, nil
}

//type SessionService struct {
//	Store SessionStore
//}
//
//func NewSessionService(store SessionStore) *SessionService {
//	return &SessionService{
//		Store: store,
//	}
//}
//
//func (s *SessionService) Authenticate(ctx context.Context, userId int) Session {
//	sid, _ := uuid.NewRandom()
//	session := Session{
//		ID:        sid.String(),
//		UserID:    userId,
//		ExpiresAt: time.Now().Add(time.Hour * 24),
//	}
//	return s.Store.CreateNewSession(ctx, session)
//}
//
//func (s *SessionService) VerifySession(ctx context.Context, sessionId string) (Session, bool) {
//	sesh, err := s.Store.GetSession(ctx, sessionId)
//	if err != nil {
//		return Session{}, false
//	}
//	return sesh, true
//}

//var (
//	ErrSessionNotFound = errors.New("session not found")
//	ErrSessionExpired  = errors.New("session expired")
//)
//
//type SessionsDB struct {
//	db *pgxpool.Pool
//}
//
//func NewSessionStoreDB(pool *pgxpool.Pool) *SessionStoreDB {
//	return &SessionStoreDB{
//		db: pool,
//	}
//}
//
//func (ss *SessionStoreDB) CreateNewSession(ctx context.Context, s Session) Session {
//	_, err := ss.db.Exec(ctx, "INSERT INTO sessions (session_id, user_id, expires) VALUES ($1, $2, $3)", s.ID, s.UserID, s.ExpiresAt)
//	if err != nil {
//		fmt.Println(err)
//		return Session{}
//	}
//	return s
//}
//
//func (ss *SessionStoreDB) GetSession(ctx context.Context, sessionId string) (Session, error) {
//	s := Session{}
//
//	err := ss.db.QueryRow(ctx, "SELECT * FROM sessions WHERE session_id = $1", sessionId).Scan(&s.ID, &s.UserID, &s.ExpiresAt)
//	if err != nil {
//		return Session{}, ErrSessionNotFound
//	}
//	if time.Now().After(s.ExpiresAt) {
//		return Session{}, ErrSessionExpired
//	}
//	return s, nil
//}
//
//func (ss *SessionStoreDB) DeleteSession(ctx context.Context, sessionId string) error {
//	if sessionId == "" {
//		return errors.New("session ID cannot be empty")
//	}
//	_, err := ss.db.Exec(ctx, "DELETE FROM sessions WHERE session_id = $1", sessionId)
//	if err != nil {
//		return fmt.Errorf("failed to delete session: %v", err)
//	}
//	return nil
//}
