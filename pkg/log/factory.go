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

	opentracing "github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// LoggerFactory defines logger factory contract
type LoggerFactory interface {
	Bg() Logger
	For(context.Context) Logger
	With(...zapcore.Field) LoggerFactory
}

// Factory is the default logging wrapper that can create
// logger instances either for a given Context or context-less.
type factory struct {
	logger *zap.Logger
}

// NewFactory creates a new Factory.
func NewFactory(logger *zap.Logger) LoggerFactory {
	return &factory{logger: logger}
}

// Bg creates a context-unaware logger.
func (b factory) Bg() Logger {
	return &logger{logger: b.logger}
}

// For returns a context-aware Logger. If the context
// contains an OpenTracing span, all logging calls are also
// echo-ed into the span.
func (b factory) For(ctx context.Context) Logger {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		return &spanLogger{span: span, logger: b.logger}
	}
	return b.Bg()
}

// With creates a child logger, and optionally adds some context fields to that logger.
func (b factory) With(fields ...zapcore.Field) LoggerFactory {
	return &factory{logger: b.logger.With(fields...)}
}
