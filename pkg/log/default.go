/*
 * Copyright (c) Continental Corporation - All Rights Reserved
 *
 * This file is a part of Entry project.
 * ITS France - Entry squad members
 *
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

package log

import (
	"context"
	l "log"

	"go.uber.org/zap"
)

var defaultFactory LoggerFactory

// -----------------------------------------------------------------------------

func init() {
	defaultLogger, err := zap.NewProduction()
	if err != nil {
		l.Fatalln(err)
	}

	defaultFactory = NewFactory(defaultLogger)
}

// SetLogger defines the default package logger
func SetLogger(instance LoggerFactory) {
	defaultFactory = instance
}

// -----------------------------------------------------------------------------

// Bg delegates a no-context logger
func Bg() Logger {
	return defaultFactory.Bg()
}

// For delegates a context logger
func For(ctx context.Context) Logger {
	return defaultFactory.For(ctx)
}

// Default returns the logger factory
func Default() LoggerFactory {
	return defaultFactory
}
