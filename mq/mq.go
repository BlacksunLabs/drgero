package mq

import (
	"fmt"

	"github.com/streadway/amqp"
)

// Client holds a RabbitMQ connection
type Client struct {
	conn *amqp.Connection
}

// Connect connects to RabbitMQ and stores an open connection in Client
func (c *Client) Connect(connectString string) error {
	if connectString == "" {
		return fmt.Errorf("cannot initialize a connection without a connection string")
	}

	var err error
	c.conn, err = amqp.Dial(fmt.Sprintf("%s/", connectString))
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	return nil
}

// PublishJSONToFanoutExchange publishes a JSON message to a fanout exchange
func (c *Client) PublishJSONToFanoutExchange(body []byte, exchangeName string) error {
	if c.conn == nil {
		return fmt.Errorf("cannot publish message without connecting to RabbitMQ")
	}

	ch, err := c.conn.Channel()
	defer ch.Close()

	err = ch.ExchangeDeclare(
		exchangeName,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare events exchange")
	}

	err = ch.Publish(
		exchangeName,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})

	if err != nil {
		return fmt.Errorf("failed to publish message to %s : %v", exchangeName, err)
	}

	return nil
}

// NewTempQueue creates a new temporary queue and returns its name
func (c *Client) NewTempQueue() (queueName string, err error) {
	if c.conn == nil {
		return "", fmt.Errorf("cannot create a temporary queue without connecting to RabbitMQ")
	}

	ch, err := c.conn.Channel()
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("failed to declare temporary queue: %v", err)
	}

	return q.Name, nil
}
