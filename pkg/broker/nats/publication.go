/*
 * Copyright (c) Continental Corporation - All Rights Reserved
 *
 * This file is a part of Entry project.
 * ITS France - Entry squad members
 *
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

package nats

import "github.com/scraly/go.common/pkg/broker"

type publication struct {
	t string
	m *broker.Message
}

// -----------------------------------------------------------------------------

func (n *publication) Topic() string {
	return n.t
}

func (n *publication) Message() *broker.Message {
	return n.m
}

func (n *publication) Ack() error {
	return nil
}

func (n *publication) Seq() uint64 {
	return 0
}

func (n *publication) Timestamp() int64 {
	return 0
}
