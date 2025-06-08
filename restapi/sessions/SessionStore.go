package sessions

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"time"
)

type Session struct {
	sessionId string
	userId    int
	expires   time.Time
}

type SessionStore interface {
	CreateNewSession(session Session) Session
	GetSession(sessionId string) Session
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
	_, err := ss.db.Exec(context.Background(), "INSERT INTO sessions (session_id, user_id, expires) VALUES ($1, $2, $3)", s.sessionId, s.userId, s.expires)
	if err != nil {
		fmt.Println(err)
		return Session{}
	}
	return s
}

func (ss *SessionStoreDB) GetSession(sessionId string) Session {
	//TODO implement me
	panic("implement me")
}

func (ss *SessionStoreDB) DeleteSession(sessionId string) {
	//TODO implement me
	panic("implement me")
}
