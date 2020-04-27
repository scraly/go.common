/*
 * Copyright (c) Continental Corporation - All Rights Reserved
 *
 * This file is a part of Entry project.
 * ITS France - Entry squad members
 *
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

package stan

import (
	"github.com/scraly/go.common/pkg/broker"
	"github.com/scraly/go.common/pkg/log"
	"github.com/scraly/go.common/pkg/storage/codec/msgpack"
	"github.com/scraly/go.common/pkg/util/runtime"

	nats "github.com/nats-io/go-nats"
	stan "github.com/nats-io/go-nats-streaming"
	"go.uber.org/zap"
)

type nBroker struct {
	addrs []string
	conn  stan.Conn
	opts  broker.Options
}

// NewBroker initializes a NATS streaming aka STAN broker instance
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
	broker.Register("stan", NewBroker)
}

// -----------------------------------------------------------------------------

func (n *nBroker) Address() string {
	if n.conn != nil && n.conn.NatsConn() != nil && n.conn.NatsConn().IsConnected() {
		return n.conn.NatsConn().ConnectedUrl()
	}
	if len(n.addrs) > 0 {
		return n.addrs[0]
	}

	return ""
}

func (n *nBroker) Connect() error {
	if n.conn != nil && n.conn.NatsConn().IsConnected() {
		return nil
	}

	conf := n.opts.Context.Value(broker.ContextKey(FullConfName)).(Configuration)

	opts := nats.DefaultOptions
	opts.Servers = n.addrs
	opts.Secure = n.opts.Secure
	opts.TLSConfig = n.opts.TLSConfig
	opts.User = conf.User
	opts.Password = conf.Password

	// secure might not be set
	if n.opts.TLSConfig != nil {
		opts.Secure = true
	}

	natsConn, err := opts.Connect()
	if err != nil {
		return err
	}

	log.Bg().Info("connection using clientID", zap.String("clientID", conf.ClientID))
	c, err := stan.Connect(conf.ClusterID, conf.ClientID, stan.NatsConn(natsConn),
		stan.SetConnectionLostHandler(func(_ stan.Conn, err error) {
			log.Bg().Error("connection lost to nats cluster", zap.Error(err), zap.Any("brokerConfig", n.opts), zap.Any("stanConfig", conf))
			errCtx := runtime.Cancel(n.opts.Context)
			if errCtx != nil {
				log.For(n.opts.Context).Error("error calling cancel", zap.Error(errCtx))
			}
		}))
	if err != nil {
		return err
	}
	n.conn = c
	return nil
}

func (n *nBroker) Disconnect() error {
	if err := n.conn.Close(); err != nil {
		log.Bg().Error("Unable to close connection to nats cluster", zap.Error(err))
	}
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

	fn := func(msg *stan.Msg) {

		log.Bg().Debug("new msg received from nats", zap.Int("length", len(msg.MsgProto.Data)))

		var m broker.Message
		if err := n.opts.Codec.Unmarshal(msg.Data, &m); err != nil {
			log.Bg().Error("Unable to decode subscription message", zap.Error(err), zap.String("subject", topic))
			return
		}
		err := handler(&publication{m: &m, r: msg})
		if err != nil {
			log.Bg().Error("Unable to register subscription handler", zap.Error(err), zap.String("subject", topic))
		}
	}

	var sub stan.Subscription
	var err error
	conf := opt.Context.Value(broker.ContextKey(SubsConfName)).(SubscribeOpts)
	stanOpts := []stan.SubscriptionOption{}

	if len(conf.DurableName) > 0 {
		stanOpts = append(stanOpts, stan.DurableName(conf.DurableName))
	}
	if len(conf.QueueGroup) > 0 && opt.Queue != conf.QueueGroup {
		opt.Queue = conf.QueueGroup
	}
	if conf.DeliverAllAvailable {
		stanOpts = append(stanOpts, stan.DeliverAllAvailable())
	}
	if conf.ManualAcks {
		stanOpts = append(stanOpts, stan.SetManualAckMode())
	}
	if conf.StartSequence > 0 {
		stanOpts = append(stanOpts, stan.StartAtSequence(conf.StartSequence))
	}

	if len(opt.Queue) > 0 {
		sub, err = n.conn.QueueSubscribe(topic, opt.Queue, fn, stanOpts...)
	} else {
		sub, err = n.conn.Subscribe(topic, fn, stanOpts...)
	}
	if err != nil {
		return nil, err
	}
	return &subscriber{s: sub, opts: opt, subject: topic}, nil
}

func (n *nBroker) String() string {
	return "stan"
}
