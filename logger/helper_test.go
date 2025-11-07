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
