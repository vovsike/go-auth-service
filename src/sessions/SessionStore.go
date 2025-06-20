package sessions

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type Session struct {
	ID        string
	UserID    int
	ExpiresAt time.Time
}

type SessionStore interface {
	CreateNewSession(ctx context.Context, session Session) Session
	GetSession(ctx context.Context, sessionId string) (Session, error)
	DeleteSession(ctx context.Context, sessionId string) error
}

var (
	ErrSessionNotFound = errors.New("session not found")
	ErrSessionExpired  = errors.New("session expired")
)

type SessionStoreDB struct {
	db *pgxpool.Pool
}

func NewSessionStoreDB(pool *pgxpool.Pool) *SessionStoreDB {
	return &SessionStoreDB{
		db: pool,
	}
}

func (ss *SessionStoreDB) CreateNewSession(ctx context.Context, s Session) Session {
	_, err := ss.db.Exec(ctx, "INSERT INTO sessions (session_id, user_id, expires) VALUES ($1, $2, $3)", s.ID, s.UserID, s.ExpiresAt)
	if err != nil {
		fmt.Println(err)
		return Session{}
	}
	return s
}

func (ss *SessionStoreDB) GetSession(ctx context.Context, sessionId string) (Session, error) {
	s := Session{}

	err := ss.db.QueryRow(ctx, "SELECT * FROM sessions WHERE session_id = $1", sessionId).Scan(&s.ID, &s.UserID, &s.ExpiresAt)
	if err != nil {
		return Session{}, ErrSessionNotFound
	}
	if time.Now().After(s.ExpiresAt) {
		return Session{}, ErrSessionExpired
	}
	return s, nil
}

func (ss *SessionStoreDB) DeleteSession(ctx context.Context, sessionId string) error {
	if sessionId == "" {
		return errors.New("session ID cannot be empty")
	}
	_, err := ss.db.Exec(ctx, "DELETE FROM sessions WHERE session_id = $1", sessionId)
	if err != nil {
		return fmt.Errorf("failed to delete session: %v", err)
	}
	return nil
}
