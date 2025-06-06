package users

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"os"
)

type UserStoreDB struct {
	db *pgx.Conn
}

func NewUserStoreDB() *UserStoreDB {
	conn, err := pgx.Connect(context.Background(), "postgres://postgres:password@localhost:5432/postgres")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return &UserStoreDB{db: conn}
}

func (u *UserStoreDB) Close() {
	_ = u.db.Close(context.Background())
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

func (u *UserStoreDB) GetAll() []User {
	//TODO implement me
	panic("implement me")
}

func (u *UserStoreDB) FindByUsername(username string) (*User, bool) {
	//TODO implement me
	panic("implement me")
}

func (u *UserStoreDB) GetById(id int) (User, error) {
	//TODO implement me
	panic("implement me")
}
