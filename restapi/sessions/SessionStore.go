package sessions

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"time"
)

type Session struct {
	SessionId string
	UserId    int
	Expires   time.Time
}

type SessionStore interface {
	CreateNewSession(session Session) Session
	GetSession(sessionId string) (Session, error)
	DeleteSession(sessionId string)
}

type SessionStoreDB struct {
	db *pgx.Conn
}

func NewSessionStoreDB(conn *pgx.Conn) *SessionStoreDB {
	return &SessionStoreDB{
		db: conn,
	}
}

func (ss *SessionStoreDB) CreateNewSession(s Session) Session {
	_, err := ss.db.Exec(context.Background(), "INSERT INTO sessions (session_id, user_id, expires) VALUES ($1, $2, $3)", s.SessionId, s.UserId, s.Expires)
	if err != nil {
		fmt.Println(err)
		return Session{}
	}
	return s
}

func (ss *SessionStoreDB) GetSession(sessionId string) (Session, error) {
	s := Session{}

	err := ss.db.QueryRow(context.Background(), "SELECT * FROM sessions WHERE session_id = $1", sessionId).Scan(&s.SessionId, &s.UserId, &s.Expires)
	if err != nil {
		return Session{}, err
	}
	return s, nil
}

func (ss *SessionStoreDB) DeleteSession(sessionId string) {
	//TODO implement me
	panic("implement me")
}
