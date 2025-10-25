package user

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Service interface {
	GetUserByName(name string) (*User, error)
	GetUserByID(id uuid.UUID) (*User, error)
	GetUserByEmail(email string) (*User, error)
	CreateNewUser(name, email, password string) (*User, error)
	Authenticate(email, password string) (string, error)
}

type InMemoryService struct {
	users Store
}

func NewInMemoryUserService(users Store) *InMemoryService {
	return &InMemoryService{
		users: users,
	}
}

func (us *InMemoryService) Authenticate(email, password string) (string, error) {
	u, err := us.users.GetByEmail(email)
	if err != nil {
		return "", fmt.Errorf("can,t authnenticate user: %v", err)
	}
	if !u.CheckPassword(password) {
		return "", fmt.Errorf("invalid password")
	}
	t, err := issueSignedToken(u)
	if err != nil {
		return "", err
	}
	return t, nil
}

func issueSignedToken(user *User) (string, error) {
	secret, ok := os.LookupEnv("SIGN_KEY")
	if !ok {
		return "", fmt.Errorf("SIGN_KEY not set")
	}
	if user == nil {
		return "", fmt.Errorf("user is nil")
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			Issuer:    "auth-service",
			Subject:   user.ID.String(),
		})

	s, err := t.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return s, nil
}

func (us *InMemoryService) GetUserByName(name string) (*User, error) {
	u, err := us.users.GetByName(name)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (us *InMemoryService) GetUserByID(id uuid.UUID) (*User, error) {
	u, err := us.users.GetByID(id)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (us *InMemoryService) GetUserByEmail(email string) (*User, error) {
	u, err := us.users.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (us *InMemoryService) CreateNewUser(name, email, password string) (*User, error) {
	_, err := us.users.GetByEmail(email)
	if err == nil {
		return nil, fmt.Errorf("user already exists")
	}
	user, err := NewUser(name, email, password)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %v", err)
	}
	err = us.users.Add(user)
	if err != nil {
		return nil, fmt.Errorf("failed to createa a new user: %v", err)
	}
	return user, nil
}
