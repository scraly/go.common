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
	"github.com/scraly/go.common/pkg/broker"
	"github.com/scraly/go.common/pkg/broker/rabbitmq/internal"
)

type subscriber struct {
	opts  broker.SubscribeOptions
	topic string
	ch    *internal.Channel
}

func (s *subscriber) Options() broker.SubscribeOptions {
	return s.opts
}

func (s *subscriber) Topic() string {
	return s.topic
}

func (s *subscriber) Unsubscribe() error {
	return s.ch.Close()
}
