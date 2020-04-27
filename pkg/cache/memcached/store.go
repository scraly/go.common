/*
 * Copyright (c) Continental Corporation - All Rights Reserved
 *
 * This file is a part of Entry project.
 * ITS France - Entry squad members
 *
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

package memcached

import (
	"time"

	api "github.com/scraly/go.common/pkg/cache"
	"github.com/scraly/go.common/pkg/storage/codec/msgpack"

	"github.com/bradfitz/gomemcache/memcache"
)

type memcachedStore struct {
	*memcache.Client

	opts api.Options
}

// NewCacheStore initializes a memcached cache
func NewCacheStore(...api.Option) api.Store {
	options := api.Options{
		// default to msgpack codec
		Codec: msgpack.NewCodec(),
	}

	return &memcachedStore{
		opts: options,
	}
}

func init() {
	api.Register("memcached", NewCacheStore)
}

// -----------------------------------------------------------------------------
func (s *memcachedStore) Name() string {
	return "memcached"
}

func (s *memcachedStore) Connect() error {
	s.Client = memcache.New(s.opts.Addrs...)

	// Return no error
	return nil
}

func (s *memcachedStore) Get(key string, value interface{}) error {
	item, err := s.Client.Get(key)
	if err != nil {
		return convertMemcacheError(err)
	}
	return s.opts.Codec.Unmarshal(item.Value, value)
}

func (s *memcachedStore) Set(key string, value interface{}, expires time.Duration) error {
	return s.invoke((*memcache.Client).Set, key, value, expires)
}

func (s *memcachedStore) Add(key string, value interface{}, expires time.Duration) error {
	return s.invoke((*memcache.Client).Add, key, value, expires)
}

func (s *memcachedStore) Replace(key string, value interface{}, expires time.Duration) error {
	return s.invoke((*memcache.Client).Replace, key, value, expires)
}

func (s *memcachedStore) Delete(key string) error {
	return convertMemcacheError(s.Client.Delete(key))
}

func (s *memcachedStore) Flush() error {
	return api.ErrNotSupported
}

// -----------------------------------------------------------------------------

func (s *memcachedStore) invoke(storeFn func(*memcache.Client, *memcache.Item) error,
	key string, value interface{}, expire time.Duration) error {

	switch expire {
	case api.DEFAULT:
		expire = s.opts.DefaultExpiration
	case api.FOREVER:
		expire = time.Duration(0)
	}

	b, err := s.opts.Codec.Marshal(value)
	if err != nil {
		return err
	}
	return convertMemcacheError(storeFn(s.Client, &memcache.Item{
		Key:        key,
		Value:      b,
		Expiration: int32(expire / time.Second),
	}))
}

func convertMemcacheError(err error) error {
	switch err {
	case nil:
		return nil
	case memcache.ErrCacheMiss:
		return api.ErrCacheMiss
	case memcache.ErrNotStored:
		return api.ErrNotStored
	}

	return err
}
