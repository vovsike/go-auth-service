package accounts

import (
	"context"
	"github.com/google/uuid"
)

type Account struct {
	ID       uuid.UUID
	UserID   uuid.UUID
	Username string
	Password string
}

type AccountService interface {
	CreateNewAccount(ctx context.Context, username string, password string) (Account, error)
}
