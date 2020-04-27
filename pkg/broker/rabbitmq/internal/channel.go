/*
 * Copyright (c) Continental Corporation - All Rights Reserved
 *
 * This file is a part of Entry project.
 * ITS France - Entry squad members
 *
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

package internal

import (
	"errors"

	uuid "github.com/satori/go.uuid"
	"github.com/streadway/amqp"
)

// Channel wraps a RabbitMQ channel
type Channel struct {
	uuid       string
	connection *amqp.Connection
	channel    *amqp.Channel
}

// NewChannel initializes a RabbitMQ Channel from the given connection
func NewChannel(conn *amqp.Connection) (*Channel, error) {
	rabbitCh := &Channel{
		uuid:       uuid.NewV4().String(),
		connection: conn,
	}
	if err := rabbitCh.Connect(); err != nil {
		return nil, err
	}
	return rabbitCh, nil
}

// Connect to the broker
func (r *Channel) Connect() (err error) {
	r.channel, err = r.connection.Channel()
	return err
}

// Close current channel
func (r *Channel) Close() error {
	if r.channel == nil {
		return errors.New("Channel is nil")
	}
	return r.channel.Close()
}

// Publish a message
func (r *Channel) Publish(exchange, key string, message amqp.Publishing) error {
	if r.channel == nil {
		return errors.New("Channel is nil")
	}
	return r.channel.Publish(exchange, key, false, false, message)
}

// DeclareExchange is used to declare a new exchange to the broker
func (r *Channel) DeclareExchange(exchange string) error {
	return r.channel.ExchangeDeclare(
		exchange, // name
		"topic",  // kind
		false,    // durable
		false,    // autoDelete
		false,    // internal
		false,    // noWait
		nil,      // args
	)
}

// DeclareQueue is used to declare a new temporary queue to the broker
func (r *Channel) DeclareQueue(queue string) error {
	_, err := r.channel.QueueDeclare(
		queue, // name
		false, // durable
		true,  // autoDelete
		false, // exclusive
		false, // noWait
		nil,   // args
	)
	return err
}

// DeclareDurableQueue is used to declare a new durable queue to the broker
func (r *Channel) DeclareDurableQueue(queue string) error {
	_, err := r.channel.QueueDeclare(
		queue, // name
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,   // args
	)
	return err
}

// DeclareReplyQueue is used to declare a temporary queue for REQ/REP pattern
func (r *Channel) DeclareReplyQueue(queue string) error {
	_, err := r.channel.QueueDeclare(
		queue, // name
		false, // durable
		true,  // autoDelete
		true,  // exclusive
		false, // noWait
		nil,   // args
	)
	return err
}

// ConsumeQueue is used to declare a consumer to the broker
func (r *Channel) ConsumeQueue(queue string, autoAck bool) (<-chan amqp.Delivery, error) {
	return r.channel.Consume(
		queue,   // queue
		r.uuid,  // consumer
		autoAck, // autoAck
		false,   // exclusive
		false,   // nolocal
		false,   // nowait
		nil,     // args
	)
}

// BindQueue is used to connect a routing key from the exchange to the given queue
func (r *Channel) BindQueue(queue, key, exchange string, args amqp.Table) error {
	return r.channel.QueueBind(
		queue,    // name
		key,      // key
		exchange, // exchange
		false,    // noWait
		args,     // args
	)
}
