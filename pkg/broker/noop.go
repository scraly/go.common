/*
 * Copyright (c) Continental Corporation - All Rights Reserved
 *
 * This file is a part of Entry project.
 * ITS France - Entry squad members
 *
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

package broker

type noopBroker struct {
}

// -----------------------------------------------------------------------------
func (b noopBroker) Options() Options {
	return Options{}
}

func (b noopBroker) Address() string {
	return ""
}

func (b noopBroker) Connect() error {
	return nil
}

func (b noopBroker) Disconnect() error {
	return nil
}

func (b noopBroker) Init(...Option) error {
	return nil
}

func (b noopBroker) Publish(string, *Message, ...PublishOption) error {
	return nil
}

func (b noopBroker) Subscribe(string, Handler, ...SubscribeOption) (Subscriber, error) {
	return nil, nil
}

func (b noopBroker) String() string {
	return "noop"
}
