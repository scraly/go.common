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
	sc "gopkg.in/bsm/sarama-cluster.v2"
)

type subscriber struct {
	s    *sc.Consumer
	t    string
	opts api.SubscribeOptions
}

// -----------------------------------------------------------------------------

func (s *subscriber) Options() api.SubscribeOptions {
	return s.opts
}

func (s *subscriber) Topic() string {
	return s.t
}

func (s *subscriber) Unsubscribe() error {
	return s.s.Close()
}
