package users

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	Store UserStore
}

func NewUserService(store UserStore) *UserService {
	return &UserService{Store: store}
}

func (s *UserService) GetAllUsers() []User {
	return s.Store.GetAll()
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

func (s *UserService) AddUser(username string, password string) (User, error) {
	u := s.Store.Add(username, password)
	if u == (User{}) {
		return User{}, errors.New("user could not be added")
	}
	return u, nil
}
