/*
 * Copyright (c) Continental Corporation - All Rights Reserved
 *
 * This file is a part of Entry project.
 * ITS France - Entry squad members
 *
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

// Package nats provides a NATS broker
package nats

import (
	"github.com/scraly/go.common/pkg/broker"

	nats "github.com/nats-io/go-nats"
)

type subscriber struct {
	s    *nats.Subscription
	opts broker.SubscribeOptions
}

func (n *subscriber) Options() broker.SubscribeOptions {
	return n.opts
}

func (n *subscriber) Topic() string {
	return n.s.Subject
}

func (n *subscriber) Unsubscribe() error {
	return n.s.Unsubscribe()
}
