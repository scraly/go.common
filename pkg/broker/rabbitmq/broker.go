/*
 * Copyright (c) Continental Corporation - All Rights Reserved
 *
 * This file is a part of Entry project.
 * ITS France - Entry squad members
 *
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

package rabbitmq

import (
	"context"
	"strings"

	"github.com/scraly/go.common/pkg/broker"
	"github.com/scraly/go.common/pkg/broker/rabbitmq/internal"
	"github.com/scraly/go.common/pkg/log"
	"github.com/scraly/go.common/pkg/storage/codec/msgpack"

	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

type rbroker struct {
	conn  *internal.Connection
	addrs []string
	opts  broker.Options
}

// NewBroker initializes a RabbitMQ broker instance
func NewBroker(opts ...broker.Option) broker.Broker {
	options := broker.Options{
		Context: context.Background(),
		// Default codec
		Codec: msgpack.NewCodec(),
	}

	for _, o := range opts {
		o(&options)
	}

	var exchange string
	if e, ok := options.Context.Value(exchangeKey{}).(string); ok {
		exchange = e
	}

	return &rbroker{
		conn:  internal.NewConnection(exchange, options.Addrs),
		addrs: options.Addrs,
		opts:  options,
	}
}

func init() {
	broker.Register("rabbitmq", NewBroker)
}

// -----------------------------------------------------------------------------

func (r *rbroker) Publish(topic string, msg *broker.Message, _ ...broker.PublishOption) error {
	// Normalize topic name
	topic = strings.Replace(topic, ":", ".", -1)

	// Prepare message
	m := amqp.Publishing{
		Body:    msg.Body,
		Headers: amqp.Table{},
	}

	for k, v := range msg.Header {
		m.Headers[k] = v
	}

	return r.conn.Publish(r.conn.Exchange, topic, m)
}

func (r *rbroker) Subscribe(topic string, handler broker.Handler, opts ...broker.SubscribeOption) (broker.Subscriber, error) {
	opt := broker.SubscribeOptions{
		AutoAck: true,
	}

	for _, o := range opts {
		o(&opt)
	}

	durableQueue := false
	if opt.Context != nil {
		durableQueue, _ = opt.Context.Value(durableQueueKey{}).(bool)
	}

	var headers map[string]interface{}
	if opt.Context != nil {
		if h, ok := opt.Context.Value(headersKey{}).(map[string]interface{}); ok {
			headers = h
		}
	}

	ch, sub, err := r.conn.Consume(
		opt.Queue,
		topic,
		headers,
		opt.AutoAck,
		durableQueue,
	)
	if err != nil {
		return nil, err
	}

	fn := func(msg amqp.Delivery) {
		header := make(map[string]string)
		for k, v := range msg.Headers {
			header[k], _ = v.(string)
		}
		m := &broker.Message{
			Header: header,
			Body:   msg.Body,
		}
		if err := handler(&publication{d: msg, m: m, t: msg.RoutingKey}); err != nil {
			log.Bg().Error("Unable to register subscription handler", zap.Error(err), zap.String("topic", topic))
		}
	}

	go func() {
		for d := range sub {
			go fn(d)
		}
	}()

	return &subscriber{ch: ch, topic: topic, opts: opt}, nil
}

func (r *rbroker) Options() broker.Options {
	return r.opts
}

func (r *rbroker) String() string {
	return "rabbitmq"
}

func (r *rbroker) Address() string {
	if len(r.addrs) > 0 {
		return r.addrs[0]
	}
	return ""
}

func (r *rbroker) Init(opts ...broker.Option) error {
	for _, o := range opts {
		o(&r.opts)
	}
	return nil
}

func (r *rbroker) Connect() error {
	return r.conn.Connect(r.opts.Secure, r.opts.TLSConfig)
}

func (r *rbroker) Disconnect() error {
	return r.conn.Close()
}
