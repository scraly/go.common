/*
 * Copyright (c) Continental Corporation - All Rights Reserved
 *
 * This file is a part of Entry project.
 * ITS France - Entry squad members
 *
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

package nats

import (
	"github.com/scraly/go.common/pkg/broker"
	"github.com/scraly/go.common/pkg/log"
	"github.com/scraly/go.common/pkg/storage/codec/msgpack"

	nats "github.com/nats-io/go-nats"
	"go.uber.org/zap"
)

type nBroker struct {
	addrs []string
	conn  *nats.Conn
	opts  broker.Options
}

// NewBroker initializes a NATS broker instance
func NewBroker(opts ...broker.Option) broker.Broker {
	options := broker.Options{
		// Default codec
		Codec: msgpack.NewCodec(),
	}

	for _, o := range opts {
		o(&options)
	}

	return &nBroker{
		addrs: setAddrs(options.Addrs),
		opts:  options,
	}
}

func init() {
	broker.Register("nats", NewBroker)
}

// -----------------------------------------------------------------------------

func (n *nBroker) Address() string {
	if n.conn != nil && n.conn.IsConnected() {
		return n.conn.ConnectedUrl()
	}
	if len(n.addrs) > 0 {
		return n.addrs[0]
	}

	return ""
}

func (n *nBroker) Connect() error {
	if n.conn != nil {
		return nil
	}

	opts := nats.DefaultOptions
	opts.Servers = n.addrs
	opts.Secure = n.opts.Secure
	opts.TLSConfig = n.opts.TLSConfig

	// secure might not be set
	if n.opts.TLSConfig != nil {
		opts.Secure = true
	}

	c, err := opts.Connect()
	if err != nil {
		return err
	}
	n.conn = c
	return nil
}

func (n *nBroker) Disconnect() error {
	n.conn.Close()
	return nil
}

func (n *nBroker) Init(opts ...broker.Option) error {
	for _, o := range opts {
		o(&n.opts)
	}
	n.addrs = setAddrs(n.opts.Addrs)
	return nil
}

func (n *nBroker) Options() broker.Options {
	return n.opts
}

func (n *nBroker) Publish(topic string, msg *broker.Message, _ ...broker.PublishOption) error {
	b, err := n.opts.Codec.Marshal(msg)
	if err != nil {
		return err
	}
	return n.conn.Publish(topic, b)
}

func (n *nBroker) Subscribe(topic string, handler broker.Handler, opts ...broker.SubscribeOption) (broker.Subscriber, error) {
	opt := broker.SubscribeOptions{
		AutoAck: true,
	}

	for _, o := range opts {
		o(&opt)
	}

	fn := func(msg *nats.Msg) {
		var m broker.Message
		if err := n.opts.Codec.Unmarshal(msg.Data, &m); err != nil {
			log.Bg().Error("Unable to decode subscription message", zap.Error(err), zap.String("topic", topic))
			return
		}
		err := handler(&publication{m: &m, t: msg.Subject})
		if err != nil {
			log.Bg().Error("Unable to register subscription handler", zap.Error(err), zap.String("topic", topic))
		}
	}

	var sub *nats.Subscription
	var err error

	if len(opt.Queue) > 0 {
		sub, err = n.conn.QueueSubscribe(topic, opt.Queue, fn)
	} else {
		sub, err = n.conn.Subscribe(topic, fn)
	}
	if err != nil {
		return nil, err
	}
	return &subscriber{s: sub, opts: opt}, nil
}

func (n *nBroker) String() string {
	return "nats"
}
