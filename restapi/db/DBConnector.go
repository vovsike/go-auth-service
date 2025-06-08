package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
)

func CreateConnectionPool() (*pgxpool.Pool, error) {
	connString := os.Getenv("DATABASE_URL")
	if connString == "" {
		connString = "postgres://postgres:password@localhost:5432/postgres"
	}
	conn, err := pgxpool.New(context.Background(), "postgres://postgres:password@localhost:5432/postgres")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}
	return conn, nil
}
