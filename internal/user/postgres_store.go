package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStore struct {
	pool *pgxpool.Pool
}

func NewPostgresStore(pool *pgxpool.Pool) *PostgresStore {
	return &PostgresStore{
		pool: pool,
	}
}

func (s *PostgresStore) Add(u *User) error {
	query := `
		INSERT INTO users (id, name, email, password_hash, joined, activated)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := s.pool.Exec(
		context.Background(),
		query,
		u.ID,
		u.Name,
		u.Email,
		u.hash,
		u.Joined,
		u.Activated,
	)
	if err != nil {
		return fmt.Errorf("failed to add user: %w", err)
	}

	return nil
}

func (s *PostgresStore) GetByName(name string) (*User, error) {
	query := `
		SELECT id, name, email, password_hash, joined, activated
		FROM users
		WHERE name = $1
	`

	var u User
	err := s.pool.QueryRow(context.Background(), query, name).Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&u.hash,
		&u.Joined,
		&u.Activated,
	)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	return &u, nil
}

func (s *PostgresStore) GetByID(id uuid.UUID) (*User, error) {
	query := `
		SELECT id, name, email, password_hash, joined, activated
		FROM users
		WHERE id = $1
	`

	var u User
	err := s.pool.QueryRow(context.Background(), query, id).Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&u.hash,
		&u.Joined,
		&u.Activated,
	)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	return &u, nil
}

func (s *PostgresStore) GetByEmail(email string) (*User, error) {
	query := `
		SELECT id, name, email, password_hash, joined, activated
		FROM users
		WHERE email = $1
	`

	var u User
	err := s.pool.QueryRow(context.Background(), query, email).Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&u.hash,
		&u.Joined,
		&u.Activated,
	)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	return &u, nil
}
