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
	"io"

	"github.com/pkg/errors"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// CheckErr handles error correctly
func CheckErr(msg string, err error, fields ...zapcore.Field) {
	if err != nil {
		fields = append(fields, zap.Error(err))
		defaultFactory.Bg().Error(msg, fields...)
	}
}

// CheckErrCtx handles error correctly
func CheckErrCtx(ctx context.Context, msg string, err error, fields ...zapcore.Field) {
	if err != nil {
		fields = append(fields, zap.Error(errors.WithStack(err)))
		defaultFactory.For(ctx).Error(msg, fields...)
	}
}

// SafeClose handles the closer error
func SafeClose(c io.Closer, msg string, fields ...zapcore.Field) {
	if cerr := c.Close(); cerr != nil {
		fields = append(fields, zap.Error(errors.WithStack(cerr)))
		defaultFactory.Bg().Error(msg, fields...)
	}
}
