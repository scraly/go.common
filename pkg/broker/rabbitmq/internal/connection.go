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
	"crypto/tls"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/scraly/go.common/pkg/tlsconfig"

	"github.com/streadway/amqp"
)

var (
	// DefaultExchange name
	DefaultExchange = "toto"
	// DefaultRabbitURL is the default RabbitMQ URI
	DefaultRabbitURL = "amqp://guest:guest@127.0.0.1:5672"

	dial    = amqp.Dial
	dialTLS = amqp.DialTLS
)

// Connection wraps a RabbitMQ connection
type Connection struct {
	Exchange string

	connection      *amqp.Connection
	channel         *Channel
	exchangeChannel *Channel
	url             string

	sync.Mutex
	connected bool
	close     chan bool
}

// NewConnection initializes a new connection wrapper
func NewConnection(exchange string, urls []string) *Connection {
	var url string

	if len(urls) > 0 && regexp.MustCompile("^amqp(s)?://.*").MatchString(urls[0]) {
		url = urls[0]
	} else {
		url = DefaultRabbitURL
	}

	if len(exchange) == 0 {
		exchange = DefaultExchange
	}

	return &Connection{
		Exchange: exchange,
		url:      url,
		close:    make(chan bool),
	}
}

// -----------------------------------------------------------------------------

// Connect to the broker
func (r *Connection) Connect(secure bool, config *tls.Config) error {
	r.Lock()

	// already connected
	if r.connected {
		r.Unlock()
		return nil
	}

	// check it was closed
	select {
	case <-r.close:
		r.close = make(chan bool)
	default:
		// no op
		// new conn
	}

	r.Unlock()

	return r.connect(secure, config)
}

// Close broker connection
func (r *Connection) Close() error {
	r.Lock()
	defer r.Unlock()

	select {
	case <-r.close:
		return nil
	default:
		close(r.close)
		r.connected = false
	}

	return r.connection.Close()
}

// Consume declares a new consumer
func (r *Connection) Consume(queue, key string, headers amqp.Table, autoAck, durableQueue bool) (*Channel, <-chan amqp.Delivery, error) {
	consumerChannel, err := NewChannel(r.connection)
	if err != nil {
		return nil, nil, err
	}

	if durableQueue {
		err = consumerChannel.DeclareDurableQueue(queue)
	} else {
		err = consumerChannel.DeclareQueue(queue)
	}

	if err != nil {
		return nil, nil, err
	}

	deliveries, err := consumerChannel.ConsumeQueue(queue, autoAck)
	if err != nil {
		return nil, nil, err
	}

	err = consumerChannel.BindQueue(queue, key, r.Exchange, headers)
	if err != nil {
		return nil, nil, err
	}

	return consumerChannel, deliveries, nil
}

// Publish a message from the exchange with given routing key
func (r *Connection) Publish(exchange, key string, msg amqp.Publishing) error {
	return r.exchangeChannel.Publish(exchange, key, msg)
}

// -----------------------------------------------------------------------------

func (r *Connection) connect(secure bool, config *tls.Config) error {
	// try connect
	if err := r.tryConnect(secure, config); err != nil {
		return err
	}

	// connected
	r.Lock()
	r.connected = true
	r.Unlock()

	// create reconnect loop
	go r.reconnect(secure, config)
	return nil
}

func (r *Connection) reconnect(secure bool, config *tls.Config) {
	// skip first connect
	var connect bool

	for {
		if connect {
			// try reconnect
			if err := r.tryConnect(secure, config); err != nil {
				time.Sleep(1 * time.Second)
				continue
			}

			// connected
			r.Lock()
			r.connected = true
			r.Unlock()
		}

		connect = true
		notifyClose := make(chan *amqp.Error)
		r.connection.NotifyClose(notifyClose)

		// block until closed
		select {
		case <-notifyClose:
			r.Lock()
			r.connected = false
			r.Unlock()
		case <-r.close:
			return
		}
	}
}

func (r *Connection) tryConnect(secure bool, config *tls.Config) error {
	var err error

	if secure || config != nil || strings.HasPrefix(r.url, "amqps://") {
		if config == nil {
			config = tlsconfig.ClientDefault()
		}

		url := strings.Replace(r.url, "amqp://", "amqps://", 1)
		r.connection, err = dialTLS(url, config)
	} else {
		r.connection, err = dial(r.url)
	}

	if err != nil {
		return err
	}

	if r.channel, err = NewChannel(r.connection); err != nil {
		return err
	}

	if err = r.channel.DeclareExchange(r.Exchange); err != nil {
		return err
	}

	r.exchangeChannel, err = NewChannel(r.connection)
	return err
}
