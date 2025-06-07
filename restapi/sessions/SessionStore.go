package sessions

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"time"
)

type SessionStore interface {
	CreateNewSession(userId int)
	VerifySession(sessionId string)
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

func (s *SessionStoreDB) CreateNewSession(userId int) {
	sid, _ := uuid.NewRandom()
	_, err := s.db.Exec(context.Background(), "INSERT INTO sessions (session_id, user_id, expires) VALUES ($1, $2, $3)", sid, userId, time.Now())
	if err != nil {
		fmt.Println(err)
	}
}

func (s *SessionStoreDB) VerifySession(sessionId string) {
	//TODO implement me
	panic("implement me")
}

func (s *SessionStoreDB) DeleteSession(sessionId string) {
	//TODO implement me
	panic("implement me")
}
