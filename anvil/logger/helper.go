// Copyright (c) 2025 The anyvoxel Authors
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// Package logger contains some helper function for slog
package logger

import (
	"context"
	"log/slog"
)

type loggerContextKeyType int

const (
	loggerContextKey loggerContextKeyType = 0
)

// NewContextWith will create a new context with the value of log,
// And user should use the context returned by this function to replace the context passed in.
// Then user can retrieve the logger from the context use `FromContext` function
// and log by Logger, this will provide contextual-log in different function and module.
func NewContextWith(ctx context.Context, log *slog.Logger) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	return context.WithValue(ctx, loggerContextKey, log)
}

// FromContext will return the logger from the context, if
// passed-in context didn't contain the logger, it will return the default logger.
// User should use this to retrieve logger to ensure log contains information(such as attrs)
// with caller function or modules.
// NOTE: we didn't provide a way to wrap `Info` and other functions to avoid unexpected behavor
// such as user didn't want contextual-log, so the user must explicit to call `FromContext`.
func FromContext(ctx context.Context) *slog.Logger {
	if ctx == nil {
		return slog.Default()
	}

	v := ctx.Value(loggerContextKey)
	if v == nil {
		return slog.Default()
	}

	l, ok := v.(*slog.Logger)
	if !ok || l == nil {
		return slog.Default()
	}

	return l
}
