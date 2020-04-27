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

	"github.com/scraly/go.common/pkg/log"
	"go.uber.org/zap"
)

// FactoryFunc is the broker constructor function
type FactoryFunc func(...Option) Store

var (
	managers = map[string]FactoryFunc{}

	// ErrInvalidCacheFactory is raised when trying to use an invalid cache factory name
	ErrInvalidCacheFactory = errors.New("cache: given cache manager is not supported")
)

// Register a broker backend factory
func Register(name string, constructor FactoryFunc) {
	if _, ok := managers[name]; !ok {
		managers[name] = constructor
	} else {
		log.Bg().Fatal("Cache manager factory already registered !", zap.String("name", name))
	}
}

// New returns a broker instance
func New(name string, opts ...Option) (Store, error) {
	if builder, ok := managers[name]; ok {
		b := builder(opts...)
		return b, b.Connect()
	}
	return nil, ErrInvalidCacheFactory
}
