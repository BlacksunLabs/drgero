package mq

// Copyright 2018 Blacksun Research Labs

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// 	http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

// BindQueueToExchange binds a queue to an exchange
func (c *Client) BindQueueToExchange(queueName string, exchangeName string) error {
	if c.conn == nil {
		return fmt.Errorf("cannot bind queue to exchange before connecting to RabbitMQ")
	}

	ch, err := c.conn.Channel()
	defer ch.Close()

	err = ch.QueueBind(
		queueName,
		"",
		exchangeName,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to bind %s to exchange %s : %v", queueName, exchangeName, err)
	}

	return nil
}

// GetChannel gets a channel from RabbitMQ
func (c *Client) GetChannel() (*amqp.Channel, error) {
	if c.conn == nil {
		return nil, fmt.Errorf("cannot get channel before connecting to RabbitMQ")
	}

	ch, err := c.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %v", err)
	}

	return ch, nil
}
