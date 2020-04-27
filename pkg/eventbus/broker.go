package eventbus

import (
	"github.com/scraly/go.common/pkg/broker"
	"github.com/scraly/go.common/pkg/log"
	"github.com/scraly/go.common/pkg/storage/codec"
	"github.com/scraly/go.common/pkg/storage/codec/json"

	"go.uber.org/zap"
)

// brokeredBus - box for handlers and callbacks.
type brokeredBus struct {
	next       EventBus
	remote     broker.Broker
	marshaller codec.Codec
}

// BrokeredBus returns new EventBus with empty handlers.
func BrokeredBus(next EventBus, remote broker.Broker) EventBus {
	return &brokeredBus{
		next:       next,
		remote:     remote,
		marshaller: json.NewCodec(),
	}
}

// Subscribe subscribes to a topic.
// Returns error if `fn` is not a function.
func (bus *brokeredBus) Subscribe(topic string, fn interface{}) error {
	return bus.next.Subscribe(topic, fn)
}

// SubscribeAsync subscribes to a topic with an asynchronous callback
// Transactional determines whether subsequent callbacks for a topic are
// run serially (true) or concurrently (false)
// Returns error if `fn` is not a function.
func (bus *brokeredBus) SubscribeAsync(topic string, fn interface{}, transactional bool) error {
	return bus.next.SubscribeAsync(topic, fn, transactional)
}

// SubscribeOnce subscribes to a topic once. Handler will be removed after executing.
// Returns error if `fn` is not a function.
func (bus *brokeredBus) SubscribeOnce(topic string, fn interface{}) error {
	return bus.next.SubscribeOnce(topic, fn)
}

// SubscribeOnceAsync subscribes to a topic once with an asyncrhonous callback
// Handler will be removed after executing.
// Returns error if `fn` is not a function.
func (bus *brokeredBus) SubscribeOnceAsync(topic string, fn interface{}) error {
	return bus.next.SubscribeOnceAsync(topic, fn)
}

// HasCallback returns true if exists any callback subscribed to the topic.
func (bus *brokeredBus) HasCallback(topic string) bool {
	return bus.next.HasCallback(topic)
}

// Unsubscribe removes callback defined for a topic.
// Returns error if there are no callbacks subscribed to the topic.
func (bus *brokeredBus) Unsubscribe(topic string, handler interface{}) error {
	return bus.next.Unsubscribe(topic, handler)
}

// Publish executes callback defined for a topic. Any additional argument will be tranfered to the callback.
func (bus *brokeredBus) Publish(topic string, args ...interface{}) {
	bus.next.Publish(topic, args...)
	go func() {
		// Encode arguments
		payload, err := bus.marshaller.Marshal(args)
		if err != nil {
			log.Bg().Error("Unable to encode event", zap.String("topic", topic), zap.Any("args", args), zap.Error(err))
			return
		}

		// Send
		if err := bus.remote.Publish(topic, &broker.Message{
			Body: payload,
		}); err != nil {
			log.Bg().Error("Unable to publish to remote broker", zap.String("topic", topic), zap.Any("args", args), zap.Error(err))
		}
	}()
}

// WaitAsync waits for all async callbacks to complete
func (bus *brokeredBus) WaitAsync() {
	bus.next.WaitAsync()
}
