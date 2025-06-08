package users

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"os"
)

type UserStoreDB struct {
	db *pgxpool.Pool
}

func NewUserStoreDB(pool *pgxpool.Pool) *UserStoreDB {
	return &UserStoreDB{db: pool}
}

func (u *UserStoreDB) Close() {
	_ = u.db.Close
}

func (u *UserStoreDB) Ping() {
	var testString string
	err := u.db.QueryRow(context.Background(), "SELECT 'test'").Scan(&testString)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(testString)
}

func (u *UserStoreDB) FindByUsername(username string) (*User, bool) {
	user := User{}
	err := u.db.QueryRow(context.Background(), "SELECT * FROM users WHERE name = $1", username).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return nil, false
	}
	return &user, true
}

func (u *UserStoreDB) GetById(id int) (User, error) {
	var user User

	err := u.db.QueryRow(context.Background(),
		"SELECT user_id, name, password FROM users WHERE user_id = $1", id).
		Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, fmt.Errorf("user with id %d not found", id)
		}
		return User{}, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

func (u *UserStoreDB) Add(username string, password string) (User, error) {
	if username == "" || password == "" {
		return User{}, errors.New("username or password is empty")
	}

	exists, _ := u.FindByUsername(username)
	if exists != nil {
		return User{}, fmt.Errorf("user with username %s already exists", username)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, fmt.Errorf("failed to hash password: %w", err)
	}

	var userID int

	err = u.db.QueryRow(context.Background(),
		"INSERT INTO users (user_id ,name, password) VALUES ($1, $2, $3) RETURNING user_id",
		2, username, hashedPassword).Scan(&userID)

	if err != nil {
		return User{}, fmt.Errorf("failed to add user: %w", err)
	}

	return User{
		ID:       userID,
		Username: username,
	}, nil
}
