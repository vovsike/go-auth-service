package accounts

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	amqp "github.com/rabbitmq/amqp091-go"
	"golang.org/x/crypto/bcrypt"
	"restapi/internal/mq"
)

type ServiceDB struct {
	dbpool   *pgxpool.Pool
	mqClient *mq.Client
	updateCh <-chan amqp.Delivery
}

func NewServiceDB(pool *pgxpool.Pool, mqClient *mq.Client) (*ServiceDB, error) {
	s := &ServiceDB{
		dbpool:   pool,
		mqClient: mqClient,
	}
	return s, nil
}

func (s *ServiceDB) CreateNewAccount(ctx context.Context, username string, password string) (Account, error) {
	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return Account{}, fmt.Errorf("failed to hash password: %w", err)
	}
	acc := Account{
		ID:       uuid.New(),
		UserID:   uuid.New(),
		Username: username,
		Password: string(passHash),
	}
	_, err = s.dbpool.Exec(
		ctx,
		"INSERT INTO accounts (accountid, userid, username, password) VALUES ($1, $2, $3, $4)",
		acc.ID,
		acc.UserID,
		acc.Username,
		acc.Password,
	)
	if err != nil {
		return Account{}, err
	}

	type User struct {
		UserID   uuid.UUID `json:"userid"`
		Username string    `json:"username"`
	}

	user := User{
		UserID:   acc.UserID,
		Username: acc.Username,
	}

	err = s.mqClient.Publish("accounts", "create", false, false, user)

	if err != nil {
		return Account{}, err
	}

	return acc, nil
}
