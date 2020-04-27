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

//go:generate mockery -name EventBus -output mock

// EventBus is the contract for all Event Bus implementations
type EventBus interface {
	Subscribe(topic string, fn interface{}) error
	SubscribeAsync(topic string, fn interface{}, transactional bool) error
	SubscribeOnce(topic string, fn interface{}) error
	SubscribeOnceAsync(topic string, fn interface{}) error
	HasCallback(topic string) bool
	Unsubscribe(topic string, handler interface{}) error
	Publish(topic string, args ...interface{})
	WaitAsync()
}
