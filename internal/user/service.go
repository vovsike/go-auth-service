package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
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

type jwtCustomClaims struct {
	Roles []string `json:"roles"`
	jwt.RegisteredClaims
}

func (us *InMemoryService) Authenticate(identifier, password string) (string, error) {
	var u *User
	var err error
	if isEmail(identifier) {
		u, err = us.users.GetByEmail(identifier)
	} else {
		u, err = us.users.GetByName(identifier)
	}
	if err != nil {
		return "", fmt.Errorf("can't authenticate user: %v", err)
	}
	if !u.CheckPassword(password) {
		return "", fmt.Errorf("invalid password")
	}

	// Try to fetch roles, but don't fail if roles service is unavailable
	roleNames := make([]string, 0)
	resp, err := http.Get("http://localhost:4001/users/" + u.ID.String() + "/roles")
	if err == nil {
		defer resp.Body.Close()

		var roles []struct {
			ID   uuid.UUID `json:"id"`
			Name string    `json:"name"`
		}

		err = json.NewDecoder(resp.Body).Decode(&roles)
		if err == nil {
			roleNames = make([]string, 0, len(roles))
			for _, role := range roles {
				roleNames = append(roleNames, role.Name)
			}
		}
	}

	t, err := issueSignedToken(u, roleNames)
	if err != nil {
		return "", err
	}
	return t, nil
}

func issueSignedToken(user *User, roles []string) (string, error) {
	secret, ok := os.LookupEnv("SIGN_KEY")
	if !ok {
		return "", fmt.Errorf("SIGN_KEY not set")
	}
	if user == nil {
		return "", fmt.Errorf("user is nil")
	}
	claims := jwtCustomClaims{
		Roles: roles,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			Issuer:    "auth-service",
			Subject:   user.ID.String(),
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(secret))
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

func isEmail(identifier string) bool {
	return identifier != "" && strings.Contains(identifier, "@")
}
