/*
 * Copyright (c) Continental Corporation - All Rights Reserved
 *
 * This file is a part of Entry project.
 * ITS France - Entry squad members
 *
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

package kafka

import (
	api "github.com/scraly/go.common/pkg/broker"
	"github.com/scraly/go.common/pkg/log"
	"github.com/scraly/go.common/pkg/storage/codec/msgpack"

	"github.com/Shopify/sarama"
	"github.com/pborman/uuid"
	"go.uber.org/zap"
	sc "gopkg.in/bsm/sarama-cluster.v2"
)

type broker struct {
	addrs []string

	c  sarama.Client
	p  sarama.SyncProducer
	sc *sc.Client

	opts api.Options
}

// NewBroker initializes a Kafka broker instance
func NewBroker(opts ...api.Option) api.Broker {
	options := api.Options{
		// default to msgpack codec
		Codec: msgpack.NewCodec(),
	}

	for _, o := range opts {
		o(&options)
	}

	var cAddrs []string
	for _, addr := range options.Addrs {
		if len(addr) == 0 {
			continue
		}
		cAddrs = append(cAddrs, addr)
	}
	if len(cAddrs) == 0 {
		cAddrs = []string{"127.0.0.1:9092"}
	}

	return &broker{
		addrs: cAddrs,
		opts:  options,
	}
}

func init() {
	api.Register("kafka", NewBroker)
}

// -----------------------------------------------------------------------------

func (k *broker) Address() string {
	if len(k.addrs) > 0 {
		return k.addrs[0]
	}
	return "127.0.0.1:9092"
}

func (k *broker) Connect() error {
	if k.c != nil {
		return nil
	}

	pconfig := sarama.NewConfig()
	// For implementation reasons, the SyncProducer requires
	// `Producer.Return.Errors` and `Producer.Return.Successes`
	// to be set to true in its configuration.
	pconfig.Producer.Return.Successes = true
	pconfig.Producer.Return.Errors = true

	c, err := sarama.NewClient(k.addrs, pconfig)
	if err != nil {
		return err
	}

	k.c = c

	p, err := sarama.NewSyncProducerFromClient(c)
	if err != nil {
		return err
	}

	k.p = p

	config := sc.NewConfig()
	// TODO: make configurable offset as SubscriberOption
	config.Config.Consumer.Offsets.Initial = sarama.OffsetNewest

	cs, err := sc.NewClient(k.addrs, config)
	if err != nil {
		return err
	}

	k.sc = cs
	// TODO: TLS
	/*
		opts.Secure = k.opts.Secure
		opts.TLSConfig = k.opts.TLSConfig

		// secure might not be set
		if k.opts.TLSConfig != nil {
			opts.Secure = true
		}
	*/
	return nil
}

func (k *broker) Disconnect() error {
	log.SafeClose(k.sc, "Error while closing Kafka client")
	log.SafeClose(k.p, "Error while closing Kafka producer")
	return k.c.Close()
}

func (k *broker) Init(opts ...api.Option) error {
	for _, o := range opts {
		o(&k.opts)
	}
	return nil
}

func (k *broker) Options() api.Options {
	return k.opts
}

func (k *broker) Publish(topic string, msg *api.Message, _ ...api.PublishOption) error {
	b, err := k.opts.Codec.Marshal(msg)
	if err != nil {
		return err
	}
	_, _, err = k.p.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(b),
	})
	return err
}

func (k *broker) Subscribe(topic string, handler api.Handler, opts ...api.SubscribeOption) (api.Subscriber, error) {
	opt := api.SubscribeOptions{
		AutoAck: true,
		Queue:   uuid.NewUUID().String(),
	}

	for _, o := range opts {
		o(&opt)
	}

	c, err := sc.NewConsumerFromClient(k.sc, opt.Queue, []string{topic})
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			select {
			case err := <-c.Errors():
				log.Bg().Error("consumer error", zap.Error(err))
			case sm := <-c.Messages():
				// ensure message is not nil
				if sm == nil {
					continue
				}
				var m api.Message
				if err := k.opts.Codec.Unmarshal(sm.Value, &m); err != nil {
					continue
				}
				if err := handler(&publication{
					m:  &m,
					t:  sm.Topic,
					c:  c,
					km: sm,
				}); err == nil && opt.AutoAck {
					c.MarkOffset(sm, "")
				}
			}
		}
	}()

	return &subscriber{s: c, opts: opt}, nil
}

func (k *broker) String() string {
	return "kafka"
}
