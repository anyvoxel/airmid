package xapp

import (
	"context"
	"log/slog"
	"os"
	"reflect"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/anyvoxel/airmid/beans"
)

type testLoggerProvider struct {
}

func (t *testLoggerProvider) GetLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))
}

func TestLoggerStartupHandler(t *testing.T) {
	t.Run("user provide logger", func(t *testing.T) {
		g := NewWithT(t)

		app := &airmidApplication{
			BeanFactory: beans.NewBeanFactory(),
		}
		app.RegisterBeanDefinition("testLoggerProvider",
			beans.MustNewBeanDefinition(
				reflect.TypeOf((*testLoggerProvider)(nil)),
			),
		)
		h := &loggerStartupHandler{}
		err := h.BeforeLoadProps(context.Background(), app, nil)
		g.Expect(err).ToNot(HaveOccurred())

		err = h.AfterLoadProps(context.Background(), app, nil)
		g.Expect(err).ToNot(HaveOccurred())
	})

	t.Run("default handler", func(t *testing.T) {
		g := NewWithT(t)
		app := &airmidApplication{
			BeanFactory: beans.NewBeanFactory(),
		}
		h := &loggerStartupHandler{}
		err := h.BeforeLoadProps(context.Background(), app, nil)
		g.Expect(err).ToNot(HaveOccurred())

		err = h.AfterLoadProps(context.Background(), app, nil)
		g.Expect(err).ToNot(HaveOccurred())
	})

	t.Run("json handler", func(t *testing.T) {
		g := NewWithT(t)
		app := &airmidApplication{
			BeanFactory: beans.NewBeanFactory(),
		}
		h := &loggerStartupHandler{}
		app.Set("airmid.logger.handler.type", "json")
		err := h.BeforeLoadProps(context.Background(), app, nil)
		g.Expect(err).ToNot(HaveOccurred())

		err = h.AfterLoadProps(context.Background(), app, nil)
		g.Expect(err).ToNot(HaveOccurred())
	})

	t.Run("invalid level", func(t *testing.T) {
		g := NewWithT(t)
		app := &airmidApplication{
			BeanFactory: beans.NewBeanFactory(),
		}
		h := &loggerStartupHandler{}
		app.Set("airmid.logger.handler.opt.level", "iii")
		err := h.BeforeLoadProps(context.Background(), app, nil)
		g.Expect(err).ToNot(HaveOccurred())

		err = h.AfterLoadProps(context.Background(), app, nil)
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(MatchRegexp(`level string "iii": unknown name`))
	})

	t.Run("invalid type", func(t *testing.T) {
		g := NewWithT(t)
		app := &airmidApplication{
			BeanFactory: beans.NewBeanFactory(),
		}
		h := &loggerStartupHandler{}
		app.Set("airmid.logger.handler.type", "iii")
		err := h.BeforeLoadProps(context.Background(), app, nil)
		g.Expect(err).ToNot(HaveOccurred())

		err = h.AfterLoadProps(context.Background(), app, nil)
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(MatchRegexp(`unknown logger handler type iii`))
	})
}
