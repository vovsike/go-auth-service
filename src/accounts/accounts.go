package accounts

import (
	"context"
	"encoding/json"
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
	ch := mqClient.GetChannel()
	q, err := ch.QueueDeclare("account_updates", false, false, true, false, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}
	err = ch.QueueBind(q.Name, "update", "accounts", false, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to bind queue: %w", err)
	}
	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to consume messages: %w", err)
	}
	s := &ServiceDB{
		dbpool:   pool,
		mqClient: mqClient,
		updateCh: msgs,
	}
	s.listenForUpdates(context.TODO())
	return &ServiceDB{dbpool: pool, mqClient: mqClient}, nil
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

func (s *ServiceDB) UpdateAccountUsername(ctx context.Context, oldUsername, newUsername string) (Account, error) {
	_, err := s.dbpool.Exec(ctx, "UPDATE accounts SET username = $1 WHERE username = $2", newUsername, oldUsername)
	if err != nil {
		return Account{}, err
	}

	return Account{}, nil
}

func (s *ServiceDB) listenForUpdates(ctx context.Context) {

	type updateUsername struct {
		OldUsername string `json:"oldUsername"`
		NewUsername string `json:"newUsername"`
	}

	var update updateUsername

	for msg := range s.updateCh {
		fmt.Printf("Received a message: %s\n", msg.Body)
		err := json.Unmarshal(msg.Body, &update)
		if err != nil {
			fmt.Println(err)
		}
		_, err = s.UpdateAccountUsername(ctx, update.OldUsername, update.NewUsername)
		if err != nil {
			fmt.Println(err)
		}
	}
}
