package users

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	Store UserStore
}

func NewUserService(store UserStore) *UserService {
	return &UserService{Store: store}
}

func (s *UserService) CheckUserPassword(un string, passwordToCheck string) (bool, error) {
	user, found := s.Store.FindByUsername(un)
	if !found {
		return false, errors.New("user not found")
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(passwordToCheck))
	if err != nil {
		return false, errors.New("password is incorrect")
	}
	return true, nil
}

func (s *UserService) CreateNewUser(username string, password string) (User, error) {
	u, err := s.Store.Add(username, password)
	if err != nil {
		return User{}, fmt.Errorf("failed to create a new user: %w", err)
	}
	return u, nil
}
