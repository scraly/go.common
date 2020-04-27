/*
 * Copyright (c) Continental Corporation - All Rights Reserved
 *
 * This file is a part of Entry project.
 * ITS France - Entry squad members
 *
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

// Package kafka provides a kafka broker using sarama cluster
package kafka

import (
	api "github.com/scraly/go.common/pkg/broker"

	"github.com/Shopify/sarama"
	sc "gopkg.in/bsm/sarama-cluster.v2"
)

type publication struct {
	t  string
	c  *sc.Consumer
	km *sarama.ConsumerMessage
	m  *api.Message
}

// -----------------------------------------------------------------------------

func (p *publication) Topic() string {
	return p.t
}

func (p *publication) Message() *api.Message {
	return p.m
}

func (p *publication) Ack() error {
	p.c.MarkOffset(p.km, "")
	return nil
}

func (p *publication) Seq() uint64 {
	return uint64(p.km.Offset)
}

func (p *publication) Timestamp() int64 {
	return 0
}
