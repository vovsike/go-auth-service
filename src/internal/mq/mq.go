package mq

import (
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
)

// Client represents a RabbitMQ client
type Client struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

// NewClient creates a new RabbitMQ client
func NewClient(url string) (*Client, error) {
	conn, ch, err := connect(url)
	if err != nil {
		return nil, err
	}

	// Enable publisher confirms
	err = ch.Confirm(false)
	if err != nil {
		return nil, fmt.Errorf("failed to enable publisher confirms: %w", err)
	}

	return &Client{
		conn:    conn,
		channel: ch,
	}, nil
}

// connect establishes a connection to RabbitMQ and creates a channel
func connect(url string) (*amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, nil, fmt.Errorf("failed to create channel: %w", err)
	}

	return conn, ch, nil
}

// DeclareExchange declares an exchange
func (c *Client) DeclareExchange(name, kind string, durable, autoDelete, internal, noWait bool) error {
	err := c.channel.ExchangeDeclare(
		name,
		kind,
		durable,
		autoDelete,
		internal,
		noWait,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare exchange: %w", err)
	}
	return nil
}

// Publish publishes a message to the exchange with the given routing key
// It returns an error if the message could not be published
func (c *Client) Publish(exchange, routingKey string, mandatory, immediate bool, msg any) error {
	preparedMsg, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	err = c.channel.Publish(
		exchange,
		routingKey,
		mandatory,
		immediate,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        preparedMsg,
		},
	)

	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}

// Close closes the connection and channel
func (c *Client) Close() {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
}

// GetChannel returns the underlying channel
func (c *Client) GetChannel() *amqp.Channel {
	return c.channel
}
