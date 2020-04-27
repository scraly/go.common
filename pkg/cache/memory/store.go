/*
 * Copyright (c) Continental Corporation - All Rights Reserved
 *
 * This file is a part of Entry project.
 * ITS France - Entry squad members
 *
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

package memory

import (
	"reflect"
	"time"

	api "github.com/scraly/go.common/pkg/cache"
	"github.com/scraly/go.common/pkg/storage/codec/msgpack"

	cache "github.com/robfig/go-cache"
)

type memoryStore struct {
	cache.Cache
	opts api.Options
}

// NewCacheStore initializes a memcached cache
func NewCacheStore(...api.Option) api.Store {
	options := api.Options{
		// default to msgpack codec
		Codec: msgpack.NewCodec(),
	}

	return &memoryStore{
		opts: options,
	}
}

func init() {
	api.Register("memory", NewCacheStore)
}

// -----------------------------------------------------------------------------
func (s *memoryStore) Name() string {
	return "memory"
}

func (s *memoryStore) Connect() error {
	s.Cache = *cache.New(s.opts.DefaultExpiration, time.Minute)
	// Return no error
	return nil
}

func (s *memoryStore) Get(key string, value interface{}) error {
	val, found := s.Cache.Get(key)
	if !found {
		return api.ErrCacheMiss
	}

	v := reflect.ValueOf(value)
	if v.Type().Kind() == reflect.Ptr && v.Elem().CanSet() {
		v.Elem().Set(reflect.ValueOf(val))
		return nil
	}

	return api.ErrNotStored
}

func (s *memoryStore) Set(key string, value interface{}, expires time.Duration) error {
	s.Cache.Set(key, value, expires)
	return nil
}

func (s *memoryStore) Add(key string, value interface{}, expires time.Duration) error {
	err := s.Cache.Add(key, value, expires)
	if err == cache.ErrKeyExists {
		return api.ErrNotStored
	}
	return err
}

func (s *memoryStore) Replace(key string, value interface{}, expires time.Duration) error {
	if err := s.Cache.Replace(key, value, expires); err != nil {
		return api.ErrNotStored
	}
	return nil
}

func (s *memoryStore) Delete(key string) error {
	if found := s.Cache.Delete(key); !found {
		return api.ErrCacheMiss
	}
	return nil
}

func (s *memoryStore) Flush() error {
	s.Cache.Flush()
	return nil
}
