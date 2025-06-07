package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"os"
)

func CreateNewConnection() *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), "postgres://postgres:password@localhost:5432/postgres")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return conn
}
