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

package logger

import (
	"context"
	"log/slog"
	"testing"

	. "github.com/onsi/gomega"
)

func TestNewContextWithLogger(t *testing.T) {
	g := NewWithT(t)

	l := &slog.Logger{}
	ctx := NewContextWith(context.Background(), l)
	actual := ctx.Value(loggerContextKey)
	g.Expect(actual).ToNot(BeNil())
	g.Expect(actual).To(Equal(l))
}

func TestGetContextLogger(t *testing.T) {
	t.Run("normal test", func(t *testing.T) {
		g := NewWithT(t)

		l := &slog.Logger{}
		ctx := NewContextWith(context.Background(), l)
		g.Expect(FromContext(ctx)).To(Equal(l))
	})

	t.Run("nil context test", func(t *testing.T) {
		g := NewWithT(t)

		var ctx context.Context
		g.Expect(FromContext(ctx)).To(Equal(slog.Default())) //nolint
	})

	t.Run("no contextKey test", func(t *testing.T) {
		g := NewWithT(t)

		g.Expect(FromContext(context.Background())).To(Equal(slog.Default()))
	})

	t.Run("nil logger test", func(t *testing.T) {
		g := NewWithT(t)

		g.Expect(FromContext(NewContextWith(context.Background(), nil))).To(Equal(slog.Default()))
	})

	t.Run("wrong logger interface test", func(t *testing.T) {
		g := NewWithT(t)

		g.Expect(FromContext(context.WithValue(context.Background(), loggerContextKey, 1))).To(Equal(slog.Default()))
	})
}
