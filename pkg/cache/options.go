/*
 * Copyright (c) Continental Corporation - All Rights Reserved
 *
 * This file is a part of Entry project.
 * ITS France - Entry squad members
 *
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

package cache

import (
	"crypto/tls"
	"time"

	"github.com/scraly/go.common/pkg/storage/codec"
)

// Options is connection option holder
type Options struct {
	Addrs             []string
	Secure            bool
	Codec             codec.Codec
	TLSConfig         *tls.Config
	Username          string
	Password          string
	DefaultExpiration time.Duration
}

// Option represents default option function
type Option func(*Options)

// Addrs sets the host addresses to be used by the cache manager
func Addrs(addrs ...string) Option {
	return func(o *Options) {
		o.Addrs = addrs
	}
}

// Codec sets the codec used for encoding/decoding
func Codec(c codec.Codec) Option {
	return func(o *Options) {
		o.Codec = c
	}
}

// Secure communication with the cache manager
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

// Username sets the username value for connection
func Username(value string) Option {
	return func(o *Options) {
		o.Username = value
	}
}

// Password sets the password value for connection
func Password(value string) Option {
	return func(o *Options) {
		o.Password = value
	}
}

// DefaultExpiration sets the expiration value for all objects in cache
func DefaultExpiration(value time.Duration) Option {
	return func(o *Options) {
		o.DefaultExpiration = value
	}
}
