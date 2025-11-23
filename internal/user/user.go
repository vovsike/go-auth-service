package user

import (
	"fmt"
	"strings"
	"time"

	"net/mail"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        uuid.UUID
	Name      string
	Email     string
	hash      []byte
	Joined    time.Time
	Activated bool
}

func NewUser(name, email, password string) (*User, error) {
	if !isValidEmail(email) {
		return nil, fmt.Errorf("invalid email")
	}
	if !isPasswordValid(password) {
		return nil, fmt.Errorf("invalid password")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return &User{
		ID:        uuid.New(),
		Name:      name,
		Email:     email,
		hash:      hash,
		Joined:    time.Now(),
		Activated: false,
	}, nil
}

func isValidEmail(email string) bool {
	// Basic checks
	if len(email) < 3 || len(email) > 254 {
		return false
	}

	// Use net/mail for RFC 5322 validation
	addr, err := mail.ParseAddress(email)
	if err != nil {
		return false
	}

	// Ensure it's just the email (no display name)
	if addr.Address != email {
		return false
	}

	// Check for @ and domain
	parts := strings.Split(email, "@")
	if len(parts) != 2 || len(parts[1]) == 0 {
		return false
	}

	return true
}

func isPasswordValid(password string) bool {
	if len(password) < 8 {
		return false
	}
	if strings.Contains(password, " ") {
		return false
	}
	return true
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword(u.hash, []byte(password))
	if err != nil {
		return false
	}
	return true
}

func (u *User) Activate() {
	u.Activated = true
}

type CreationDTO struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type DTO struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Activated bool      `json:"activated"`
	Joined    time.Time `json:"joined"`
}

type PasswordWrapper struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

type TokenWrapper struct {
	Token string `json:"token"`
}
