package accounts

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"restapi/internal/mq"
)

type ServiceDB struct {
	dbpool   *pgxpool.Pool
	mqClient *mq.Client
}

func NewServiceDB(pool *pgxpool.Pool, mqClient *mq.Client) *ServiceDB {
	return &ServiceDB{dbpool: pool, mqClient: mqClient}
}

func (s *ServiceDB) CreateNewAccount(ctx context.Context, username string, password string) (Account, error) {
	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return Account{}, fmt.Errorf("failed to hash password: %w", err)
	}
	acc := Account{
		ID:       uuid.New(),
		Username: username,
		Password: string(passHash),
	}
	_, err = s.dbpool.Exec(ctx, "INSERT INTO accounts (id, username, password) VALUES ($1, $2, $3)", acc.ID, acc.Username, acc.Password)
	if err != nil {
		return Account{}, err
	}

	type User struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}

	user := User{
		Username: acc.Username,
		Email:    acc.Username + "@test.test",
	}

	err = s.mqClient.Publish("accounts", "", true, false, user)

	if err != nil {
		return Account{}, err
	}

	return acc, nil
}
