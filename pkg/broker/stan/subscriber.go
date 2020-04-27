/*
 * Copyright (c) Continental Corporation - All Rights Reserved
 *
 * This file is a part of Entry project.
 * ITS France - Entry squad members
 *
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

// Package stan provides a NATS Streaming broker
package stan

import (
	"github.com/scraly/go.common/pkg/broker"

	stan "github.com/nats-io/go-nats-streaming"
)

type subscriber struct {
	s       stan.Subscription
	opts    broker.SubscribeOptions
	subject string
}

func (n *subscriber) Options() broker.SubscribeOptions {
	return n.opts
}

func (n *subscriber) Topic() string {
	return n.subject
}

func (n *subscriber) Unsubscribe() error {
	return n.s.Unsubscribe()
}
