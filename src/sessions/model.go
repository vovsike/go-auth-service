package sessions

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type Session struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	ExpiresAt time.Time
}

type SessionService interface {
	Authenticate(ctx context.Context, userId uuid.UUID) (Session, error)
}
