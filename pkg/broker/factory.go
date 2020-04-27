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
	"errors"

	"github.com/scraly/go.common/pkg/log"
	"go.uber.org/zap"
)

// FactoryFunc is the broker constructor function
type FactoryFunc func(...Option) Broker

var (
	brokers = map[string]FactoryFunc{}

	// ErrInvalidBrokerFactory is raised when trying to use an invalid broker factory name
	ErrInvalidBrokerFactory = errors.New("broker: given broker not supported")
)

// Register a broker backend factory
func Register(name string, constructor FactoryFunc) {
	if _, ok := brokers[name]; !ok {
		brokers[name] = constructor
	} else {
		log.Bg().Fatal("Broker factory already registered !", zap.String("name", name))
	}
}

// Registered returns the list of registered brokers
func Registered() []string {
	// Convert map to slice of keys.
	keys := []string{}
	for key := range brokers {
		keys = append(keys, key)
	}
	return keys
}

// New returns a broker instance
func New(name string, opts ...Option) (Broker, error) {
	if builder, ok := brokers[name]; ok {
		b := builder(opts...)
		return b, b.Connect()
	}
	return nil, ErrInvalidBrokerFactory
}
