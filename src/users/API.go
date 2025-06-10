package users

import "context"

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Store interface {
	FindByUsername(ctx context.Context, username string) (User, error)
	GetById(ctx context.Context, id int) (User, error)
	Add(ctx context.Context, username string, password string) (User, error)
}

type Service interface {
	CheckUserPassword(ctx context.Context, un string, passwordToCheck string) error
	CreateNewUser(ctx context.Context, username string, password string) (User, error)
}
