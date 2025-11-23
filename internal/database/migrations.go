package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// RunMigrations executes all database migrations on startup
func RunMigrations(ctx context.Context, pool *pgxpool.Pool) error {
	// Create users table
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY,
		name VARCHAR(255) NOT NULL UNIQUE,
		email VARCHAR(255) NOT NULL UNIQUE,
		password_hash BYTEA NOT NULL,
		joined TIMESTAMP NOT NULL DEFAULT NOW(),
		activated BOOLEAN NOT NULL DEFAULT FALSE
	);

	CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
	CREATE INDEX IF NOT EXISTS idx_users_name ON users(name);
	`

	_, err := pool.Exec(ctx, createUsersTable)
	if err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	return nil
}
