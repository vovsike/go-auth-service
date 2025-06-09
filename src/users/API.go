package users

import "context"

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserStore interface {
	FindByUsername(ctx context.Context, username string) (*User, bool)
	GetById(ctx context.Context, id int) (User, error)
	Add(ctx context.Context, username string, password string) (User, error)
}
