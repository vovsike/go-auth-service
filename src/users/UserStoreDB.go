package users

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
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

func (u *UserStoreDB) FindByUsername(ctx context.Context, username string) (User, error) {
	user := User{}
	err := u.db.QueryRow(ctx, "SELECT * FROM users WHERE name = $1", username).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return User{}, fmt.Errorf("failed to find user: %w", err)
	}
	return user, nil
}

func (u *UserStoreDB) GetById(ctx context.Context, id int) (User, error) {
	var user User

	err := u.db.QueryRow(ctx,
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

func (u *UserStoreDB) Add(ctx context.Context, username string, password string) (User, error) {
	if username == "" || password == "" {
		return User{}, errors.New("username or password is empty")
	}

	_, err := u.FindByUsername(ctx, username)
	if err != nil {
		return User{}, fmt.Errorf("user with username %s already exists", username)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, fmt.Errorf("failed to hash password: %w", err)
	}

	var userID int

	err = u.db.QueryRow(ctx,
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
