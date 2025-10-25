package user

import (
	"fmt"

	"github.com/google/uuid"
)

type Store interface {
	GetByName(string) (*User, error)
	GetByID(uuid.UUID) (*User, error)
	GetByEmail(string) (*User, error)
	Add(*User) error
}

type InMemStore struct {
	usersByName  map[string]*User
	usersByID    map[uuid.UUID]*User
	usersByEmail map[string]*User
}

func (r InMemStore) Add(u *User) error {
	r.usersByID[u.ID] = u
	r.usersByName[u.Name] = u
	r.usersByEmail[u.Email] = u
	return nil
}

func NewInMemStore() *InMemStore {
	r := InMemStore{
		usersByName:  make(map[string]*User),
		usersByID:    make(map[uuid.UUID]*User),
		usersByEmail: make(map[string]*User),
	}

	initialUsers := []struct {
		name, email, password string
	}{
		{"admin", "admin@example.com", "password"},
		{"testuser", "test@example.com", "anotherPassword"},
	}

	for _, userData := range initialUsers {
		u, err := NewUser(userData.name, userData.email, userData.password)
		if err != nil {
			panic(fmt.Sprintf("failed to create user %s: %v", userData.email, err))
		}
		if err := r.Add(u); err != nil {
			panic(fmt.Sprintf("failed to add user %s: %v", userData.email, err))
		}
	}

	return &r
}

func (r InMemStore) GetByName(name string) (*User, error) {
	user, ok := r.usersByName[name]
	if !ok {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

func (r InMemStore) GetByID(id uuid.UUID) (*User, error) {
	user, ok := r.usersByID[id]
	if !ok {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

func (r InMemStore) GetByEmail(email string) (*User, error) {
	user, ok := r.usersByEmail[email]

	if !ok {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}
