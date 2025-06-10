package users

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	Store Store
}

func NewUserService(store Store) *UserService {
	return &UserService{Store: store}
}

func (s *UserService) CheckUserPassword(ctx context.Context, un string, passwordToCheck string) error {
	user, err := s.Store.FindByUsername(ctx, un)
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(passwordToCheck))
	if err != nil {
		return errors.New("password is incorrect")
	}
	return nil
}

func (s *UserService) CreateNewUser(ctx context.Context, username string, password string) (User, error) {
	u, err := s.Store.Add(ctx, username, password)
	if err != nil {
		return User{}, fmt.Errorf("failed to create a new user: %w", err)
	}
	return u, nil
}
