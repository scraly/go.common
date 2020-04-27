/*
 * Copyright (c) Continental Corporation - All Rights Reserved
 *
 * This file is a part of Entry project.
 * ITS France - Entry squad members
 *
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

package stan

import (
	"github.com/scraly/go.common/pkg/broker"
	stan "github.com/nats-io/go-nats-streaming"
)

type publication struct {
	m *broker.Message
	r *stan.Msg
}

// -----------------------------------------------------------------------------

//Topic:  The NATS Streaming delivery subject
func (n *publication) Topic() string {
	return n.r.Subject
}

func (n *publication) Message() *broker.Message {
	return n.m
}

func (n *publication) Ack() error {
	return n.r.Ack()
}

//Seq: a globally ordered sequence number for the subject’s channel
func (n *publication) Seq() uint64 {
	return n.r.Sequence
}

//Timestamp: the received timestamp, in nanoseconds
func (n *publication) Timestamp() int64 {
	return n.r.Timestamp
}
