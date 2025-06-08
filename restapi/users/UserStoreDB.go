package users

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
	"os"
)

type UserStoreDB struct {
	db *pgx.Conn
}

func NewUserStoreDB(conn *pgx.Conn) *UserStoreDB {
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
	user := User{}
	err := u.db.QueryRow(context.Background(), "SELECT * FROM users WHERE name = $1", username).Scan(&user.Id, &user.Username, &user.Password)
	if err != nil {
		return nil, false
	}
	return &user, true
}

func (u *UserStoreDB) GetById(id int) (User, error) {
	//TODO implement me
	panic("implement me")
}

func (u *UserStoreDB) Add(username string, password string) User {
	phash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	_, err := u.db.Exec(context.Background(), "INSERT INTO users (user_id ,name, password) VALUES ($1, $2, $3) RETURNING user_id", 2, username, phash)
	if err != nil {
		fmt.Println(err)
		return User{}
	}
	return User{
		Id:       2,
		Username: username,
	}
}
