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

import (
	"context"
	"crypto/tls"

	"github.com/scraly/go.common/pkg/storage/codec"
)

// ContextKey is used to generate keys for Context value
type ContextKey string

func (c ContextKey) String() string {
	return "bCtxKeys" + string(c)
}

// Options is connection option holder
type Options struct {
	Addrs     []string
	Secure    bool
	Codec     codec.Codec
	TLSConfig *tls.Config
	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

// PublishOptions is publication option holder
type PublishOptions struct {
	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

// SubscribeOptions is subscription option holder
type SubscribeOptions struct {
	// AutoAck defaults to true. When a handler returns
	// with a nil error the message is acked.
	AutoAck bool
	// Subscribers with the same queue name
	// will create a shared subscription where each
	// receives a subset of messages.
	Queue string

	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

// Option represents default option function
type Option func(*Options)

// PublishOption represents publication option function
type PublishOption func(*PublishOptions)

// SubscribeOption represents subscription option function
type SubscribeOption func(*SubscribeOptions)

// Addrs sets the host addresses to be used by the broker
func Addrs(addrs ...string) Option {
	return func(o *Options) {
		o.Addrs = addrs
	}
}

// Codec sets the codec used for encoding/decoding used where
// a broker does not support headers
func Codec(c codec.Codec) Option {
	return func(o *Options) {
		o.Codec = c
	}
}

// DisableAutoAck will disable auto acking of messages
// after they have been handled.
func DisableAutoAck() SubscribeOption {
	return func(o *SubscribeOptions) {
		o.AutoAck = false
	}
}

// Queue sets the name of the queue to share messages on
func Queue(name string) SubscribeOption {
	return func(o *SubscribeOptions) {
		o.Queue = name
	}
}

// Secure communication with the broker
func Secure(b bool) Option {
	return func(o *Options) {
		o.Secure = b
	}
}

// TLSConfig sets the TLS connection settings
func TLSConfig(t *tls.Config) Option {
	return func(o *Options) {
		o.TLSConfig = t
	}
}
