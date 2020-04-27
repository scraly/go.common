/*
 * Copyright (c) Continental Corporation - All Rights Reserved
 *
 * This file is a part of Entry project.
 * ITS France - Entry squad members
 *
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

package eventbus

import (
	"fmt"
	"reflect"
	"sync"
)

// EventBus - box for handlers and callbacks.
type eventBus struct {
	handlers map[string][]*eventHandler
	lock     sync.Mutex // a lock for the map
	wg       sync.WaitGroup
}

type eventHandler struct {
	callBack      reflect.Value
	flagOnce      bool
	async         bool
	transactional bool
	called        bool
	sync.Mutex    // lock for an event handler - useful for running async callbacks serially
}

// NewLocal returns new EventBus with empty handlers.
func NewLocal() EventBus {
	return &eventBus{
		make(map[string][]*eventHandler),
		sync.Mutex{},
		sync.WaitGroup{},
	}
}

// Subscribe subscribes to a topic.
// Returns error if `fn` is not a function.
func (bus *eventBus) Subscribe(topic string, fn interface{}) error {
	bus.lock.Lock()
	defer bus.lock.Unlock()
	if !(reflect.TypeOf(fn).Kind() == reflect.Func) {
		return fmt.Errorf("%s is not of type reflect.Func", reflect.TypeOf(fn).Kind())
	}
	v := reflect.ValueOf(fn)
	bus.handlers[topic] = append(bus.handlers[topic], &eventHandler{
		v, false, false, false, false, sync.Mutex{},
	})
	return nil
}

// SubscribeAsync subscribes to a topic with an asynchronous callback
// Transactional determines whether subsequent callbacks for a topic are
// run serially (true) or concurrently (false)
// Returns error if `fn` is not a function.
func (bus *eventBus) SubscribeAsync(topic string, fn interface{}, transactional bool) error {
	bus.lock.Lock()
	defer bus.lock.Unlock()
	if !(reflect.TypeOf(fn).Kind() == reflect.Func) {
		return fmt.Errorf("%s is not of type reflect.Func", reflect.TypeOf(fn).Kind())
	}
	v := reflect.ValueOf(fn)
	bus.handlers[topic] = append(bus.handlers[topic], &eventHandler{
		v, false, true, transactional, false, sync.Mutex{},
	})
	return nil
}

// SubscribeOnce subscribes to a topic once. Handler will be removed after executing.
// Returns error if `fn` is not a function.
func (bus *eventBus) SubscribeOnce(topic string, fn interface{}) error {
	bus.lock.Lock()
	defer bus.lock.Unlock()
	if !(reflect.TypeOf(fn).Kind() == reflect.Func) {
		return fmt.Errorf("%s is not of type reflect.Func", reflect.TypeOf(fn).Kind())
	}
	v := reflect.ValueOf(fn)
	bus.handlers[topic] = append(bus.handlers[topic], &eventHandler{
		v, true, false, false, false, sync.Mutex{},
	})
	return nil
}

// SubscribeOnceAsync subscribes to a topic once with an asyncrhonous callback
// Handler will be removed after executing.
// Returns error if `fn` is not a function.
func (bus *eventBus) SubscribeOnceAsync(topic string, fn interface{}) error {
	bus.lock.Lock()
	defer bus.lock.Unlock()
	if !(reflect.TypeOf(fn).Kind() == reflect.Func) {
		return fmt.Errorf("%s is not of type reflect.Func", reflect.TypeOf(fn).Kind())
	}
	v := reflect.ValueOf(fn)
	bus.handlers[topic] = append(bus.handlers[topic], &eventHandler{
		v, true, true, false, false, sync.Mutex{},
	})
	return nil
}

// HasCallback returns true if exists any callback subscribed to the topic.
func (bus *eventBus) HasCallback(topic string) bool {
	bus.lock.Lock()
	defer bus.lock.Unlock()
	_, ok := bus.handlers[topic]
	if ok {
		return len(bus.handlers[topic]) > 0
	}
	return false
}

// Unsubscribe removes callback defined for a topic.
// Returns error if there are no callbacks subscribed to the topic.
func (bus *eventBus) Unsubscribe(topic string, handler interface{}) error {
	bus.lock.Lock()
	defer bus.lock.Unlock()
	if _, ok := bus.handlers[topic]; ok && len(bus.handlers[topic]) > 0 {
		bus.removeHandler(topic, reflect.ValueOf(handler))
		return nil
	}
	return fmt.Errorf("topic %s doesn't exist", topic)
}

// Publish executes callback defined for a topic. Any additional argument will be tranfered to the callback.
func (bus *eventBus) Publish(topic string, args ...interface{}) {
	bus.lock.Lock() // will unlock if handler is not found or always after setUpPublish
	defer bus.lock.Unlock()
	if handlers, ok := bus.handlers[topic]; ok {
		for _, handler := range handlers {
			if !handler.async {
				bus.doPublish(handler, topic, args...)
			} else {
				bus.wg.Add(1)
				go bus.doPublishAsync(handler, topic, args...)
			}
		}
	}
}

func (bus *eventBus) doPublish(handler *eventHandler, topic string, args ...interface{}) {
	passedArguments := bus.setUpPublish(topic, args...)
	if handler.flagOnce {
		bus.removeHandler(topic, handler.callBack)
		if handler.called {
			return
		}
	}
	handler.called = true
	handler.callBack.Call(passedArguments)
}

func (bus *eventBus) doPublishAsync(handler *eventHandler, topic string, args ...interface{}) {
	defer bus.wg.Done()
	if handler.transactional {
		handler.Lock()
		defer handler.Unlock()
	}
	bus.doPublish(handler, topic, args...)
}

func (bus *eventBus) findHandlerIdx(topic string, callback reflect.Value) int {
	if _, ok := bus.handlers[topic]; ok {
		for idx, handler := range bus.handlers[topic] {
			if handler.callBack == callback {
				return idx
			}
		}
	}
	return -1
}

func (bus *eventBus) removeHandler(topic string, callback reflect.Value) {
	i := bus.findHandlerIdx(topic, callback)
	if i >= 0 {
		bus.handlers[topic] = append(bus.handlers[topic][:i], bus.handlers[topic][i+1:]...)
	}
}

func (bus *eventBus) setUpPublish(_ string, args ...interface{}) []reflect.Value {

	var passedArguments = make([]reflect.Value, 0)
	for _, arg := range args {
		passedArguments = append(passedArguments, reflect.ValueOf(arg))
	}
	return passedArguments
}

// WaitAsync waits for all async callbacks to complete
func (bus *eventBus) WaitAsync() {
	bus.wg.Wait()
}
