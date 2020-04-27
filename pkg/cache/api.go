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
	"errors"
	"time"
)

//go:generate mockery -name Store

// Store is an interface used for cache backend operations
type Store interface {
	Name() string
	Connect() error
	Get(key string, value interface{}) error
	Set(key string, value interface{}, expire time.Duration) error
	Add(key string, value interface{}, expire time.Duration) error
	Replace(key string, data interface{}, expire time.Duration) error
	Delete(key string) error
	Flush() error
}

const (
	// DEFAULT stores value according default cache value
	DEFAULT = time.Duration(0)
	// FOREVER stores value forever, no eviction
	FOREVER = time.Duration(-1)
)

var (
	// CachePrefix defines key spec prefix
	CachePrefix = "default"

	// ErrCacheMiss is raised when key is not found in cache
	ErrCacheMiss = errors.New("cache: key not found")
	// ErrNotStored is raised when trying to insert an object in the cache and failure occurs
	ErrNotStored = errors.New("cache: not stored")
	// ErrNotSupported is raise dwhen trying to invoke a not supported operation
	ErrNotSupported = errors.New("cache: operation not supported")
)
