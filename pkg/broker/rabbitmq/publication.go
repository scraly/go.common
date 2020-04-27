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

	"github.com/streadway/amqp"
)

type publication struct {
	d amqp.Delivery
	m *broker.Message
	t string
}

func (p *publication) Ack() error {
	return p.d.Ack(false)
}

func (p *publication) Topic() string {
	return p.t
}

func (p *publication) Message() *broker.Message {
	return p.m
}

func (p *publication) Seq() uint64 {
	return p.d.DeliveryTag
}

func (p *publication) Timestamp() int64 {

	ts := p.d.Timestamp.UnixNano()

	return ts
}
