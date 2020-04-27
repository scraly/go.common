/*
 * Copyright (c) Continental Corporation - All Rights Reserved
 *
 * This file is a part of Entry project.
 * ITS France - Entry squad members
 *
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

package redis

import (
	"time"

	api "github.com/scraly/go.common/pkg/cache"
	"github.com/scraly/go.common/pkg/log"
	"github.com/scraly/go.common/pkg/storage/codec/msgpack"

	"github.com/garyburd/redigo/redis"
)

type redisStore struct {
	pool *redis.Pool

	opts api.Options
}

// NewCacheStore initializes a redis cache
func NewCacheStore(...api.Option) api.Store {
	options := api.Options{
		// default to msgpack codec
		Codec: msgpack.NewCodec(),
	}

	return &redisStore{
		opts: options,
	}
}

func init() {
	api.Register("redis", NewCacheStore)
}

// -----------------------------------------------------------------------------
func (s *redisStore) Name() string {
	return "redis"
}

func (s *redisStore) Connect() error {
	s.pool = &redis.Pool{
		MaxIdle:     5,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", s.opts.Addrs[0])
			if err != nil {
				return nil, err
			}

			if len(s.opts.Password) > 0 {
				if _, errAuth := c.Do("AUTH", s.opts.Password); errAuth != nil {
					log.SafeClose(c, "Unable to close redis connection")
					return nil, errAuth
				}
			} else {
				// check with PING
				if _, errPing := c.Do("PING"); errPing != nil {
					log.SafeClose(c, "Unable to close redis connection")
					return nil, errPing
				}
			}

			return c, nil
		},
		// custom connection test method
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	// Return no error
	return nil
}

func (s *redisStore) Get(key string, value interface{}) error {
	conn := s.pool.Get()
	defer func(conn redis.Conn) {
		log.SafeClose(conn, "Unable to close redis connection")
	}(conn)

	// Retrieve from backend
	raw, err := conn.Do("GET", key)
	if raw == nil {
		return api.ErrCacheMiss
	}

	// Decode value as byte array
	item, err := redis.Bytes(raw, err)
	if err != nil {
		return err
	}

	// Defer to codec for unmarshalling responsibility
	return s.opts.Codec.Unmarshal(item, value)
}

func (s *redisStore) Set(key string, value interface{}, expires time.Duration) error {
	return s.invoke(s.pool.Get().Do, key, value, expires)
}

func (s *redisStore) Add(key string, value interface{}, expires time.Duration) error {
	conn := s.pool.Get()
	defer func(conn redis.Conn) {
		log.SafeClose(conn, "Unable to close redis connection")
	}(conn)

	ok, err := exists(conn, key)
	if err != nil {
		return err
	}
	if ok {
		return api.ErrNotStored
	}

	return s.invoke(conn.Do, key, value, expires)
}

func (s *redisStore) Replace(key string, value interface{}, expires time.Duration) error {
	conn := s.pool.Get()
	defer func(conn redis.Conn) {
		log.SafeClose(conn, "Unable to close redis connection")
	}(conn)

	ok, err := exists(conn, key)
	if err != nil {
		return err
	}
	if !ok {
		return api.ErrNotStored
	}

	err = s.invoke(conn.Do, key, value, expires)
	if value == nil {
		return api.ErrNotStored
	}

	return err
}

func (s *redisStore) Delete(key string) error {
	conn := s.pool.Get()
	defer func(conn redis.Conn) {
		log.SafeClose(conn, "Unable to close redis connection")
	}(conn)

	ok, err := exists(conn, key)
	if err != nil {
		return err
	}
	if !ok {
		return api.ErrCacheMiss
	}

	_, err = conn.Do("DEL", key)
	return err
}

func (s *redisStore) Flush() error {
	conn := s.pool.Get()
	defer func(conn redis.Conn) {
		log.SafeClose(conn, "Unable to close redis connection")
	}(conn)

	_, err := conn.Do("FLUSHALL")
	return err
}

// -----------------------------------------------------------------------------

func exists(conn redis.Conn, key string) (bool, error) {
	return redis.Bool(conn.Do("EXISTS", key))
}

func (s *redisStore) invoke(f func(string, ...interface{}) (interface{}, error),
	key string, value interface{}, expires time.Duration) error {

	switch expires {
	case api.DEFAULT:
		expires = s.opts.DefaultExpiration
	case api.FOREVER:
		expires = time.Duration(0)
	}

	b, err := s.opts.Codec.Marshal(value)
	if err != nil {
		return err
	}

	conn := s.pool.Get()
	defer func(conn redis.Conn) {
		log.SafeClose(conn, "Unable to close redis connection")
	}(conn)
	if expires > 0 {
		_, errSetEx := f("SETEX", key, int32(expires/time.Second), b)
		return errSetEx
	}

	_, err = f("SET", key, b)
	return err
}
